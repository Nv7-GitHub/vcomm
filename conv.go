package vcomm

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Nv7-Github/vcomm/definitions"
)

func (c *VComm) CreateDefinitions() (*definitions.Definitions, error) {
	t := c.server.Type()
	out := &definitions.Definitions{
		Functions: make([]*definitions.Function, t.NumMethod()),
	}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fn, err := c.getMethod(m)
		if err != nil {
			return nil, err
		}
		out.Functions[i] = fn
	}

	return out, nil
}

func (c *VComm) getMethod(val reflect.Method) (*definitions.Function, error) {
	// in type (note: the first value is the struct itself)
	if val.Type.NumIn() > 2 {
		return nil, errors.New("vcomm: method must have none or 1 parameter")
	}

	var inType definitions.Type
	inName := ""
	if val.Type.NumIn() == 2 {
		var err error
		inType, err = c.getType(val.Type.In(1))
		if err != nil {
			return nil, err
		}
		inName = "value"
		// TODO: figure out a better way
	}

	// out type
	var outType definitions.Type
	if val.Type.NumOut() > 2 || val.Type.NumOut() < 1 {
		return nil, errors.New("vcomm: method must have none or 1 return value, and an error")
	}

	if val.Type.NumOut() == 1 {
		// Check if it is error
		ret := val.Type.Out(0)
		if ret.Kind() != reflect.Interface {
			return nil, errors.New("vcomm: method must return an error")
		}
		errType := reflect.TypeOf((*error)(nil)).Elem()
		if !ret.Implements(errType) {
			return nil, errors.New("vcomm: method must return an error")
		}
	}

	if val.Type.NumOut() == 2 {
		// Get first
		var err error
		outType, err = c.getType(val.Type.Out(0))
		if err != nil {
			return nil, err
		}

		// Check if second is error
		ret := val.Type.Out(1)
		if ret.Kind() != reflect.Interface {
			return nil, errors.New("vcomm: method must return an error")
		}
		errType := reflect.TypeOf((*error)(nil)).Elem()
		if !ret.Implements(errType) {
			return nil, errors.New("vcomm: method must return an error")
		}
	}

	return &definitions.Function{
		Name:       val.Name,
		ParamName:  inName,
		ParamType:  inType,
		ReturnType: outType,
	}, nil
}

func (c *VComm) getType(val reflect.Type) (definitions.Type, error) {
	if val.Kind() != reflect.Struct {
		// Basic type
		switch val.Kind() {
		case reflect.Int:
			return definitions.INT, nil

		case reflect.Float64:
			return definitions.FLOAT, nil

		case reflect.String:
			return definitions.STRING, nil

		case reflect.Slice:
			elTyp, err := c.getType(val.Elem())
			return &definitions.ArrayType{
				ElemType: elTyp,
			}, err

		case reflect.Map:
			keyTyp, err := c.getType(val.Key())
			if err != nil {
				return nil, err
			}
			valTyp, err := c.getType(val.Elem())
			return &definitions.MapType{
				KeyType:   keyTyp,
				ValueType: valTyp,
			}, err

		default:
			return nil, fmt.Errorf("vcomm: unknown type: %s", val)
		}
	}

	// Struct type
	fields := make([]definitions.StructField, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typ, err := c.getType(field.Type)
		if err != nil {
			return nil, err
		}
		fields[i] = definitions.StructField{
			Name: field.Name,
			Type: typ,
		}
	}
	return &definitions.StructType{
		Name:   val.Name(),
		Fields: fields,
	}, nil
}
