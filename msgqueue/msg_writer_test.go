package msgqueue

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMsgWriter(t *testing.T) {
	w, err := createMsgWriter("nop")
	require.NoError(t, err)
	require.Equal(t, "nop", w.String())

	w, err = createMsgWriter("os:stdout")
	require.NoError(t, err)
	require.Equal(t, "stdout", w.String())
	require.NoError(t, w.Close())

	defer os.Remove("messages.txt")
	w, err = createMsgWriter("file:messages.txt")
	require.NoError(t, err)
	require.Equal(t, "file", w.String())
	require.NoError(t, w.Close())

	w, err = createMsgWriter("kafka:a,b,c")
	require.Error(t, err)
	require.Equal(t, "kafka", w.String())

	w, err = createMsgWriter("db:mongo")
	require.Error(t, err)

	w, err = createMsgWriter("dir:tmp")
	require.NoError(t, err)
	defer os.RemoveAll("tmp")
	require.Equal(t, "dir", w.String())
	require.NoError(t, w.Close())
}

func TestNopMsgWriter(t *testing.T) {
	w := NewNopMsgWriter()
	require.Equal(t, "nop", w.String())
	require.NoError(t, w.WriteKV(nil, nil))
	require.NoError(t, w.Close())
}

func TestFileMsgWriter(t *testing.T) {
	// new file
	defer os.Remove("messages.txt")
	w, err := NewFileMsgWriter("messages.txt")
	require.NoError(t, err)

	require.NoError(t, w.WriteKV([]byte("k1"), []byte("v1")))
	require.NoError(t, w.WriteKV([]byte("k2"), []byte("v2")))
	require.NoError(t, w.Close())

	data, err := ioutil.ReadFile("messages.txt")
	require.NoError(t, err)
	require.Equal(t, "k1#v1\r\nk2#v2\r\n", string(data))

	// existed file
	w, err = NewFileMsgWriter("messages.txt")
	require.NoError(t, err)
	require.NoError(t, w.Close())

	// dir
	w, err = NewFileMsgWriter(".")
	require.Error(t, err)
	//require.Equal(t, "Need to give the file path ", err.Error())
}

func TestDirMsgWriter(t *testing.T) {
	w, err := NewDirMsgWriter("tmp")
	require.NoError(t, err)
	defer os.RemoveAll("tmp")

	MaxFileSize = 20
	err = w.WriteKV([]byte("hello world"), []byte("nihao"))
	require.NoError(t, err)
	err = w.WriteKV([]byte("hello world"), []byte("nihao"))
	require.NoError(t, err)
	err = w.WriteKV([]byte("hello world"), []byte("nihao"))
	require.NoError(t, err)

	info, err := os.Stat("tmp/" + filePrefix + "0")
	require.NoError(t, err)
	require.EqualValues(t, 11+5+3, info.Size())
	bz, err := ioutil.ReadFile("tmp/" + filePrefix + "0")
	require.NoError(t, err)
	require.EqualValues(t, []byte("hello world#nihao\r\n"), bz)

	info, err = os.Stat("tmp/" + filePrefix + "1")
	require.NoError(t, err)
	require.EqualValues(t, 11+5+3, info.Size())
	bz, err = ioutil.ReadFile("tmp/" + filePrefix + "1")
	require.NoError(t, err)
	require.EqualValues(t, []byte("hello world#nihao\r\n"), bz)

	info, err = os.Stat("tmp/" + filePrefix + "2")
	require.NoError(t, err)
	require.EqualValues(t, 11+5+3, info.Size())
	bz, err = ioutil.ReadFile("tmp/" + filePrefix + "2")
	require.NoError(t, err)
	require.EqualValues(t, []byte("hello world#nihao\r\n"), bz)

	require.Nil(t, w.Close())
}
