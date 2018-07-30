package service

import (
	"fmt"
	"github.com/hidevopsio/hicicd/pkg/orch/k8s"
	"github.com/hidevopsio/hicicd/pkg/orch/openshift"
	"os"
	"strings"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"encoding/json"
	"github.com/kevholditch/gokong"
	authorization_v1 "github.com/openshift/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/jinzhu/copier"
)

type PipelineService struct {
	entity.Pipeline
}


func (p *PipelineService) CreateSecret(username, password string, isToken bool) (string, error) {
	log.Debug("Pipeline.CreateSecret()")
	if username == "" {
		return "", fmt.Errorf("unkown username")
	}
	// Create secret
	secretName := username + "-secret"
	secret := k8s.NewSecret(secretName, username, password, p.BuildConfigs.Namespace, isToken)
	err := secret.Create()

	return secretName, err
}

func (p *PipelineService) CreateProject() error {
	project, err := openshift.NewProject(p.Namespace, "", "", p.NodeSelector)
	if err != nil {
		log.Error("create namespace err :", err)
		return err
	}
	newProject, err := project.Create()
	if err != nil {
		return err
	}
	//init namespace
	err = p.InitProject()
	if err != nil {
		return err
	}
	log.Debug("create project debug", newProject)
	return nil
}

func (p *PipelineService) CreateRoleBinding(username, metaName, roleRefName string) error {
	roleBinding, err := openshift.NewRoleBinding(metaName, p.Namespace)
	if err != nil {
		return err
	}
	r, err := roleBinding.Get()
	if err != nil {
		role := &authorization_v1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metaName,
				Namespace: p.Namespace,
			},
			RoleRef: corev1.ObjectReference{
				Name: roleRefName,
			},
			Subjects: []corev1.ObjectReference{
				{
					Name: username,
					Kind: "User",
				},
			},
		}
		_, err = roleBinding.Create(role)
		return err
	}
	for _, value := range r.Subjects {
		if value.Name == username {
			return nil
		}
	}
	reference := corev1.ObjectReference{
		Name: username,
		Kind: "User",
	}
	r.Subjects = append(r.Subjects, reference)
	r.UserNames = append(r.UserNames, username)
	_, err = roleBinding.Create(r)
	return err
}

func (p *PipelineService) Build(secret string, completedHandler func() error) error {
	log.Debug("Pipeline.Build()")

	if !p.BuildConfigs.Enable {
		return completedHandler()
	}

	scmUrl := p.CombineScmUrl()
	buildConfig, err := openshift.NewBuildConfig(p.BuildConfigs.Namespace, p.App, scmUrl, p.BuildConfigs.Branch, secret, p.Version, p.BuildConfigs.ImageStream, p.BuildConfigs.Rebuild)
	if err != nil {
		return err
	}
	_, err = buildConfig.Create()
	if err != nil {
		return err
	}
	// Build image stream
	build, err := buildConfig.Build(p.BuildConfigs.Env)

	if err != nil {
		log.Error("buildConfig.Build(p.BuildConfigs.Env)", err)
		return err
	}

	err = buildConfig.Watch(build, completedHandler)
	log.Info("buildConfig.Watch", err)
	return err
}

func (p *PipelineService) CombineScmUrl() string {
	scmUrl := p.Scm.Url + "/" + p.Project + "/" + p.App + "." + p.Scm.Type
	return scmUrl
}

func (p *PipelineService) RunUnitTest() error {
	log.Debug("Pipeline.RunUnitTest()")
	return nil
}

func (p *PipelineService) RunIntegrationTest() error {
	log.Debug("Pipeline.RunIntegrationTest()")
	return nil
}

func (p *PipelineService) Analysis() error {
	log.Debug("Pipeline.Analysis()")
	return nil
}

