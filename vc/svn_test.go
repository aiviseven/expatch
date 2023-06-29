package vc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRegexp(t *testing.T) {
	s := "123:235"
	r := regexp.MustCompile(`(?P<StartRevision>\w+):(?P<Revision>\w+)`)
	if r.MatchString(s) {
		m := r.FindStringSubmatch(s)
		n := r.SubexpNames()

		result := make(map[string]string)
		for i, name := range n {
			if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
				result[name] = m[i]
			}
		}
		result["BuildTime"] = time.Now().Format("2006-01-02 15:04:05")
		result["SVN"] = "https://127.0.0.1/svn"
		result["Author"] = "one"
		prettyResult, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("%s\n", prettyResult)

		ns := []string{"StartRevision", "Revision", "BuildTime", "Author", "SVN"}
		var bf strings.Builder
		for i, v := range ns {
			bf.WriteString(v)
			bf.WriteString("=")
			bf.WriteString(result[v])
			if i != len(ns)-1 {
				bf.WriteString("\n")
			}
		}
		fmt.Println(bf.String())
	}
}

func TestReplaceHost(t *testing.T) {
	//s := blurHostAndPort("http://117.107.139.62")
	//s := blurHostAndPort("http://117.107.139.62:80")
	//s := blurHostAndPort("https://www.baidu.com")
	//s := blurHostAndPort("https://www.baidu.com:443/")
	s := blurHostAndPort("https://117.107.139.62/ams/AMS2.0/AMS_New_Product/AdminPortal/2.0.0x.0911A/Tags/3.0.1.20121217/IAMServer")
	fmt.Println(s)

}
