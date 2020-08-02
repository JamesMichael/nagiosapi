package xdata

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Decoder is used to Decode a xdata file.
type Decoder struct {
	// When set, IgnoreInvalidTypes will not cause a decode error if a type cannot be converted from
	// the string value into the reciever value.
	IgnoreInvalidTypes bool

	// When set, IgnoreInvalidLines will not cause a decode error if an unexpected line is encountered.
	IgnoreInvalidLines bool

	r *bufio.Reader
}

var (
	reIgnore     = regexp.MustCompile(`^\s*(#|$)`)
	reBlockStart = regexp.MustCompile(`^\s*(\w+)\s*{`)
	reBlockEnd   = regexp.MustCompile(`^\s*}`)
	reKV         = regexp.MustCompile(`^\s*(.*?)\s*=\s*(.*)\s*$`)
)

// NewDecoder takes a reader and returns a Decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}

// Decode attempts to decode the file into the receiever.
//
// An ErrDecodeError is returned if the file cannot be read into the provided reciever.
func (dec *Decoder) Decode(s interface{}) error {
	st := reflect.TypeOf(s)
	if st.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer input")
	}

	sv := reflect.ValueOf(s).Elem()
	if sv.Type().Kind() != reflect.Struct {
		fmt.Println(sv.Type().Kind().String())
		return fmt.Errorf("expected pointer to struct input")
	}

	blocks := map[string]reflect.Value{}
	for i := 0; i < sv.NumField(); i++ {
		fieldName := sv.Type().Field(i).Name
		configName := strings.ToLower(fieldName)

		field := sv.Field(i)
		if field.IsNil() && field.Kind() == reflect.Slice {
			field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		}
		blocks[configName] = field
	}

	for {
		line, err := dec.r.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if reIgnore.Match(line) {
			continue
		}

		block := reBlockStart.FindSubmatch(line)
		if len(block) == 2 {
			field, ok := blocks[string(block[1])]
			if !ok {
				for {
					line, err := dec.r.ReadBytes('\n')
					if err == io.EOF {
						return nil
					}
					if err != nil {
						return err
					}

					if reBlockEnd.Match(line) {
						break
					}
				}
				// continue until the end of the block
				continue
			}

			if err := dec.decodeBlock(field); err != nil {
				return err
			}
			continue
		}

		if !dec.IgnoreInvalidLines {
			return fmt.Errorf("invalid line encountered '%s'", line)
		}
	}

	return nil
}

func (dec *Decoder) decodeBlock(v reflect.Value) error {
	var result reflect.Value

	switch v.Kind() {
	case reflect.Slice:
		sliceValue := v.Type().Elem()
		switch sliceValue.Kind() {
		case reflect.Ptr:
			result = reflect.New(sliceValue.Elem())
			v.Set(reflect.Append(v, result))
			result = reflect.Indirect(result)
		default:
			return fmt.Errorf("invalid reciever type '%s'", sliceValue.Kind().String())
		}

	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		result = v.Elem()

	case reflect.Struct:
		result = v
	}

	if result.Kind() != reflect.Struct {
		return fmt.Errorf("invalid receiver type '%s', expected struct", result.Kind().String())
	}

	fields := map[string]reflect.Value{}
	for i := 0; i < result.NumField(); i++ {
		fieldName := result.Type().Field(i).Name
		configName := toSnakeCase(fieldName)
		fields[configName] = result.Field(i)
	}

	for {
		line, err := dec.r.ReadBytes('\n')
		if err != nil {
			return err
		}

		if reIgnore.Match(line) {
			continue
		}

		if reBlockEnd.Match(line) {
			break
		}

		kv := reKV.FindSubmatch(line)
		if len(kv) != 3 {
			if dec.IgnoreInvalidLines {
				return nil
			}
			return fmt.Errorf("invalid line, expected KEY=VALUE, got '%s'", line)
		}

		field, ok := fields[string(kv[1])]
		if !ok {
			continue
		}

		switch field.Kind() {

		case reflect.String:
			field.SetString(string(kv[2]))

		case reflect.Int:
			i, err := strconv.Atoi(string(kv[2]))
			if err != nil {
				if !dec.IgnoreInvalidTypes {
					return fmt.Errorf("unable to convert '%s' to int", kv[2])
				}
				continue
			}
			field.SetInt(int64(i))

		case reflect.Bool:
			switch string(kv[2]) {
			case "0":
				field.SetBool(false)
			case "1":
				field.SetBool(true)
			default:
				if !dec.IgnoreInvalidTypes {
					return fmt.Errorf("unable to convert '%s' to bool", kv[2])
				}
			}

		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(string(kv[2]), 64)
			if err != nil {
				if !dec.IgnoreInvalidTypes {
					return fmt.Errorf("unable to convert '%s' to float", kv[2])
				}
			} else {
				field.SetFloat(f)
			}

		default:
			if !dec.IgnoreInvalidTypes {
				return fmt.Errorf("unable to parse value of type '%s'", field.Kind().String())
			}
		}
	}

	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
