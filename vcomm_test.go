package vcomm

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type TestServer struct {
	Messages []string
}

type TestReturnType struct {
	Val  string
	Val2 int
	Val3 float64
}

func (t *TestServer) Hi(val string) (TestReturnType, error) {
	return TestReturnType{}, nil
}

func TestVComm(t *testing.T) {
	serv := &TestServer{}
	comm := NewVComm(serv)
	def, err := comm.CreateDefinitions()
	if err != nil {
		t.Error(err)
		return
	}

	spew.Dump(def)
}
