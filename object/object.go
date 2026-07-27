package object

import (
	"fmt"

)


type ObjType string

const (
	INTIGER = "intiger"
	FLOAT 	= "float"
	STRING 	= "string"
	BOOL    = "boolean"
	ERROR 	= "error"
	RETURN  = "return_value"
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

type Float struct {
	Value float64
}
func (f *Float) Inspect() string {
	return fmt.Sprintf("%v", f.Value)
}
func (f *Float) Type() ObjType {
	return FLOAT
}


type String struct {
	Value string
}
func (s *String) Inspect() string {
	return fmt.Sprintf("%s", s.Value)
}
func (s *String) Type() ObjType {
	return STRING
}


type Bool struct {
	Value bool
}
func (b *Bool) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}
func (b *Bool) Type() ObjType {
	return BOOL
}

type Error struct {
	Message string
}
func (e *Error) Inspect() string {
	return e.Message
}
func (e *Error) Type() ObjType {
	return ERROR
}

type ReturnValue struct {
	Value Object
}
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
func (r *ReturnValue) Type() ObjType {
	return RETURN
}

