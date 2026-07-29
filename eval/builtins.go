package eval

import (
	"fmt"

	"github.com/myselfBZ/bzscript/object"
)

var builtIns = map[string]*object.BuiltIn{
	"print":{
		Fn: PrintBuiltin,
	},
	"len":{
		Fn: LenBuiltin,
	},
	"float":{
		Fn: FloatBuiltin,
	},
	"int":{
		Fn: IntBuiltin,
	},
}


func IntBuiltin(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("int() requires 1 argument, got %d", len(objs))
	}
	switch t := objs[0].(type) {
	case *object.Intiger:
		return &object.Intiger{Value: t.Value}
	case *object.Float:
		return &object.Intiger{Value: int64(t.Value)}
	default:
		return newError("%s cannot be converted to an integer", objs[0].Type())
	}
}

func FloatBuiltin(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("float() requires 1 argument, got %d", len(objs))
	}
	switch t := objs[0].(type) {
	case *object.Intiger:
		return &object.Float{Value: float64(t.Value)}
	case *object.Float:
		return &object.Float{Value: t.Value}
	default:
		return newError("%s cannot be converted to float", objs[0].Type())
	}
}

func LenBuiltin(objs ...object.Object) object.Object {
	if len(objs) != 1 {
		return newError("len() requires 1 argument, got %d", len(objs))
	}

	switch t := objs[0].(type) {
	case *object.String:
		return &object.Intiger{Value: int64(len(t.Value))}
	case *object.Array:
		return &object.Intiger{Value: int64(len(t.Elements))}
	default:
		return newError("length of type %s cannot be measured", objs[0].Type())
	}
}

func PrintBuiltin(objs ...object.Object) object.Object {
	for _, obj := range objs {
		if obj == Nothing {
			continue
		}
		switch t := obj.(type) {
		case *object.Intiger:
			fmt.Print(t.Value)
		case *object.Bool:
			fmt.Print(t.Value)
		case *object.String:
			fmt.Print(t.Value)
		default:
			fmt.Printf("non-printable object: %s", obj.Type())
		}
		fmt.Print(" ")
	}
	fmt.Println()
	return Nothing
}
