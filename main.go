//java项目的增量打包工具
package main

import (
	"flag"
	"fmt"
	"github.com/aiviseven/expatch/ide"
	"github.com/aiviseven/expatch/util"
	"github.com/aiviseven/expatch/vc"
	"os"
	"strings"
)

var currPath, projectType, projectConfigFilePath, patchOutPath, ignoreConfFilePath, author string

func main() {
	//获取当前路径的绝对路径
	currPath = util.GetAbsolutePath(util.GetCurrentDirectory())

	//解析参数
	flag.StringVar(&patchOutPath, "out", currPath+"/patch"+util.GetNowStr(), "输出目录，默认为当前目录下的[patch+当前时间]目录\n")
	flag.StringVar(&ignoreConfFilePath, "ignore", ".expatch_ignore", "忽略文件的配置文件路径，需要忽略的文件路径按行填写")
	flag.StringVar(&projectType, "type", "idea", "项目类型，可选(idea,eclipse)，默认为idea")
	flag.StringVar(&projectConfigFilePath, "conf", "", "idea模块配置文件路径，当type=idea时候有效")
	flag.StringVar(&author, "author", "", "补丁创建人")
	svnArgs := flag.String("svn", "", "需要对比的两个版本，版本号用英文冒号(:)分隔，带有文件路径则用空格分割，例：'100:95'、'100 .'，为空时则对比当前目录与svn最新版本")
	flag.Parse()

	//解析项目配置文件，得到src、output、WebRoot等路径
	javaProject, err := ide.GetInstance(strings.ToLower(projectType), projectConfigFilePath).AnalysisProjectConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)
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
		fmt.Printf("error: %v\n", err)
		return
	}

	//获取需要忽略的文件，并从需要复制的文件中剔除需要忽略的文件
	ignoreSet := util.ReadFileToSet(ignoreConfFilePath)
	if ignoreSet != nil {
		ignoreSet.Each(func(i interface{}) bool {
			p := i.(string)
			p = util.GetAbsolutePath(p)
			fmt.Printf("ignore file: %s\n", p)
			if copySet.Contains(p) {
				copySet.Remove(p)
			}
			return false
		})
	}

	//获取补丁输出路径
	patchOutPath, err = util.Mkdir(patchOutPath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	//复制差异文件
	err = ide.ConvertPatchDir(copySet, javaProject, patchOutPath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	//生成补丁信息文件，当补丁创建人不为空时默认生成补丁记录信息文件
	if author != "" {
		revision, info := svn.GetSvnInfo(author)
		if info != "" {
			pp, err := util.Mkdir(patchOutPath + "/WEB-INF/patch/")
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			pd, err := os.OpenFile(pp+"/patch."+revision+".properties", os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			defer pd.Close()
			pd.WriteString(info)
		}
	}

	fmt.Println("Done!")
}
