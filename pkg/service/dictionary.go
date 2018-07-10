package service

import (
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hicicd/pkg/admin"
	"encoding/json"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type DictionaryService struct {
	Repository db.KVRepository `inject:"repository,dataSourceType=bolt,namespace=dictionary"`
}

func (ds *DictionaryService) Add(dictionary *admin.Dictionary) error {
	d, err := json.Marshal(dictionary)
	if err == nil {
		ds.Repository.Put([]byte(dictionary.Id), d)
	}
	return nil
}

func (ds *DictionaryService) Get(id string) (*admin.Dictionary, error) {
	d, err := ds.Repository.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	var dictionary = &admin.Dictionary{}
	err = json.Unmarshal(d, dictionary)
	return dictionary, err
}


func (ds *DictionaryService) Delete(id string) error  {
	err := ds.Repository.Delete([]byte(id))
	log.Debug("Delete Dictionary Service:", err)
	return err
}



