package msgqueue

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFileSize(t *testing.T) {
	filePath := "tmp.out"
	require.EqualValues(t, -1, GetFileSize(filePath))

	file, err := os.Create(filePath)
	require.Nil(t, err)
	require.EqualValues(t, 0, GetFileSize(filePath))
	defer os.Remove(filePath)
	lenth, err := file.Write([]byte("hello world"))
	require.Nil(t, err)
	require.EqualValues(t, lenth, GetFileSize(filePath))
}

func TestGetFileName(t *testing.T) {
	require.EqualValues(t, "dir/backup-0", GetFileName("dir", 0))
	require.EqualValues(t, "dir/backup-192", GetFileName("dir", 192))
	require.NotEqual(t, "dir/backup-0", GetFileName("dir", 1))
}

func TestGetFilePathAndFileIndexFromDir(t *testing.T) {
	dirPath := "tmp"
	filePath, fileIndex, err := GetFilePathAndFileIndexFromDir(dirPath, 10)
	require.Nil(t, err)
	require.Equal(t, "tmp/backup-0", filePath)
	require.Equal(t, 0, fileIndex)
	defer os.RemoveAll(dirPath)

	filePath, fileIndex, err = GetFilePathAndFileIndexFromDir(dirPath, 10)
	require.Nil(t, err)
	require.Equal(t, "tmp/backup-0", filePath)
	require.Equal(t, 0, fileIndex)

	file1 := "tmp/backup-2"
	file2 := "tmp/backup-3"
	_, err = os.Create(file1)
	require.Nil(t, err)
	defer os.Remove(file1)
	f2, err := os.Create(file2)
	require.Nil(t, err)
	defer os.Remove(file2)
	filePath, fileIndex, err = GetFilePathAndFileIndexFromDir(dirPath, 10)
	require.Nil(t, err)
	require.Equal(t, "tmp/backup-3", filePath)
	require.Equal(t, 3, fileIndex)

	_, err = f2.Write([]byte("hello world\n"))
	require.Nil(t, err)
	filePath, fileIndex, err = GetFilePathAndFileIndexFromDir(dirPath, 10)
	require.Nil(t, err)
	require.Equal(t, "tmp/backup-4", filePath)
	require.Equal(t, 4, fileIndex)

}

func TestOpenFile(t *testing.T) {
	_, err := openFile(".")
	require.Error(t, err)

	file, err := openFile("test.txt")
	require.Nil(t, err)
	defer os.Remove("test.txt")
	err = file.Close()
	require.Nil(t, err)

	info, err := os.Stat("test.txt")
	require.Nil(t, err)
	require.EqualValues(t, false, info.IsDir())

}
