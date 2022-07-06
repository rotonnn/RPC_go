package protocol

type Request struct {
	ServicePath string
	Payload     []byte
}

func (r Request) Reset() {
	//TODO implement me
	panic("implement me")
}

func (r Request) String() string {
	//TODO implement me
	panic("implement me")
}

func (r Request) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}

type Response struct {
	ErrCode int
	ErrTips string

	Payload []byte
}

func (r Response) Reset() {
	//TODO implement me
	panic("implement me")
}

func (r Response) String() string {
	//TODO implement me
	panic("implement me")
}

func (r Response) ProtoMessage() {
	//TODO implement me
	panic("implement me")
}
