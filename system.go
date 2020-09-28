package ecs

type System interface {
	Priority() int
	Update(dt float64, entities []Entity)
}

type SystemRegistration struct {
	System     System
	Components []Component
}
