package xdata

// ContactStatus represents a 'contactstatus' entry in the Nagios state.dat file.
type ContactStatus struct {
	ContactName                 string
	HostNotificationPeriod      string
	HostNotificationsEnabled    bool
	LastHostNotification        int
	LastServiceNotification     int
	ModifiedAttributes          int
	ModifiedHostAttributes      int
	ModifiedServiceAttributes   int
	ServiceNotificationPeriod   string
	ServiceNotificationsEnabled bool
}
