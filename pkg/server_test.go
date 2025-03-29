package pkg

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"log"
	"net"
	"testing"
)

func init() {
	s := NewServer()

	log.Println("Running server")
	err := s.Run()

	if err != nil {
		log.Fatal("Error starting server:", err)
		return
	}

}

func TestNewServer(t *testing.T) {

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		t.Error("Error connecting to server:", err)
		return
	}

	t.Log("Connected to", conn.RemoteAddr())
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.Error("Error closing connection:", err)
		}
	}(conn)

	msg := Message{
		Size:    1,
		RoomId:  uuid.New().String(),
		Content: "Hello World!",
	}
	data := new(bytes.Buffer)

	enc := gob.NewEncoder(data)
	err = enc.Encode(msg)
	if err != nil {
		t.Error(err)
		return
	}

	err = enc.Encode(msg)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = conn.Write(data.Bytes())
	if err != nil {
		t.Error(err)
		return
	}

}
