package flagsender

type Submitter interface {
	Send(flag string) (Result, error)
	Close() error
}

type Status string

type Result struct {
	Success bool
	Status  string
	Msg     string
}
