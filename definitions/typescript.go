package definitions

import (
	"fmt"
	"strings"

	_ "embed"
)

//go:embed base.ts
var base string

func GenerateTypescript(definitions *Definitions) string {
	added := make(map[string]bool)
	types := &strings.Builder{}
	class := &strings.Builder{}

	for _, s := range definitions.Functions {
		addFn(types, class, added, s)
	}
	txt := strings.Replace(base, `"__types__"`, types.String(), 1)
	txt = strings.Replace(txt, `"__code__"`, strings.TrimSpace(class.String()), 1)
	return txt
}

func addFn(types *strings.Builder, class *strings.Builder, added map[string]bool, fn *Function) {
	fmt.Fprintf(class, "\tasync %s(", fn.Name)
	if fn.ParamType != nil {
		fmt.Fprintf(class, "%s: %s", fn.ParamName, addType(types, fn.ParamType, added))
	}
	class.WriteString("): Promise<")

	// Return type
	var responseTypeName string
	if fn.ReturnType != nil {
		// Create response type
		fmt.Fprintf(types, `export type %sResponse = {
	value: %s,
	error?: Error
}

`, fn.Name, addType(types, fn.ReturnType, added))
		// Return
		responseTypeName = fn.Name + "Response"
	} else {
		responseTypeName = "OptionalError"
	}
	fmt.Fprintf(class, "%s> {\n", responseTypeName)

	// Body
	txt := ", " + fn.ParamName
	if fn.ParamType == nil {
		txt = ""
	}
	fmt.Fprintf(class, "\t\tlet res = await this.createMessage(\"%s\"%s);\n", fn.Name, txt)
	fmt.Fprintf(class, "\t\treturn res as %s;\n", responseTypeName)

	// Finish
	class.WriteString("\t}\n\n")
}

func addType(types *strings.Builder, typ Type, added map[string]bool) string {
	switch typ.BasicType() {
	case INT:
		return "number"

	case FLOAT:
		return "number"

	case STRING:
		return "string"

	case ARRAY:
		return fmt.Sprintf("Array<%s>", addType(types, typ.(*ArrayType).ElemType, added))

	case MAP:
		return fmt.Sprintf("Record<%s, %s>", addType(types, typ.(*MapType).KeyType, added), addType(types, typ.(*MapType).ValueType, added))

	case STRUCT:
		break

	default:
		return "unknown"
	}

	// Struct
	s := typ.(*StructType)
	if added[s.Name] {
		return s.Name
	}
	// Type pass
	for _, f := range s.Fields {
		addType(types, f.Type, added)
	}
	// Create struct
	added[s.Name] = true
	fmt.Fprintf(types, "export type %s = {\n", s.Name)
	// Add fields
	for _, f := range s.Fields {
		fmt.Fprintf(types, "\t%s: %s;\n", f.Name, addType(types, f.Type, added))
	}
	types.WriteString("}\n\n")
	return s.Name
}
