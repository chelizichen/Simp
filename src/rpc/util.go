package rpc

func queryStructLen(bytes []byte) int32 {
	return int32(bytes[0]) | int32(bytes[1])<<8 | int32(bytes[2])<<16 | int32(bytes[3])<<24
}
