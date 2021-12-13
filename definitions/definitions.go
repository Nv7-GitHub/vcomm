package definitions

import (
	"fmt"
	"strings"
)

type Type interface {
	BasicType() BasicType
	Equal(Type) bool
	String() string
}

type BasicType int

const (
	INT BasicType = iota
	FLOAT
	STRING
	BOOL
	ARRAY
	MAP
	STRUCT
)

var typeNames = map[BasicType]string{
	INT:    "int",
	FLOAT:  "float",
	STRING: "string",
	BOOL:   "bool",
}

func (b BasicType) BasicType() BasicType {
	return b
}

func (b BasicType) String() string {
	return typeNames[b]
}

func (b BasicType) Equal(t Type) bool {
	return b == t.BasicType()
}

type StructField struct {
	Name string
	Type Type
}

type ArrayType struct {
	ElemType Type
}

func (a *ArrayType) BasicType() BasicType {
	return ARRAY
}

func (a *ArrayType) Equal(b Type) bool {
	if b.BasicType() != ARRAY {
		return false
	}

	if !a.ElemType.Equal(b.(*ArrayType).ElemType) {
		return false
	}

	return true
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[]%s", a.ElemType.String())
}

type MapType struct {
	KeyType   Type
	ValueType Type
}

func (m *MapType) BasicType() BasicType {
	return MAP
}

func (m *MapType) Equal(b Type) bool {
	if b.BasicType() != MAP {
		return false
	}

	if !m.KeyType.Equal(b.(*MapType).KeyType) {
		return false
	}

	if !m.ValueType.Equal(b.(*MapType).ValueType) {
		return false
	}

	return true
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.KeyType.String(), m.ValueType.String())
}

type StructType struct {
	Name   string
	Fields []StructField
}

func (s *StructType) BasicType() BasicType {
	return STRUCT
}

func (s *StructType) Equal(t Type) bool {
	if t.BasicType() != STRUCT {
		return false
	}
	st := t.(*StructType)
	if len(s.Fields) != len(st.Fields) {
		return false
	}
	for i, f := range s.Fields {
		if !f.Type.Equal(st.Fields[i].Type) {
			return false
		}
		if f.Name != st.Fields[i].Name {
			return false
		}
	}
	return true
}

func (s *StructType) String() string {
	out := &strings.Builder{}
	out.WriteString("struct {")
	for i, f := range s.Fields {
		out.WriteString(f.Name)
		out.WriteString(": ")
		out.WriteString(f.Type.String())
		if i != len(s.Fields)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("}")
	return out.String()
}

type FunctionParameter struct {
	Name string
	Type Type
}

type Function struct {
	Name       string
	ParamType  Type // may be nil, indicating no parameters
	ReturnType Type // may be nil, indicating no return
	// Functions also have error returns
}

type Definitions struct {
	Functions []*Function
}
