package monitorsDTO

type CreateMonitorRequest struct {
	Url             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
}
