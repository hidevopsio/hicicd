package service

import (
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hicicd/pkg/protobuf"
	"time"
	"golang.org/x/net/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/jinzhu/copier"
)

type Profile struct {
	Id           string      `json:"id"`            //主键
	Name         string      `json:"name"`          //名称
	ProfileId    string      `json:"profile_id"`    //profileId
	Namespace    string      `json:"namespace"`     //主键
	Host         string      `json:"host"`          //kong 的host example  http://devcloud.vpclub.cn http://devgw.vpclub.cn
	Rebuild      bool        `json:"re_build"`      // 是否重新获取最新镜像 true or false
	Deploy       bool        `json:"deploy"`        // 是否 重新部署 true or false
	Gateway      bool        `json:"gateway"`       // 是否配置网关 true or false
	BuildEnable  bool        `json:"build_enable"`  //是否重新build
	RemoteEnable bool        `json:"remote_enable"` //是否远程部署
	TagEnable    bool        `json:"tag_enable"`
	ImageTags    []*ImageTag `json:"image_string_tags"`
	ForceUpdate  bool        `json:"force_update"`
	TagFrom      string      `json:"tag_from"`
	ScmRef       string      `json:"scm_ref"`      //分支
	IstioEnable  bool        `json:"istio_enable"` //istio
	CreatedTime  int64       `json:"created_time"` //创建时间
	CreatedBy    int64       `json:"created_by"`   //创建人
	UpdatedTime  int64       `json:"updated_time"` //修改时间
	UpdatedBy    int64       `json:"updated_by"`   //修改人
	Deleted      int8        `json:"deleted"`      //是否删除  删除 1 没有删除 2
}

type ImageTag struct {
	Id         string `json:"id"`
	ProfileId  string `json:"profile_id"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Text       string `json:"text"`
	Value      string `json:"value"`
	Name       string `json:"name"`
}

type ProfileService struct {
	profileClient protobuf.ProfileServiceClient
}

func (p *ProfileService) Init(profileClient protobuf.ProfileServiceClient) {
	p.profileClient = profileClient
}

func init() {
	grpc.RegisterClient("hiadmin-client", protobuf.NewProfileServiceClient)
}

func (p *ProfileService) Get(id string) (*Profile, error) {
	request := &protobuf.ProfileRequest{}
	request.Id = id
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	response, err := p.profileClient.Get(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	profile := &Profile{}
	copier.Copy(profile, response.Data)
	return profile, nil
}
