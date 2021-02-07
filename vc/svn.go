package vc

import (
	"bufio"
	"fmt"
	"github.com/bsh100220/expatch/util"
	set "github.com/deckarep/golang-set"
	"io"
	"os/exec"
	"strings"
)

//
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
	svnArgsArr := []string{"diff", "-r"}
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
	buf, err := cmd.Output()
	return fmt.Sprintf("%s\n", buf), err
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

		if strings.HasPrefix(line, diffLinePre) {
			line = line[strings.Index(line, diffLinePre)+len(diffLinePre):]
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
