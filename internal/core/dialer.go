package core

type Dialer interface {
	Name() string
	Dial(addr string) (Session, error)
}
