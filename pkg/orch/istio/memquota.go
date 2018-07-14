package istio

import (
	"istio.io/istio/pilot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/log"
	mccpb "istio.io/api/mixer/v1/config/client"
)

const (
	MemquotaType    = "quota-spec"
	MemquotaVersion = "v1alpha2"
	MemquotaGroup   = "config.istio.io"
)

type Memquota struct {
	Client
}

func (memquota *Memquota) getConfig() (model.Config, error) {
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        MemquotaType,
			Version:     MemquotaVersion,
			Group:       MemquotaGroup,
			Name:        memquota.Name,
			Namespace:   memquota.Namespace,
			Labels:      memquota.Labels,
			Annotations: memquota.Annotations,
		},
		Spec: &mccpb.QuotaSpec{
			Rules: []*mccpb.QuotaRule{
				{
					Match: []*mccpb.AttributeMatch{{
						Clause: map[string]*mccpb.StringMatch{"": {
						},
						},
					},
					}},
			},
		},
	}
	return config, nil
}

func (memquota *Memquota) Create() (string, error) {
	log.Debug("create rule :", memquota)
	config, err := memquota.getConfig()
	con, exists := memquota.Get(MemquotaType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := memquota.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := memquota.Crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}
