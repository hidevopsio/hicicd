package service

import (
	"github.com/hidevopsio/hicicd/pkg/entity"
		"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/data/bolt"
)

type BoltRepository bolt.Repository

type DictionaryService struct {
	repository BoltRepository
}

func (ds *DictionaryService) Init(repository BoltRepository)  {
	ds.repository = repository
}

func (ds *DictionaryService) Add(dictionary *entity.Dictionary) error {
	ds.repository.Put(dictionary)
	return nil
}

func (ds *DictionaryService) Get(id string) (*entity.Dictionary, error) {
	var dictionary entity.Dictionary
	err := ds.repository.Get(id, &dictionary)
	return &dictionary, err
}


func (ds *DictionaryService) Delete(id string) error  {
	err := ds.repository.Delete([]byte(id))
	log.Debug("Delete Dictionary Service:", err)
	return err
}



