//java项目的增量打包工具
package main

import (
	"flag"
	"fmt"
	"github.com/bsh100220/expatch/ide"
	"github.com/bsh100220/expatch/util"
	"github.com/bsh100220/expatch/vc"
	"strings"
)

var currPath, projectType, projectConfigFilePath, patchOutPath string

func main() {
	//获取当前路径的绝对路径
	currPath = util.GetAbsolutePath(util.GetCurrentDirectory())

	//解析参数
	flag.StringVar(&patchOutPath, "out", currPath+"/patch"+util.GetNowStr(), "输出目录，默认为当前目录下的[patch+当前时间]目录\n")
	flag.StringVar(&projectType, "type", "idea", "项目类型，可选(idea,eclipse)，默认为idea")
	flag.StringVar(&projectConfigFilePath, "conf", "", "idea模块配置文件路径，当type=idea时候有效")
	svnArgs := flag.String("svn", "", "需要对比的两个版本，版本号用英文冒号(:)分隔，带有文件路径则用空格分割，例：'100:95'、'100 .'，为空时则对比当前目录与svn最新版本")
	flag.Parse()

	//解析项目配置文件，得到src、output、WebRoot等路径
	javaProject, err := ide.GetInstance(strings.ToLower(projectType), projectConfigFilePath).AnalysisProjectConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	if javaProject == nil {
		fmt.Println("解析项目配置信息异常")
		return
	}

	if javaProject.ClassesOutPath == "" {
		fmt.Println("未找到class文件输出目录")
		return
	}

	if javaProject.JavaSrcPaths == nil || len(javaProject.JavaSrcPaths) == 0 {
		fmt.Println("未找到java源文件目录")
		return
	}

	//解析svn差异信息得到差异文件路径
	svn := &vc.SVN{CurrPath: currPath, Args: *svnArgs}
	copySet, err := svn.GetDiffSet()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//获取补丁输出路径
	patchOutPath, err = util.Mkdir(patchOutPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//复制差异文件
	err = ide.ConvertPatchDir(copySet, javaProject, patchOutPath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Println("Done!")
}
