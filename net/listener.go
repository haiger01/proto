package net

type Listener interface {
	Accept() (Conn, error)
	Close() error
	// Addr()
}
