package internal

import (
	"encoding/binary"
	"github.com/epicseven-cup/magic-server/pkg"
	"github.com/google/uuid"
	"log"
	"net"
)

func Encode(conn net.Conn) (*pkg.Message, error) {

	size := make([]byte, 4)
	roomId := make([]byte, 16)

	n, err := conn.Read(size)
	log.Println("Read: ", n)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("size", size)

	msgType := binary.LittleEndian.Uint16(size)

	n, err = conn.Read(roomId)
	log.Println("Read: ", n)
	log.Println("uuid", roomId)

	content := make([]byte, 2048)
	n, err = conn.Read(content)
	log.Println("Read: ", n)

	msg := pkg.Message{
		Type:    msgType,
		RoomId:  uuid.UUID(roomId),
		Content: [2048]byte(content),
	}

	return &msg, nil
}
