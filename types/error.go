package types

type ErrorMonitor struct {
	Title   string
	Details string
}

func (e ErrorMonitor) Error() string {
	return e.Title
}
