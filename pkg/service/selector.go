package service

import (
	"github.com/hidevopsio/hicicd/pkg/entity"
	"encoding/json"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type SelectorService struct {
	repository BoltRepository
}

func (ss *SelectorService) Init(repository BoltRepository)  {
	ss.repository = repository
}

func (ss *SelectorService) Add(selectors *entity.Selector) error {
	s, err := json.Marshal(selectors)
	if err == nil {
		ss.repository.Put([]byte(selectors.Id), s)
	}
	return nil
}

func (ss *SelectorService) Get(id string) (*entity.Selector, error) {
	var selectors = entity.Selector{}
	err := ss.repository.Get(id, &selectors)
	if err != nil {
		return nil, err
	}

	return &selectors, err
}


func (ss *SelectorService) Delete(id string) error  {
	err := ss.repository.Delete([]byte(id))
	log.Debug("Delete Selector Service:", err)
	return err
}





