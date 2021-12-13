package vcomm

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type sendDataVal struct {
	value string
}

func (s *sendDataVal) UnmarshalJSON(data []byte) error {
	s.value = string(data)
	return nil
}

type sendData struct {
	Fn    string
	Value sendDataVal
}

func (v *VComm) AcceptMessage(message string) string {
	var dat sendData
	err := json.Unmarshal([]byte(message), &dat)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}

	// Get method
	methType, exists := v.server.Type().MethodByName(dat.Fn)
	if !exists {
		return fmt.Sprintf(`{"error": "vcomm: method %s not found"}`, dat.Fn)
	}
	method := v.server.MethodByName(dat.Fn)
	inps := make([]reflect.Value, 0)
	if methType.Type.NumIn() == 2 {
		// Get input
		v := reflect.Zero(methType.Type.In(1)).Interface()
		err = json.Unmarshal([]byte(dat.Value.value), &v)
		if err != nil {
			return fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
		inps = append(inps, reflect.ValueOf(v))
	}
	res := method.Call(inps)

	// Create return
	out := make(map[string]interface{})

	// Get error
	errV := res[len(res)-1].Interface()
	if errV != nil {
		out["error"] = errV.(error).Error()
	}
	if methType.Type.NumOut() == 2 {
		val := res[0].Interface()
		out["value"] = val
	}

	ret, err := json.Marshal(out)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(ret)
}
