package service

import (
	"github.com/hidevopsio/hicicd/pkg/protobuf"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"golang.org/x/net/context"
	"github.com/hidevopsio/hiboot/pkg/utils/copier"
	"time"
	"github.com/prometheus/common/log"
	"github.com/hidevopsio/hioak/pkg"
	"net/http"
)

type RemoteDeploymentConfigs struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Token          string `json:"token"`
	Namespace      string `json:"namespace"`
	Project        string `json:"project"`
	Version        string `json:"version"`
	Profile        string `json:"profile"`
	RegistryUrl    string `json:"registry_url"`
	RemoteUrl      string `json:"remote_url"`
	ApplicationId  string `json:"application_id"`
	DockerRegistry string `json:"docker_registry"` //svc 地址
	Secret         string `json:"secret"`
	CallBackUrl    string `json:"call_back_url"`
	CreatedTime    int64  `json:"created_time"`
	CreatedBy      string `json:"created_by"`
	UpdatedTime    int64  `json:"updated_time"`
	UpdatedBy      string `json:"updated_by"`
}

type RemoteDeploymentConfigsService struct {
	remoteDeploymentClient protobuf.RemoteDeploymentServiceClient
	deployClient           protobuf.DeployServiceClient
}

func (r *RemoteDeploymentConfigsService) Init(remoteDeploymentClient protobuf.RemoteDeploymentServiceClient, deployClient protobuf.DeployServiceClient) {
	r.remoteDeploymentClient = remoteDeploymentClient
	r.deployClient = deployClient
}

func init() {
	grpc.RegisterClient("hiadmin-client", protobuf.NewRemoteDeploymentServiceClient)
	grpc.RegisterClient("hiagent-client", protobuf.NewDeployServiceClient)
}

func (r *RemoteDeploymentConfigsService) InitRemote(id, namespace, profile, app, version string) (*RemoteDeploymentConfigs, error) {
	remote := &RemoteDeploymentConfigs{}
	request := &protobuf.RemoteDeploymentRequest{
		Id: id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	response, err := r.remoteDeploymentClient.Get(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	copier.Copy(remote, response.Data)
	cli := orch.GetClientInstance()
	token := cli.Config().BearerToken
	remote.Token = token
	remote.Profile = profile
	remote.Namespace = namespace
	remote.Version = version
	remote.Project = app
	return remote, nil
}

func (r *RemoteDeploymentConfigsService) Run(remoteDeploy *RemoteDeploymentConfigs) error {
	request := &protobuf.DeployRequest{}
	copier.Copy(request, remoteDeploy)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	response, err := r.deployClient.Run(ctx, request)
	if err != nil && response.Code == http.StatusOK {
		log.Error(err)
		return err
	}
	return nil
}
