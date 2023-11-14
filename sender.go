package flagsender

type Sender interface {
	Send(flag string) (Result, error)
	Close() error
}

type Result struct {
	Success bool
	Status  string
	Msg     string
}
