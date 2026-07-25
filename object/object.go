package object

import "fmt"


type ObjType string

const (
	INTIGER = "intiger"
)

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		outer: outer,
		store: map[string]Object{},
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		store: map[string]Object{},
	}
}

type Environment struct {
	outer *Environment
	store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil{
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) {
	e.store[name] = obj
} 

type Object interface {
	Inspect() string
	Type() ObjType
}

type Intiger struct {
	Value int64
}
func (i *Intiger) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Intiger) Type() ObjType {
	return INTIGER
}
