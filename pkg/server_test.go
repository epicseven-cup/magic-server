package pkg

import (
	"bytes"
	"encoding/gob"
	"github.com/google/uuid"
	"log"
	"net"
	"testing"
	"time"
)

func init() {
	s := NewServer()

	log.Println("Running server")
	go func() {
		err := s.Run()
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)

}

func TestNewServer(t *testing.T) {

	conn, err := net.Dial("tcp", ":8000")
	//deadline := 40 * time.Second
	//err = conn.SetDeadline(time.Now().Add(deadline))
	if err != nil {
		log.Fatal("Error setting write deadline:", err)
		return
	}
	if err != nil {
		t.Error("Error connecting to server:", err)
		return
	}

	t.Log("Connected to", conn.RemoteAddr())
	content := [2048]byte{}
	copy(content[:2048], "Hello World")
	id := uuid.New()

	msg := Message{
		Size:    1,
		RoomId:  id,
		Content: content,
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
	log.Println("Written", len(data.Bytes()), err)

	log.Println("Write to", data.String())
	if err != nil {

		t.Error("Error writing to server:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			t.Error("Error closing connection:", err)
		}
	}(conn)

	t.Log("Message Size:", msg.Size)
	t.Log("Message RoomId:", msg.RoomId)
	t.Log("Message Content:", msg.Content)

}
