package ide

import (
	"encoding/xml"
	"errors"
	"github.com/bsh100220/expatch/util"
)

//普通eclipse配置
//对应.classpath文件
type EclipseNormalProjectConfig struct {
	Classpath        xml.Name         `xml:"classpath"`
	ClasspathEntries []ClasspathEntry `xml:"classpathentry"`
}
type ClasspathEntry struct {
	Kind string `xml:"kind,attr"`
	Path string `xml:"path,attr"`
}

//eclipse Web配置
//对应.settings/org.eclipse.wst.common.component文件
type EclipseWebProjectConfig struct {
	ProjectModuleId string    `xml:"id,attr"`
	ProjectVersion  string    `xml:"project-version,attr"`
	WebModule       WebModule `xml:"wb-module"`
}
type WebModule struct {
	DeployName       string            `xml:"deploy-name,attr"`
	WebResources     []WebResource     `xml:"wb-resource"`
	DependentModules []DependentModule `xml:"dependent-module"`
	Properties       []Property        `xml:"property"`
}
type WebResource struct {
	DeployPath string `xml:"deploy-path,attr"`
	SourcePath string `xml:"source-path,attr"`
	Tag        string `xml:"tag,attr"`
}
type DependentModule struct {
	ArchiveName    string `xml:"archiveName,attr"`
	DeployPath     string `xml:"deploy-path,attr"`
	Handle         string `xml:"handle,attr"`
	DependencyType string `xml:"dependency-type"`
}
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type EclipseProject struct {
	JavaProject
}

func (e *EclipseProject) AnalysisProjectConfig() (*JavaProject, error) {
	err := e.analysisEclipseNormalProjectConfig()
	if err != nil {
		return nil, err
	}
	err = e.analysisEclipseWebProjectConfig()
	if err != nil {
		return nil, err
	}
	return &e.JavaProject, nil
}

//解析eclipse普通工程配置文件
func (e *EclipseProject) analysisEclipseNormalProjectConfig() error {
	projectConfigPath := e.ProjectPath + "/.classpath"
	data := util.ReadXmlFile(projectConfigPath)
	eclipseNormalConfig := EclipseNormalProjectConfig{}
	if data == nil {
		return errors.New("eclipse project config is blank")
	}
	err := xml.Unmarshal(data, &eclipseNormalConfig)
	if err != nil {
		return err
	}
	for _, classpathEntry := range eclipseNormalConfig.ClasspathEntries {
		if classpathEntry.Kind == "src" {
			e.JavaSrcPaths = append(e.JavaSrcPaths, e.ProjectPath+"/"+classpathEntry.Path)
		} else if classpathEntry.Kind == "output" {
			e.ClassesOutPath = e.ProjectPath + "/" + classpathEntry.Path
		}
	}
	return nil
}

//解析eclipseWeb工程的配置文件
func (e *EclipseProject) analysisEclipseWebProjectConfig() error {
	webConfigPath := e.ProjectPath + "/.settings/org.eclipse.wst.common.component"
	isWebProject, _ := util.PathExists(webConfigPath)
	if isWebProject {
		data := util.ReadXmlFile(webConfigPath)
		eclipseWebConfig := EclipseWebProjectConfig{}
		if data == nil {
			return errors.New("未找到eclipse工程配置信息")
		}
		err := xml.Unmarshal(data, &eclipseWebConfig)
		if err != nil {
			return err
		}
		for _, webResource := range eclipseWebConfig.WebModule.WebResources {
			if webResource.DeployPath == "/" {
				e.WebRootPaths = append(e.WebRootPaths, e.ProjectPath+webResource.SourcePath)
			}
		}
	}
	return nil
}
