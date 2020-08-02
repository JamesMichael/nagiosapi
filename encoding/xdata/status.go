package xdata

type Status struct {
	Info           *Info
	ProgramStatus  *ProgramStatus
	HostStatus     []*HostStatus
	HostComment    []*HostComment
	ServiceStatus  []*ServiceStatus
	ServiceComment []*ServiceComment
}
