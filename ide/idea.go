package ide

import (
	"encoding/xml"
	"errors"
	"github.com/bsh100220/expatch/util"
	"io/ioutil"
	"strings"
)

//idea配置
//对应后缀为iml的文件
type IdeaProjectConfig struct {
	Module     xml.Name    `xml:"module"`
	Type       string      `xml:"type,attr"`
	Version    string      `xml:"version,attr"`
	Components []Component `xml:"component"`
}
type Component struct {
	Name          string       `xml:"name,attr"`
	LanguageLevel string       `xml:"LANGUAGE_LEVEL,attr"`
	Facets        []Facet      `xml:"facet"`
	Output        Output       `xml:"output"`
	TestOutput    Output       `xml:"output-test"`
	Content       Content      `xml:"content"`
	OrderEntries  []OrderEntry `xml:"orderEntry"`
}
type Facet struct {
	Name          string        `xml:"name,attr"`
	Type          string        `xml:"type,attr"`
	Configuration Configuration `xml:"configuration"`
}
type Configuration struct {
	Descriptors Descriptors `xml:"descriptors"`
	WebRoots    WebRoots    `xml:"webroots"`
	SourceRoots SourceRoots `xml:"sourceRoots"`
}
type WebRoots struct {
	Roots []Root `xml:"root"`
}
type SourceRoots struct {
	Roots []Root `xml:"root"`
}
type Descriptors struct {
	DeploymentDescriptors []DeploymentDescriptor `xml:"deploymentDescriptor"`
}
type Root struct {
	Url      string `xml:"url,attr"`
	Relative string `xml:"relative,attr"`
}
type DeploymentDescriptor struct {
	Name string `xml:"name,attr"`
	Url  string `xml:"url,attr"`
}
type Output struct {
	Url string `xml:"url,attr"`
}
type Content struct {
	Url            string         `xml:"url,attr"`
	SourceFolders  []SourceFolder `xml:"sourceFolder"`
	ExcludeFolders []SourceFolder `xml:"excludeFolders"`
}
type SourceFolder struct {
	Url          string `xml:"url,attr"`
	Type         string `xml:"type,attr"`
	IsTestSource string `xml:"isTestSource,attr"`
}
type OrderEntry struct {
	Type       string  `xml:"type,attr"`
	Name       string  `xml:"name,attr"`
	Level      string  `xml:"level,attr"`
	Scope      string  `xml:"scope,attr"`
	ForTests   string  `xml:"forTests,attr"`
	ModuleName string  `xml:"module-name,attr"`
	Library    Library `xml:"library"`
	JdkName    string  `xml:"jdkName,attr"`
	JdkType    string  `xml:"jdkType,attr"`
}
type Library struct {
	Name           string         `xml:"name,attr"`
	Classes        Classes        `xml:"CLASSES"`
	Javadoc        string         `xml:"JAVADOC"`
	Sources        string         `xml:"SOURCES"`
	JarDirectories []JarDirectory `xml:"jarDirectory"`
}
type Classes struct {
	Root Root `xml:"root"`
}
type JarDirectory struct {
	Url       string `xml:"url,attr"`
	Recursive string `xml:"recursive,attr"`
	Type      string `xml:"type,attr"`
}

type IdeaProject struct {
	JavaProject
	ProjectConfigPath string
}

func (idea *IdeaProject) AnalysisProjectConfig() (*JavaProject, error) {
	err := idea.analysisIdeaProjectConfig()
	if err != nil {
		return nil, err
	}
	return &idea.JavaProject, nil
}

//解析idea模块配置文件
func (idea *IdeaProject) analysisIdeaProjectConfig() error {
	if idea.ProjectConfigPath == "" {
		var ideaConfigFileName string
		configFileCount := 0
		files, _ := ioutil.ReadDir(idea.ProjectPath)
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".iml") {
				ideaConfigFileName = f.Name()
				configFileCount++
			}
		}
		if configFileCount <= 0 {
			return errors.New("项目未包含idea模块配置文件，请先确认项目是否是idea模块。")
		}
		if configFileCount > 1 {
			return errors.New("项目路径中包含多个idea模块配置文件，请指定一个。")
		}
		idea.ProjectConfigPath = idea.ProjectPath + "/" + ideaConfigFileName
	} else {
		idea.ProjectConfigPath = util.GetAbsolutePath(idea.ProjectConfigPath)
	}

	data := util.ReadXmlFile(idea.ProjectConfigPath)
	ideaProjectConfig := IdeaProjectConfig{}
	if data == nil {
		return errors.New("未找到idea模块配置信息")
	}
	err := xml.Unmarshal(data, &ideaProjectConfig)
	if err != nil {
		return err
	}

	for _, component := range ideaProjectConfig.Components {
		if component.Name == "FacetManager" {
		sec:
			for _, facet := range component.Facets {
				if facet.Type == "web" {
					for _, root := range facet.Configuration.WebRoots.Roots {
						idea.WebRootPaths = append(idea.WebRootPaths, idea.ideaPathToReal(root.Url))
						break sec
					}
				}
			}
		} else if component.Name == "NewModuleRootManager" {
			idea.ClassesOutPath = idea.ideaPathToReal(component.Output.Url)
			for _, srcFolder := range component.Content.SourceFolders {
				if srcFolder.IsTestSource == "false" ||
					srcFolder.Type == "java-resource" {
					idea.JavaSrcPaths = append(idea.JavaSrcPaths, idea.ideaPathToReal(srcFolder.Url))
				}
			}
		}
	}

	return nil
}

//替换idea模块配置文件中的地址变量为实际值
func (idea *IdeaProject) ideaPathToReal(path string) (realPath string) {
	return strings.Replace(path, "file://$MODULE_DIR$", idea.ProjectPath, len(path))
}
