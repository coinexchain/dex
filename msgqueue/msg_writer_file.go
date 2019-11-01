package msgqueue

import (
	"bytes"
	"io"
	"os"
)

var _ MsgWriter = fileMsgWriter{}

type MkFifoFunc func(path string, mode uint32) (err error)

var mkFifoFunc MkFifoFunc

func SetMkFifoFunc(mkFifo MkFifoFunc) {
	mkFifoFunc = mkFifo
}

type fileMsgWriter struct {
	io.WriteCloser
}

func NewStdOutMsgWriter() MsgWriter {
	return fileMsgWriter{os.Stdout}
}

func NewPipeMsgWriter(pipe string) (MsgWriter, error) {
	file, err := os.OpenFile(pipe, os.O_RDWR, 0666)
	if os.IsNotExist(err) {
		err := mkFifoFunc(pipe, 0666)
		if err != nil {
			return fileMsgWriter{}, err
		}
		file, err = os.OpenFile(pipe, os.O_RDWR, 0666)
		if err != nil {
			return fileMsgWriter{}, err
		}
	} else if err != nil {
		return fileMsgWriter{}, err
	}
	return fileMsgWriter{file}, nil
}

func NewFileMsgWriter(filePath string) (MsgWriter, error) {
	file, err := openFile(filePath)
	if err != nil {
		return fileMsgWriter{}, err
	}
	return fileMsgWriter{file}, nil
}

func (w fileMsgWriter) WriteKV(k, v []byte) error {
	buffer := bytes.NewBuffer(nil)
	buffer.Write(k)
	buffer.Write([]byte("#"))
	buffer.Write(v)
	buffer.Write([]byte("\r\n"))
	if _, err := w.WriteCloser.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func (w fileMsgWriter) Close() error {
	if w.WriteCloser == os.Stdout {
		return nil
	}
	return w.WriteCloser.Close()
}

func (w fileMsgWriter) String() string {
	if w.WriteCloser == os.Stdout {
		return "stdout"
	}
	return "file"
}
