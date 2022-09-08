package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(fmt.Sprintf("GetCurrentDirectory: %s", err))
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func CopyFile(srcPath, dstPath string) error {
	file1, err1 := os.Open(srcPath)
	if err1 != nil {
		return err1
	}
	file2, err2 := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err2 != nil {
		return err2
	}
	defer file1.Close()
	defer file2.Close()
	_, err3 := io.Copy(file2, file1)
	if err3 != nil {
		return err3
	}
	return nil
}
