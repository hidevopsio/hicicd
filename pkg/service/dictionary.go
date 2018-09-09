package service

import (
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hicicd/pkg/protobuf"
	"time"
	"golang.org/x/net/context"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/copier"
)

type Dictionary struct {
	Id          string `json:"id"`
	AppId       int64  `json:"app_id"`
	Name        string `json:"name"`
	Type        int32    `json:"type"`
	Value       string `json:"value"`
	CreatedBy   int64  `json:"created_by"`
	UpdatedBy   int64  `json:"updated_by"`
	CreatedTime int64  `json:"created_time"`
	UpdateTime  int64  `json:"update_time"`
	Deleted     int8   `json:"deleted"`
}


type DictionaryService struct {
	dictionaryClient protobuf.DictionaryServiceClient
}

func (d *DictionaryService) Init(dictionaryClient protobuf.DictionaryServiceClient) {
	d.dictionaryClient = dictionaryClient
}

func init() {
	grpc.RegisterClient("hiadmin-client", protobuf.NewDictionaryServiceClient)
}

func (d *DictionaryService) GetType(id int32) (*[]Dictionary, error) {
	request := &protobuf.DictionaryRequest{}
	request.Type = id
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	response, err := d.dictionaryClient.Get(ctx, request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	dictionary := &[]Dictionary{}
	copier.Copy(dictionary, response.Data)
	return dictionary, nil
}