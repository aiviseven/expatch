package vc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aiviseven/expatch/util"
	set "github.com/deckarep/golang-set"
	"io"
	"os"
	"os/exec"
	"strings"
)

//diff文件行前缀，用于解析变更文件的路径
const diffLinePre = "Index: "

type SVN struct {
	CurrPath string
	Args     string
}

//获取版本差异文件
func (svn *SVN) GetDiffSet() (set.Set, error) {
	//获取svn差异信息
	svnDiffInfo, err := svn.getSvnDiffInfo(svn.Args)
	if err != nil {
		fmt.Printf("svn diff error: %v\n", err)
		return nil, err
	}
	copySet := svn.convertSvnDiffList(strings.NewReader(svnDiffInfo))
	return copySet, nil
}

//通过svn diff命令获取差异信息
func (svn *SVN) getSvnDiffInfo(svnArgs string) (string, error) {
	svnArgsArr := []string{"diff", "--summarize", "-r"}
	if svnArgs == "" {
		svnArgsArr = append(svnArgsArr, "HEAD")
		if util.IsCygwin {
			svnArgsArr = append(svnArgsArr, util.WinPathToCyg(svn.CurrPath))
		} else {
			svnArgsArr = append(svnArgsArr, svn.CurrPath)
		}
	} else {
		for _, v := range strings.Split(svnArgs, " ") {
			svnArgsArr = append(svnArgsArr, v)
		}
	}

	cmd := exec.Command("svn", svnArgsArr...)

	//防止出现需要交互输入密码的情况会导致报错
	//buf, err := cmd.Output()
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Start()
	cmd.Wait()

	rst := fmt.Sprintf("%s\n", string(stdout.Bytes()))
	//fmt.Printf("svn diff info:\n%s", rst)
	return rst, nil
}

//读取svn差异文件得到有变更的文件列表
func (svn *SVN) convertSvnDiffList(reader io.Reader) set.Set {
	copyFileSet := set.NewSet()
	delFileSet := set.NewSet()

	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		line = strings.TrimSpace(line)

		if line != "" {
			line = strings.TrimSpace(line[1:])
			line = util.GetAbsolutePath(line)
			if isExist, _ := util.PathExists(line); isExist {
				copyFileSet.Add(line)
			} else {
				delFileSet.Add(line)
			}
		}
	}

	//fmt.Println(copyFileSet)
	//fmt.Println(delFileSet)

	return copyFileSet
}
