package service

import (
	"github.com/hidevopsio/hiboot/pkg/starter/db"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"encoding/json"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type SelectorService struct {
	Selector db.KVRepository `inject:"selector,dataSourceType=bolt,namespace=selector"`
}



func (ss *SelectorService) Add(selectors *entity.Selector) error {
	s, err := json.Marshal(selectors)
	if err == nil {
		ss.Selector.Put([]byte(selectors.Id), s)
	}
	return nil
}

func (ss *SelectorService) Get(id string) (*entity.Selector, error) {
	s, err := ss.Selector.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	var selectors = &entity.Selector{}
	err = json.Unmarshal(s, selectors)
	return selectors, err
}


func (ss *SelectorService) Delete(id string) error  {
	err := ss.Selector.Delete([]byte(id))
	log.Debug("Delete Selector Service:", err)
	return err
}





