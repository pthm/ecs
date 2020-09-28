package ecs

import (
	"sync"
	"sync/atomic"
)

type Entity interface {
	GetID() uint64
	GetComponents() []Component
	AddComponents(c ...Component)
	RemoveComponents(c ...Component)
}

type BasicEntity struct {
	mux        *sync.RWMutex
	ID         uint64
	Components []Component
}

func (be BasicEntity) GetID() uint64 {
	be.mux.RLock()
	defer be.mux.RUnlock()
	return be.ID
}

func (be BasicEntity) GetComponents() []Component {
	be.mux.RLock()
	defer be.mux.RUnlock()
	return be.Components
}

func (be *BasicEntity) AddComponents(cs ...Component) {
	be.mux.Lock()
	defer be.mux.Unlock()
	for _, c := range cs {
		be.Components = append(be.Components, c)
	}
}

func (be *BasicEntity) RemoveComponents(cs ...Component) {
	be.mux.Lock()
	defer be.mux.Unlock()
	for _, c := range cs {
		var delete int = -1
		for index, component := range be.Components {
			if c == component {
				delete = index
				break
			}
		}
		if delete >= 0 {
			be.Components = append(be.Components[:delete], be.Components[delete+1:]...)
		}
	}
}

func NewBasicEntity() BasicEntity {
	return BasicEntity{ID: atomic.AddUint64(&IDInc, 1), mux: &sync.RWMutex{}}
}
