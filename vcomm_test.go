package vcomm

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Nv7-Github/vcomm/definitions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TestServer struct {
}

type Param struct {
	Val1 int
	Val2 []string
}

func (t *TestServer) Receive(val Param) error {
	fmt.Println(val)
	return errors.New("uh oh")
}

func TestVComm(t *testing.T) {
	serv := &TestServer{}
	comm := NewVComm(serv)
	def, err := comm.CreateDefinitions()
	if err != nil {
		t.Error(err)
		return
	}

	f, err := os.Create("test/src/generated.ts")
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

	// Dev server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			res := comm.Message(string(message))

			err = c.WriteMessage(mt, []byte(res))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
