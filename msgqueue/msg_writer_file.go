package msgqueue

import (
	"io"
	"os"
	"syscall"
)

var _ MsgWriter = fileMsgWriter{}

type fileMsgWriter struct {
	io.WriteCloser
}

func NewStdOutMsgWriter() MsgWriter {
	return fileMsgWriter{os.Stdout}
}

func NewPipeMsgWriter(pipe string) (MsgWriter, error) {
	file, err := os.OpenFile(pipe, os.O_RDWR, 0666)
	if os.IsNotExist(err) {
		err := syscall.Mkfifo(pipe, 0666)
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
	if _, err := w.WriteCloser.Write(k); err != nil {
		return err
	}
	if _, err := w.WriteCloser.Write([]byte("#")); err != nil {
		return err
	}
	if _, err := w.WriteCloser.Write(v); err != nil {
		return err
	}
	if _, err := w.WriteCloser.Write([]byte("\r\n")); err != nil {
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
