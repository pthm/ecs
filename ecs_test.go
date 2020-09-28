package ecs

import (
	"fmt"
	"testing"
)

type Player struct {
	BasicEntity
	Name string
}

type TalkComponent struct {
	text string
}

type TalkerSystem struct {
}

func (ts TalkerSystem) Update(dt float64, entities []Entity) {
	for _, talker := range entities {
		fmt.Println(talker)
		talkerComponents := (talker).GetComponents()
		for _, component := range talkerComponents {
			switch typedComponent := component.(type) {
			case *TalkComponent:
				fmt.Println(typedComponent.text)
			}
		}
	}
}

func Test_ECS(t *testing.T) {
	w := NewWorld()
	p := &Player{
		BasicEntity: NewBasicEntity(),
		Name:        "PT",
	}

	p2 := &Player{
		BasicEntity: NewBasicEntity(),
		Name:        "Stef",
	}

	tc := &TalkComponent{
		text: "Hello",
	}

	w.AddEntity(p)
	w.AddEntity(p2)

	p.AddComponents(tc)

	ts := &TalkerSystem{}

	w.AddSystem(ts, &TalkComponent{})

	w.Update(0)

	p.RemoveComponents(tc)

	w.Update(1)
}