func (p *PipelineService) CreateDeploymentConfig(force bool, injectSidecar func(in interface{}) (interface{}, error)) error {
	log.Debug("Pipeline.CreateDeploymentConfig()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace, p.Version)
	if err != nil {
		return err
	}
	var l map[string]string
	labels, _ := json.Marshal(p.DeploymentConfigs.Labels)
	err = json.Unmarshal(labels, &l)
	err = dc.Create(&p.DeploymentConfigs.Env, l, &p.Ports, p.DeploymentConfigs.Replicas, force, p.DeploymentConfigs.HealthEndPoint, p.NodeSelector, injectSidecar)
	if err != nil {
		log.Error("dc.Create ", err)
		return err
	}

	return nil
}

func (p *PipelineService) InstantiateDeploymentConfig() error {
	log.Debug("Pipeline.Deploy()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace, p.Version)
	if err != nil {
		return err
	}

	d, err := dc.Instantiate()
	log.Debug(d.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p *PipelineService) CreateKongGateway(upstreamUrl string) error {
	log.Debug("Pipeline.CreateKongGateway()")
	uris := p.GatewayConfigs.Uri
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", p.Profile, -1)
	hosts := strings.Split(host, ",")
	apiRequest := &gokong.ApiRequest{
		Name:                   p.App + "-" + p.Project,
		Hosts:                  hosts,
		Uris:                   []string{uris},
		UpstreamUrl:            "http://" + upstreamUrl,
		StripUri:               true,
		PreserveHost:           true,
		Retries:                "3",
		UpstreamConnectTimeout: 1000,
		UpstreamSendTimeout:    2000,
		UpstreamReadTimeout:    3000,
		HttpsOnly:              false,
		HttpIfTerminated:       true,
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", p.Profile, -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	_, err := gokong.NewClient(config).Apis().Create(apiRequest)
	return err
}

func (p *PipelineService) CreateService() error {
	log.Debug("Pipeline.CreateService()")

	// new dc instance
	svc := k8s.NewService(p.App, p.Namespace)

	err := svc.Create(&p.Ports)
	log.Debug("Pipeline.CreateService svc.Create result: ", err)
	return err
}

func (p *PipelineService) CreateRoute() (string, error) {
	log.Debug("Pipeline.CreateRoute()")
	upstreamUrl := ""
	route, err := openshift.NewRoute(p.App, p.Namespace)
	if err != nil {
		return upstreamUrl, err
	}

	upstreamUrl, err = route.Create(8080)
	return upstreamUrl, err
}

func (p *PipelineService) GetImageStreamTag() error {
	log.Debug("pipeline get image stream tag :")
	ist, err := openshift.NewImageStreamTags(p.App, p.Version, p.BuildConfigs.Project+"-"+p.BuildConfigs.TagFrom)
	if err != nil {
		log.Error("Pipeline.CreateImageStreamTag.NewImageStreamTags", err)
		return err
	}
	image, err := ist.Get()
	if err != nil {
		log.Error("images not found :", err)
		return err
	}
	log.Debug("pipeline get image stream tag name:", image.Name)
	return nil
}

func (p *PipelineService) CreateImageStreamTag() error {
	log.Debug("Pipeline.CreateImageStreamTag")
	ist, err := openshift.NewImageStreamTags(p.App, p.Version, p.Namespace)
	if err != nil {
		log.Error("Pipeline.CreateImageStreamTag.NewImageStreamTags", err)
		return err
	}
	_, err = ist.Create(p.BuildConfigs.Project + "-" + p.BuildConfigs.TagFrom)
	return err
}

func (p *PipelineService) InitProject() error {
	//init Groups
	roleBinding := &openshift.RoleBinding{
		Namespace: p.Namespace,
	}
	//init image builders
	err := roleBinding.InitImageBuilders()
	if err != nil {
		return err
	}
	//init image pullers
	err = roleBinding.InitImagePullers()
	if err != nil {
		return err
	}
	//init system deployers
	err = roleBinding.InitSystemDeployers()
	if err != nil {
		return err
	}
	return nil
}

func (p *PipelineService) Deploy() error {
	log.Info("p.Deploy()")
	if p.DeploymentConfigs.Enable {

		// create dc - deployment config
		err := p.CreateDeploymentConfig(p.DeploymentConfigs.ForceUpdate, func(in interface{}) (interface{}, error) {
			if !p.IstioConfigs.Enable {
				return in, nil
			}
			injector := &istio.Injector{}
			copier.Copy(injector, p.IstioConfigs)
			return injector.Inject(in)
		})
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("failed on CreateDeploymentConfig! %s", err.Error())
		}
		rc := k8s.NewReplicationController(p.App, p.Namespace)
		// rc.Watch(message, handler)
		err = rc.Watch(func() error {
			log.Debug("Completed!")
			return nil
		})
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("failed on watch rc! %s", err.Error())
		}
	}

	// create service
	err := p.CreateService()
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed on CreateService! %s", err.Error())
	}

	// create route
	upstreamUrl, err := p.CreateRoute()
	log.Debug("CreateRoute get upstream url :", upstreamUrl)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed on CreateRoute! %s", err.Error())
	}

	//create kong-gateway
	if p.GatewayConfigs.Enable {
		err = p.CreateKongGateway(p.GatewayConfigs.UpstreamUrl)
	}
	return err
}

func (p *PipelineService) Run(username, password, token string, uid int, isToken bool) error {
	log.Debug("Pipeline.Run()")
	// TODO: check if the same app in the same namespace is already in running status.
	permission := &auth.Permission{}
	metaName, roleRefName, accessLevelValue, err := permission.Get(p.Scm.Url, token, p.App, p.BuildConfigs.Project, uid)
	// TODO: accessLevelValue permission 30
	if err != nil || accessLevelValue < auth.DeveloperPermissions {
		return err
	}
	// TODO: first, let's check if namespace is exist or not
	//create  namespace
	err = p.CreateProject()
	if err != nil {
		log.Error("Pipeline run new namespace err:", err)
		return err
	}
	//Authorize the user
	err = p.CreateRoleBinding(username, metaName, roleRefName)
	if err != nil {
		log.Error("Pipeline run Create RoleBinding err :", err)
		return err
	}
	if p.BuildConfigs.TagFrom == p.Profile && p.BuildConfigs.Project == p.Project {

		// create secret for building image
		secret, err := p.CreateSecret(username, password, isToken)
		if err != nil {
			return fmt.Errorf("failed on CreateSecret! %s", err.Error())
		}

		// build image
		err = p.Build(secret, func() error {
			return p.Deploy()
		})
	} else {
		//TODO check if tag from is exist or not
		err = p.GetImageStreamTag()
		if err != nil {
			return nil
		}
		err = p.CreateImageStreamTag()
		if err != nil {
			log.Error("Pipeline.Run.CreateImageStreamTag error", err)
			return err
		}
		p.Deploy()
	}

	if err != nil {
		return fmt.Errorf("failed on Build! %s", err.Error())
	}

	// interact with developer
	// deploy to test ? yes/no

	// finally, all steps are done well, let tell the client ...
	return nil
}

