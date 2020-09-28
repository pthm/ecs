package ecs

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"
)

type World struct {
	mux      *sync.RWMutex
	Systems  []SystemRegistration
	Entities []Entity
	nextid   uint64
}

func NewWorld() *World {
	return &World{
		mux:      &sync.RWMutex{},
		Systems:  make([]SystemRegistration, 0),
		Entities: make([]Entity, 0),
	}
}

func (w *World) AddSystem(s System, sc ...Component) {
	w.mux.Lock()
	defer w.mux.Unlock()
	sr := SystemRegistration{
		System:     s,
		Components: sc,
	}
	w.Systems = append(w.Systems, sr)
}

func (w *World) AddEntity(e Entity) {
	w.mux.Lock()
	defer w.mux.Unlock()
	w.Entities = append(w.Entities, e)
}

func (w *World) RemoveEntity(e Entity) {
	w.mux.Lock()
	defer w.mux.Unlock()
	var delete int = -1
	for index, entity := range w.Entities {
		if e == entity {
			delete = index
			break
		}
	}
	if delete >= 0 {
		w.Entities = append(w.Entities[:delete], w.Entities[delete+1:]...)
	}
}

func (w *World) GetEntitiesForSystemRegistration(sr SystemRegistration) []Entity {
	w.mux.RLock()
	defer w.mux.RUnlock()
	if len(sr.Components) == 1 {
		return w.GetEntitiesWithComponent(sr.Components[0])
	} else {
		return w.GetEntitiesWithAllComponents(sr.Components...)
	}
	panic("System not present in world")
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (w *World) Update(dt float64) {
	var wg sync.WaitGroup
	var singleWg sync.WaitGroup
	priorityGroups := make(map[int][]SystemRegistration)

	w.mux.RLock()
	defer w.mux.RUnlock()

	for _, sr := range w.Systems {
		priorityGroups[sr.System.Priority()] = append(priorityGroups[sr.System.Priority()], sr)
	}

	keys := make([]int, 0)
	for k, _ := range priorityGroups {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		srs := priorityGroups[k]
		wg.Add(1)
		for _, sr := range srs {
			singleWg.Add(1)
			entities := w.GetEntitiesForSystemRegistration(sr)
			go func(dt float64, entities []Entity, wg *sync.WaitGroup) {
				start := time.Now()
				sr.System.Update(dt, entities)
				t := time.Since(start).Microseconds()
				if t > 1200 {
					fmt.Printf("updated %v in %vus\n", getType(sr.System), t)
				}
				singleWg.Done()
			}(dt, entities, &wg)
		}
		singleWg.Wait()
		wg.Done()
	}

	wg.Wait()
}

func (w *World) GetEntitiesWithComponent(component Component) []Entity {
	w.mux.RLock()
	defer w.mux.RUnlock()
	results := make([]Entity, 0)
	// Iterate over all entities in world
	for _, entity := range w.Entities {
		entityComponents := entity.GetComponents()
		if len(entityComponents) == 0 {
			continue
		}
		// Iterate components of entity to find component
		for _, entityComponent := range entityComponents {
			if entityComponent.GetName() == component.GetName() {
				results = append(results, entity)
				break
			}
		}
	}
	return results
}

func (w *World) GetEntitiesWithAllComponents(componentsQuery ...Component) []Entity {
	w.mux.RLock()
	defer w.mux.RUnlock()
	results := make([]Entity, 0)
	searchLen := len(componentsQuery)
	// Iterate over all entities in world
	for _, entity := range w.Entities {
		entityComponents := entity.GetComponents()
		if len(entityComponents) == 0 {
			continue
		}
		// Test for existance of all components in query
		found := 0
		for _, testComponent := range componentsQuery {
			// Iterate components of entity to find component in query
			for _, entityComponent := range entityComponents {
				if entityComponent.GetName() == testComponent.GetName() {
					found++
				}
			}
		}
		if found == searchLen {
			results = append(results, entity)
		}
	}
	return results
}

func (w *World) GetEntityByID(id uint64) (Entity, error) {
	w.mux.RLock()
	defer w.mux.RUnlock()
	for _, entity := range w.Entities {
		if entity.GetID() == id {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("entity %v does not exist", id)
}
