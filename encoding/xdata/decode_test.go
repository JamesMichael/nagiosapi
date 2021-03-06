package xdata

import (
	"fmt"
	"strings"
	"testing"
)

const sampleInput = `
########################################
#          NAGIOS STATUS FILE
#
# THIS FILE IS AUTOMATICALLY GENERATED
# BY NAGIOS.  DO NOT MODIFY THIS FILE!
########################################

info {
	created=12345
	version=4.0.0
	last_update_check=67890
	update_available=1
	last_version=4.0.0
	new_version=4.4.5
}

programstatus {
	next_comment_id=123
}

hoststatus {
	host_name=host1
}

hoststatus {
	host_name=host2
}

servicestatus {
	host_name=host1
	service_description=service 1
	check_interval=5.0000
}

servicestatus {
	host_name=host1
	service_description=service 2
	notifications_enabled=0
}

hostcomment {
	host_name=host1
}

servicecomment {
	host_name=host1
	service_description=service 1
}
`

func TestDecoder_Decode(t *testing.T) {
	dec := NewDecoder(strings.NewReader(sampleInput))

	var res Status
	if err := dec.Decode(&res); err != nil {
		t.Errorf("unable to decode sample input: %s", err)
		return
	}

	if res.Info == nil {
		t.Errorf("failed to parse info block")
	} else {
		if got := res.Info.Created; got != 12345 {
			t.Errorf("incorrect info.created, got: %d, expected: %d", got, 12345)
		}

		if got := res.Info.Version; got != "4.0.0" {
			t.Errorf("incorrect info.version, got: %s, expected: %s", got, "4.0.0")
		}

		if got := res.Info.LastUpdateCheck; got != 67890 {
			t.Errorf("incorrect info.last_update_check, got: %d, expected: %d", got, 67890)
		}

		if got := res.Info.UpdateAvailable; !got {
			t.Errorf("incorrect info.update_available, got: %t, expected: %t", got, true)
		}

		if got := res.Info.LastVersion; got != "4.0.0" {
			t.Errorf("incorrect info.last_update_check, got: %s, expected: %s", got, "4.0.0")
		}

		if got := res.Info.NewVersion; got != "4.4.5" {
			t.Errorf("incorrect info.new_version, got: %s, expected: %s", got, "4.4.5")
		}
	}

	if res.ProgramStatus == nil {
		t.Errorf("failed to parse programstatus block")
	} else {
		if got := res.ProgramStatus.NextCommentID; got != 123 {
			t.Errorf("incorrect programstatus.next_comment_id, got: %d, expected: %d", got, 123)
		}
	}

	if res.HostStatus == nil || len(res.HostStatus) != 2 {
		t.Errorf("failed to parse hoststatus blocks")
	} else {
		if got := res.HostStatus[0].HostName; got != "host1" {
			t.Errorf("incorrect hoststatus.host_name, got: %s, expected: %s", got, "host1")
		}

		if got := res.HostStatus[1].HostName; got != "host2" {
			t.Errorf("incorrect hoststatus.host_name, got: %s, expected: %s", got, "host2")
		}
	}

	if res.ServiceStatus == nil || len(res.ServiceStatus) != 2 {
		t.Errorf("failed to parse servicestatus blocks")
	} else {
		if got := res.ServiceStatus[0].HostName; got != "host1" {
			t.Errorf("incorrect servicestatus.host_name, got: %s, expected: %s", got, "host1")
		}

		if got := res.ServiceStatus[0].CheckInterval; got != 5.0 {
			t.Errorf("incorrect servicestatus.check_interval, got: %f, expected: %f", got, 5.0)
		}

		if got := res.ServiceStatus[1].ServiceDescription; got != "service 2" {
			t.Errorf("incorrect servicestatus.service_description, got: %s, expected: %s", got, "service 2")
		}
	}

	if res.HostComment == nil || len(res.HostComment) != 1 {
		t.Errorf("failed to parse hostcomment blocks")
	} else {
		if got := res.HostComment[0].HostName; got != "host1" {
			t.Errorf("incorrect hostcomment.host_name, got: %s, expected: %s", got, "host1")
		}
	}

	if res.ServiceComment == nil || len(res.ServiceComment) != 1 {
		t.Errorf("failed to parse servicecomment blocks")
	} else {
		if got := res.ServiceComment[0].ServiceDescription; got != "service 1" {
			t.Errorf("incorrect servicecomment.service_name, got: %s, expected: %s", got, "service 1")
		}
	}
}

func TestDecoder_Decode_InvalidReciever(t *testing.T) {
	dec := NewDecoder(strings.NewReader(sampleInput))

	var res string
	if err := dec.Decode(&res); err == nil {
		t.Errorf("expected error")
		return
	}
}

func TestDecoder_Decode_NonPointerReceiver(t *testing.T) {
	dec := NewDecoder(strings.NewReader(sampleInput))

	var res string
	if err := dec.Decode(res); err == nil {
		t.Errorf("expected error")
		return
	}
}

func TestDecoder_Decode_UnknownBlock(t *testing.T) {
	const sampleInput = `
		info {
			created=12345
		}

		unknown {
		}
	`
	dec := NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidLines = true

	var res Status
	if err := dec.Decode(&res); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if got := res.Info.Created; got != 12345 {
		t.Errorf("incorrect info.created, got: %d, expected: %d", got, 12345)
	}

	dec = NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidLines = false
	err := dec.Decode(&res)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDecoder_Decode_InvalidInt(t *testing.T) {
	const sampleInput = `
		info {
			created=today
		}
	`

	dec := NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = true

	var res Status
	if err := dec.Decode(&res); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if got := res.Info.Created; got != 0 {
		t.Errorf("unexpected value for info.created, got: %d, expected: %d", got, 0)
	}

	dec = NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = false

	err := dec.Decode(&res)
	if err == nil {
		t.Errorf("expected error")
		return
	}
}

func TestDecoder_Decode_InvalidFloat(t *testing.T) {
	const sampleInput = `
		servicestatus {
			check_interval=hourly
		}
	`

	dec := NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = true

	var res Status
	if err := dec.Decode(&res); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if got := res.ServiceStatus[0].CheckInterval; got != 0.0 {
		t.Errorf("unexpected value for info.created, got: %f, expected: %f", got, 0.0)
	}

	dec = NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = false

	err := dec.Decode(&res)
	if err == nil {
		t.Errorf("expected error")
		return
	}
}

func TestDecoder_Decode_InvalidBool(t *testing.T) {
	const sampleInput = `
		servicestatus {
			notifications_enabled=true
		}
	`

	dec := NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = true

	var res Status
	if err := dec.Decode(&res); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if got := res.ServiceStatus[0].NotificationsEnabled; got != false {
		t.Errorf("unexpected value for info.created, got: %t, expected: %t", got, false)
	}

	dec = NewDecoder(strings.NewReader(sampleInput))
	dec.IgnoreInvalidTypes = false

	err := dec.Decode(&res)
	if err == nil {
		t.Errorf("expected error")
		return
	}
}

func ExampleDecoder_Decode() {
	const input = `
# NAGIOS STATUS FILE
info {
	version=4.0.0
}

servicestatus {
	host_name=example.com
}
`

	var res Status
	if err := NewDecoder(strings.NewReader(input)).Decode(&res); err != nil {
		panic(err)
	}
	fmt.Println(res.ServiceStatus[0].HostName)
}
