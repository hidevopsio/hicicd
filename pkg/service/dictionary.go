package service

import (
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hicicd/pkg/admin"
	"encoding/json"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type DictionaryService struct {
	Repository db.KVRepository `component:"repository" dataSourceType:"bolt"`
}

func (ds *DictionaryService) CreateDictionaryService(dictionary *admin.Dictionary) error {
	d, err := json.Marshal(dictionary)
	if err == nil {
		ds.Repository.Put([]byte("dictionary"), []byte(dictionary.Id), d)
	}
	return nil
}

func (ds *DictionaryService) GetDictionaryService(id string) (*admin.Dictionary, error) {
	d, err := ds.Repository.Get([]byte("dictionary"), []byte(id))
	if err != nil {
		return nil, err
	}
	var dictionary = &admin.Dictionary{}
	err = json.Unmarshal(d, dictionary)
	return dictionary, err
}


func (ds *DictionaryService) DeleteDictionaryService(id string) error  {
	err := ds.Repository.Delete([]byte("dictionary"), []byte(id))
	log.Debug("Delete Dictionary Service:", err)
	return err
}



