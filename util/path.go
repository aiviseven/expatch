package util

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var IsCygwin bool

func init()  {
	checkOs()
}

//检查运行的系统，主要针对windows中使用cygwin时路径有差异的问题
func checkOs() {
	cmd := exec.Command("uname")
	buf, _ := cmd.Output()
	if strings.Contains(strings.ToUpper(string(buf)), "CYGWIN") {
		IsCygwin = true
	}
}

//获取当前路径
func GetCurrentDirectory() string {
	//dir, _ := os.Executable()
	//exePath := filepath.Dir(dir)
	exePath, _ := os.Getwd()
	return exePath
}

//获取真实的绝对路径，并把路径分割符换成'/'
func GetAbsolutePath(path string) (result string) {
	if IsCygwin {
		cmd := exec.Command("cygpath", "--absolute", "--windows", path)
		buf, _ := cmd.Output()
		result = strings.TrimSpace(string(buf))
	} else {
		result, _ = filepath.Abs(path)
	}
	result = filepath.ToSlash(result)
	return result
}

//将windows路径转换为cygwin路径
func WinPathToCyg(winPath string) string {
	if !IsCygwin {
		return winPath
	}
	cmd := exec.Command("cygpath", "-p", winPath, "-a", "-u")
	buf, _ := cmd.Output()
	result := strings.TrimSpace(string(buf))
	return result
}

//检查路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//判断是否是src路径，如果是则返回src路径
func PathContains(paths []string, filePath string) (bool, string) {
	for _, path := range paths {
		if strings.Contains(filePath, path) {
			return true, path
		}
	}
	return false, ""
}

//创建目录
func Mkdir(dirPath string) (string, error){
	dirPath = GetAbsolutePath(dirPath)
	//如果输出路径不存在则新建
	if isExist, _ := PathExists(dirPath); !isExist {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return "", errors.New("创建输出路径出错！")
		}
	}
	return dirPath, nil
}
