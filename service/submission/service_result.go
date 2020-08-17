package submission

// ServiceResult represents the result of a passive service check.
type ServiceResult struct {
	Time        int64  `json:"time"`
	Hostname    string `json:"hostname"`
	ServiceName string `json:"service_name"`
	Status      uint   `json:"status"`
	Body        string `json:"body"`
}
