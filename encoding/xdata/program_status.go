package xdata

// ProgramStatus represents a 'programstatus' entry in the Nagios state.dat file.
type ProgramStatus struct {
	ActiveHostChecksEnabled          bool
	ActiveOndemandHostCheckStats     string
	ActiveOndemandServiceCheckStats  string
	ActiveScheduledHostCheckStats    string
	ActiveScheduledServiceCheckStats string
	ActiveServiceChecksEnabled       bool
	CachedHostCheckStats             string
	CachedServiceCheckStats          string
	CheckHostFreshness               bool
	CheckServiceFreshness            bool
	DaemonMode                       bool
	EnableEventHandlers              bool
	EnableFlapDectection             bool
	EnableFlapDetection              bool
	EnableNotifications              bool
	ExternalCommandStats             string
	GlobalHostEventHandler           string
	GlobalServiceEventHandler        string
	LastLogRotation                  int
	ModifiedHostAttributes           ModifiedAttribute
	ModifiedServiceAttributes        ModifiedAttribute
	NagiosPID                        int
	NextCommentID                    int
	NextDowntimeID                   int
	NextEventID                      int
	NextNotificationID               int
	NextProblemID                    int
	ObsessOverHosts                  bool
	ObsessOverServices               bool
	ParallelHostCheckStats           string
	PassiveHostCheckStats            string
	PassiveHostChecksEnabled         bool
	PassiveServiceCheckStats         string
	PassiveServiceChecksEnabled      bool
	ProcessPerformanceData           bool
	ProgramStart                     int
	SerialHostCheckStats             string
}
