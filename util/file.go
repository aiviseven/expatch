package util

import (
	"bufio"
	"fmt"
	set "github.com/deckarep/golang-set"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// CopyFile 复制文件
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	//判断目标文件路径是否存在，不存在则创建
	if lastSepIndex := strings.LastIndex(dst, "/"); lastSepIndex >= 0 {
		if isExist, _ := PathExists(dst[:lastSepIndex]); !isExist {
			err := os.MkdirAll(dst[:lastSepIndex], os.ModePerm)
			if err != nil {
				return 0, err
			}
		}
	}

	//destination, err := os.Create(dst)
	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// ReadFile 读取xml文件为[]byte
func ReadFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return data
}

// ReadFileToSet 按行读取文件为set
func ReadFileToSet(filePath string) set.Set {
	if filePath == "" {
		return nil
	}

	s := set.NewSet()
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}

	r := bufio.NewReader(f)
	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		s.Add(string(l))
	}
	return s
}
