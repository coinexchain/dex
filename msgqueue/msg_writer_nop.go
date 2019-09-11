package msgqueue

var _ MsgWriter = nopMsgWriter{}

type nopMsgWriter struct {
}

func NewNopMsgWriter() MsgWriter {
	return nopMsgWriter{}
}

func (w nopMsgWriter) WriteKV(k, v []byte) error {
	// do nothing
	return nil
}

func (w nopMsgWriter) Close() error {
	// do nothing
	return nil
}

func (w nopMsgWriter) String() string {
	return "nop"
}
