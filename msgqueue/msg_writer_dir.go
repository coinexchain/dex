package msgqueue

import (
	"fmt"
	"io"
)

const (
	filePrefix = "backup-"
)

var MaxFileSize = 1024 * 1024 * 100

var _ MsgWriter = (*dirMsgWriter)(nil)

type dirMsgWriter struct {
	io.WriteCloser
	haveWriteSize int
	fileIndex     int
	dir           string
}

func NewDirMsgWriter(dir string) (MsgWriter, error) {
	filePath, fileIndex, err := GetFilePathAndFileIndexFromDir(dir, MaxFileSize)
	if err != nil {
		return &dirMsgWriter{}, err
	}
	file, err := openFile(filePath)
	if err != nil {
		return &dirMsgWriter{}, err
	}
	fileSize := GetFileSize(filePath)
	if fileSize < 0 {
		return &dirMsgWriter{}, fmt.Errorf("The parameter passed in is not the correct file path. ")
	}
	return &dirMsgWriter{
		WriteCloser:   file,
		fileIndex:     fileIndex,
		dir:           dir,
		haveWriteSize: fileSize,
	}, nil
}

func (w *dirMsgWriter) WriteKV(k, v []byte) error {
	if len(k)+len(v)+3+w.haveWriteSize > MaxFileSize {
		if err := w.Close(); err != nil {
			return err
		}
		file, err := openFile(GetFileName(w.dir, w.fileIndex+1))
		if err != nil {
			return err
		}
		w.WriteCloser = file
		w.fileIndex++
		w.haveWriteSize = 0
	}
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
	w.haveWriteSize += len(k) + len(v) + 3
	return nil
}

func (w *dirMsgWriter) Close() error {
	return w.WriteCloser.Close()
}

func (w *dirMsgWriter) String() string {
	return "dir"
}
