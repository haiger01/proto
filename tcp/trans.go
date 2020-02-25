package tcp

type RetransmissionQueue struct {
	Queue chan []byte
}

type SendSequence struct {
	UNA uint32 // send unacknowladged
	NXT uint32 // send next
	WND uint32 // send window
	UP  uint32 // send urgent pointer
	WL1 uint32 // segment sequence number used for last window update
	WL2 uint32 // segment acknowledgement number used for last window update
	ISS uint32 // initial send sequence number
}

type ReceiveSequence struct {
	NXT uint32 // receive next
	WND uint32 // receive window
	UP  uint32 // receive urgent pointer
	IRS uint32 // initial receive sequence number
}
