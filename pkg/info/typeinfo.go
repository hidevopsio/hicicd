package info

import (
	"encoding/base64"
	"github.com/ghodss/yaml"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
	"github.com/kataras/iris/core/errors"
	"encoding/xml"
)

const (
	JGroups = "axon-jgroups"
	CQRS = "java-cqrs"
	WAR  = "war"
	JAVA = "java"
	JAVA_POM = "java-pom"
)


type TypeInfo struct {
	ClusterName string `json:"cluster_name"`
	AppType     string `json:"app_type"`
}

//app.cluster-name
type Application struct {
	App    App    `yaml:"app"`
	Spring Spring `yaml:"spring"`
}

type App struct {
	Project     string `yaml:"project"`
	ClusterName string `yaml:"cluster-name"`
}

type Spring struct {
	Profiles Profiles `yaml:"profiles"`
}

type Profiles struct {
	Include []string `yaml:"include,flow"`
}

type AppType struct {
	Name       string `json:"name"`
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	configFile string `json:"config_file"`
}

var AppTypeLUT = []AppType{
	{Name: "java", FileName: "pom.xml", FileType: "blob", configFile: "src/main/resources/application.yml"},
	{Name: "nodejs-dist", FileName: "dist", FileType: "tree",},
	{Name: "nodejs", FileName: "package.json", FileType: "blob"},
}

type ResourceString struct {
	XMLName      xml.Name `xml:"project"`
	ModelVersion string   `xml:"modelVersion"`
	GroupId      string   `xml:"groupId"`
	ArtifactId   string   `xml:"artifactId"`
	Version      string   `xml:"version"`
	Packaging    string   `xml:"packaging"`
}

func (t *TypeInfo) RepositoryType(url, token, ref string, id int) error {
	scmFactory := new(factories.ScmFactory)
	r, err := scmFactory.NewResitory(factories.GitlabScmType)
	if err != nil {
		log.Error("type info repository type:", err)
		return err
	}
	treeNodes, err := r.ListTree(url, token, ref, id)
	for _, appType := range AppTypeLUT {
		for _, treeNode := range treeNodes {
			if appType.FileType == treeNode.Type && appType.FileName == treeNode.Name {
				log.Info(appType)
				t.AppType = appType.Name
				if appType.Name == JAVA {
					//TODO if pom.xml <packaging>jar</packaging>  packagin is jar TypeInfo.typeapp=java   typeapp=java-war
					pomContext, err := r.GetRepository(url, token, "pom.xml", ref, id)
					if err != nil {
						log.Error("type app repository pom context:", err)
						return nil
					}
					err = t.ParsePom(pomContext)
					if err != nil {
						log.Error("pom.xml parse err:", err)
						return err
					}
					//TODO parse application.yaml  culsterName
					context, err := r.GetRepository(url, token, appType.configFile, ref, id)
					log.Debug("type info GetRepository:", err)
					err = t.Parse(context)
					return err
				}
				return nil
			}
		}
	}
	return errors.New("not fount file name ")
}

func (t *TypeInfo) Parse(context string) error {
	application, err := ParseYaml(context)
	if err != nil {
		log.Error("type info RepositoryType parse err: ", err)
		return err
	}
	for _, include := range application.Spring.Profiles.Include {
		if include == JGroups {
			t.AppType = CQRS
			t.ClusterName = application.App.ClusterName
			return nil
		}
	}
	return err
}

func ParseYaml(content string) (*Application, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}
	application := &Application{}
	err = yaml.Unmarshal([]byte(data), application)
	log.Info("application:{}", application)
	return application, nil
}

func (t *TypeInfo) ParsePom(content string) error {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	resource := &ResourceString{}
	err = xml.Unmarshal([]byte(data), resource)
	log.Info("pom.xml parse resource:", resource.Packaging)
	if resource.Packaging == WAR || err != nil {
		t.AppType = "java-war"
	} else if resource.Packaging == JAVA_POM{
		t.AppType = JAVA_POM
	}
	return err
}
