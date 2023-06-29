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
	"regexp"
	"strings"
	"time"
)

//diff文件行前缀，用于解析变更文件的路径
const diffLinePre = "Index: "

type SVN struct {
	CurrPath string
	Args     string
}

//获取版本差异文件
func (s *SVN) GetDiffSet() (set.Set, error) {
	//获取svn差异信息
	svnDiffInfo, err := s.getSvnDiffInfo()
	if err != nil {
		fmt.Printf("svn diff error: %v\n", err)
		return nil, err
	}
	copySet := s.convertSvnDiffList(strings.NewReader(svnDiffInfo))
	return copySet, nil
}

//通过svn diff命令获取差异信息
func (s *SVN) getSvnDiffInfo() (string, error) {
	svnArgsArr := []string{"diff", "--summarize", "-r"}
	if s.Args == "" {
		svnArgsArr = append(svnArgsArr, "HEAD")
		if util.IsCygwin {
			svnArgsArr = append(svnArgsArr, util.WinPathToCyg(s.CurrPath))
		} else {
			svnArgsArr = append(svnArgsArr, s.CurrPath)
		}
	} else {
		for _, v := range strings.Split(s.Args, " ") {
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
func (s *SVN) convertSvnDiffList(reader io.Reader) set.Set {
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

func getSvnUrl() string {
	cmd := exec.Command("svn", "info", "--show-item", "url")
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Start()
	cmd.Wait()
	return string(stdout.Bytes())
}

func (s *SVN) GetSvnInfo(author string) (string, string) {
	if s.Args != "" {
		a := strings.Split(s.Args, " ")
		if vn := strings.TrimSpace(a[0]); vn != "" {
			r := regexp.MustCompile(`(?P<StartRevision>\w+):(?P<Revision>\w+)`)
			if r.MatchString(vn) {
				m := r.FindStringSubmatch(vn)
				n := r.SubexpNames()
				patchInfo := make(map[string]string)
				for i, name := range n {
					if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
						patchInfo[name] = m[i]
					}
				}
				patchInfo["BuildTime"] = time.Now().Format("2006-01-02 15:04:05")
				patchInfo["SVN"] = blurHostAndPort(getSvnUrl())
				patchInfo["Author"] = author

				ns := []string{"StartRevision", "Revision", "BuildTime", "Author", "SVN"}
				var bf strings.Builder
				for i, v := range ns {
					bf.WriteString(v)
					bf.WriteString("=")
					bf.WriteString(patchInfo[v])
					if i != len(ns)-1 {
						bf.WriteString("\n")
					}
				}
				return patchInfo["Revision"], bf.String()
			}
		}
	}
	return "", ""
}

func blurHostAndPort(url string) string {
	// 正则表达式模式
	pattern := `(https?://)([^:/\s]+)(:\d+)?(/.*)?`
	// 使用正则表达式提取主机和端口
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(url)
	if len(matches) > 0 {
		protocol := matches[1]
		host := "XX.XX.XX.XX"
		port := "" // 可以选择保留端口号或替换为模糊化标记
		path := matches[4]

		// 组合模糊化的URL
		blurredURL := protocol + host + port + path
		return blurredURL
	}

	return url
}
