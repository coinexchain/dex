package msgqueue

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetFileName(dir string, fileIndex int) string {
	return dir + "/" + filePrefix + strconv.Itoa(fileIndex)
}

func GetFilePathAndFileIndexFromDir(dir string, maxFileSize int) (filePath string, fileIndex int, err error) {
	files, err := ioutil.ReadDir(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return "", -1, err
		}
		return GetFileName(dir, 0), 0, nil
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasPrefix(name, filePrefix) {
				fileNames = append(fileNames, name)
			}
		}
	}
	if len(fileNames) == 0 {
		return GetFileName(dir, 0), 0, nil
	}

	fileIndex = 0
	for _, fileName := range fileNames {
		vals := strings.Split(fileName, "-")
		if len(vals) == 2 {
			if index, err := strconv.Atoi(vals[1]); err == nil {
				if index > fileIndex {
					fileIndex = index
				}
			}
		}
	}

	fileSize := GetFileSize(GetFileName(dir, fileIndex))
	if fileSize < maxFileSize {
		return GetFileName(dir, fileIndex), fileIndex, nil
	}
	return GetFileName(dir, fileIndex+1), fileIndex + 1, nil
}

func GetFileSize(filePath string) int {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		return -1
	}

	return int(info.Size())
}

func openFile(filePath string) (*os.File, error) {
	if s, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.Create(filePath)
	} else if s.IsDir() {
		return nil, fmt.Errorf("Need to give the file path ")
	} else {
		return os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	}
}
