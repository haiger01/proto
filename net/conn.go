package net


// Conn interface will contain tcpConn and udpConn
type Conn interface {
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	Close() error
	// LocalAddr() 
	// RemoteAddr()
}