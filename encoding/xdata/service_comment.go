package xdata

// ServiceComment represents a 'servicecomment' entry in the Nagios state.dat file.
type ServiceComment struct {
	Author             string
	CommentData        string
	CommentID          int
	EntryTime          int
	EntryType          int
	ExpireTime         int
	Expires            bool
	HostName           string
	Persistent         bool
	ServiceDescription string
	Source             int
}
