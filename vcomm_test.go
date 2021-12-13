package vcomm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Nv7-Github/vcomm/definitions"
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
	fmt.Println(t.Messages)
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

	f, err := os.Create("test.ts")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(definitions.GenerateTypescript(def))
	if err != nil {
		t.Error(err)
		return
	}

	// Test receive
	res := comm.AcceptMessage(`{"Fn": "Hi", "Value": "Hello"}`)
	fmt.Println(res)
}
