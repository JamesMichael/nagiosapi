package xdata

// Info represents an 'info' entry in the Nagios state.dat file.
type Info struct {
	Created         int
	LastUpdateCheck int
	LastVersion     string
	NewVersion      string
	UpdateAvailable bool
	Version         string
}
