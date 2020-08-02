package xdata

// HostComment represents a 'hostcomment' entry in the Nagios state.dat file.
type HostComment struct {
	Author      string
	CommentData string
	CommentID   int
	EntryTime   int
	EntryType   int
	ExpireTime  int
	Expires     bool
	HostName    string
	Persistent  bool
	Source      int
}
