package ecs

var IDInc uint64

type Component interface {
	GetName() string
}
