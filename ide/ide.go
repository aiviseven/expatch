package ide

import (
	"fmt"
	"github.com/aiviseven/expatch/util"
	set "github.com/deckarep/golang-set"
	"io/ioutil"
	"strings"
)

type IDE interface {
	// AnalysisProjectConfig 解析项目配置
	AnalysisProjectConfig() (*JavaProject, error)
}

type JavaProject struct {
	ProjectPath    string   //项目根路径
	JavaSrcPaths   []string //src源文件根路径
	WebRootPaths   []string //WebRoot路径
	ClassesOutPath string   //class文件路径
}

func GetInstance(typeName, projectConfigPath string) IDE {
	currPath := util.GetAbsolutePath(util.GetCurrentDirectory())
	if typeName == "idea" {
		return &IdeaProject{JavaProject: JavaProject{ProjectPath: currPath}, ProjectConfigPath: projectConfigPath}
	} else if typeName == "eclipse" {
		return &EclipseProject{JavaProject: JavaProject{ProjectPath: currPath}}
	}

	return &IdeaProject{JavaProject: JavaProject{ProjectPath: currPath}, ProjectConfigPath: projectConfigPath}
}

// ConvertPatchDir 复制源文件到指定位置
func ConvertPatchDir(copySet set.Set, javaProject *JavaProject, patchOutPath string) error {
	copySet.Each(func(i interface{}) bool {
		copyFilePath := i.(string)
		isDir, err := util.IsDir(copyFilePath)
		if err != nil || isDir {
			return false
		}
		fmt.Println(copyFilePath)

		_, javaSrcPath := util.PathContains(javaProject.JavaSrcPaths, copyFilePath)
		_, webRootPath := util.PathContains(javaProject.WebRootPaths, copyFilePath)

		if index := strings.Index(copyFilePath, javaSrcPath); javaSrcPath != "" && index >= 0 {
			//java src文件
			relativeFilePath := copyFilePath[index+len(javaSrcPath):]
			classFileSrcPath := javaProject.ClassesOutPath + relativeFilePath
			classFileDstPath := patchOutPath + "/WEB-INF/classes/" + relativeFilePath

			suffixIndex := strings.LastIndex(relativeFilePath, ".")
			if relativeFilePath[suffixIndex:] == ".java" {
				//java文件需要复制.class文件
				classFileName := relativeFilePath[strings.LastIndex(relativeFilePath, "/")+1 : suffixIndex]

				relativeFilePath = relativeFilePath[:suffixIndex] + ".class"

				if lastSepIndex := strings.LastIndex(classFileSrcPath, "/"); lastSepIndex >= 0 {
					classSrcFiles, _ := ioutil.ReadDir(classFileSrcPath[:lastSepIndex])
					for _, f := range classSrcFiles {
						if !f.IsDir() &&
							(f.Name() == classFileName+".class" ||
								strings.HasPrefix(f.Name(), classFileName+"$")) {
							classFileSrcPath = classFileSrcPath[:lastSepIndex] + "/" + f.Name()
							classFileDstPath = classFileDstPath[:strings.LastIndex(classFileDstPath, "/")] + "/" + f.Name()
							util.CopyFile(classFileSrcPath, classFileDstPath)
						}
					}
				}
			} else {
				//除java文件以外，直接复制文件
				util.CopyFile(classFileSrcPath, classFileDstPath)
			}
		} else if index := strings.Index(copyFilePath, webRootPath); webRootPath != "" && index >= 0 {
			//除了java src文件以外的静态资源文件，直接复制到对应目录
			relativeFilePath := copyFilePath[index+len(webRootPath):]
			dstFilePath := patchOutPath + relativeFilePath
			util.CopyFile(copyFilePath, dstFilePath)
		}
		return false
	})
	return nil
}
