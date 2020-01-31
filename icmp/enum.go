package icmp

const (
	EchoReply              ICMPType = 8
	DestinationUnreachable ICMPType = 3
	SourceQuench           ICMPType = 4
	Redirect               ICMPType = 5
	Echo                   ICMPType = 0
	RouterAdvertisement    ICMPType = 9
	RouterSolicitation     ICMPType = 10
	TimeExceeded           ICMPType = 11
	ParameterProblem       ICMPType = 12
	Timestamp              ICMPType = 13
	TimestampReply         ICMPType = 14
	InformationRequest     ICMPType = 15
	InformationReply       ICMPType = 16
	AddressMaskRequest     ICMPType = 17
	AddressMaskReply       ICMPType = 18
)

const (
	EchoReplyCode uint8 = 0

	EchoRequestCode uint8 = 0

	DestinationNetworkUnreachableCode  uint8 = 0
	DestinationHostUnreachableCode     uint8 = 1
	DestinationProtocolUnreachableCode uint8 = 2
	DestinationPortUnreachableCode     uint8 = 3
	FragmentationRequiredCode          uint8 = 4
	SourceRouteFailedCode              uint8 = 5
	DestinationNetworkUnknownCode      uint8 = 6
	DestinationHostUnknown             uint8 = 7

	TTLExpired           uint8 = 0
	FragmentTimeExceeded uint8 = 1
)
