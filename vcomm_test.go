package vcomm

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TestServer struct {
	Messages []string
}

type TestReturnType struct {
	Val    string
	Val2   int
	Val3   float64
	Nested TestReturnTypeNested
}

type TestReturnTypeNested struct {
	Arr []int
	Map map[string]string
}

func (t *TestServer) Hi(val string) (TestReturnType, error) {
	return TestReturnType{}, nil
}

func (t *TestServer) Receive(msg string) error {
	t.Messages = append(t.Messages, msg)
	return nil
}

func TestVComm(t *testing.T) {
	serv := &TestServer{}
	comm := NewVComm(serv)
	def, err := comm.CreateDefinitions()
	if err != nil {
		t.Error(err)
		return
	}

	spew.Config.DisablePointerAddresses = true
	spew.Config.Indent = "    "
	spew.Dump(def)
}
