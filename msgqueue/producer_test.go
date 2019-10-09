package msgqueue

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProducer(t *testing.T) {
	defer os.Remove("messages.txt")
	p := NewProducerFromConfig([]string{"file:messages.txt", "dir:tmp", "os:stdout"}, "bank,auth", true, nil)
	require.True(t, p.IsOpenToggle())
	require.True(t, p.IsSubscribed("bank"))
	require.True(t, p.IsSubscribed("auth"))
	require.False(t, p.IsSubscribed("gov"))
	require.Equal(t, []string{"file", "dir", "stdout"}, p.GetMode())

	p.SendMsg([]byte("foo"), []byte("bar"))
	p.Close()

	data, err := ioutil.ReadFile("messages.txt")
	require.NoError(t, err)
	require.Equal(t, "foo#bar\r\n", string(data))
}
