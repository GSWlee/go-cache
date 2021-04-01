package mycache

type Byteviews struct {
	b []byte
}

func (receiver Byteviews) Len() int {
	return len(receiver.b)
}

func (receiver Byteviews) cloneByte() []byte {
	temp:=make([]byte,len(receiver.b))
	copy(temp,receiver.b)
	return temp
}

func (receiver Byteviews) ByteSlice() []byte {
	return receiver.cloneByte()
}

func (receiver Byteviews) String() string {
	return string(receiver.b)
}


