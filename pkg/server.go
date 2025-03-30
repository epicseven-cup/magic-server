package pkg

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"log"
	"net"
)

var MessageMaxSize = 2048

type Server struct {
	chatrooms map[uuid.UUID]*ChatRoom
}

func NewServer() *Server {
	s := &Server{
		chatrooms: make(map[uuid.UUID]*ChatRoom),
	}
	return s
}

func (s *Server) JoinChatroom(roomId uuid.UUID, clientAddress string) error {
	chatroom, ok := s.chatrooms[roomId]
	if !ok {
		return errors.New("room not found")
	}

	err := chatroom.Join(clientAddress)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) CreateRoom() (uuid.UUID, error) {
	roomId := uuid.New()
	_, ok := s.chatrooms[roomId]
	if ok {
		return [16]byte{}, errors.New("room already exists")
	}
	s.chatrooms[roomId] = NewChatRoom()
	return roomId, nil
}

func (s *Server) GetRoom(roomId uuid.UUID) (*ChatRoom, error) {
	room, ok := s.chatrooms[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (s *Server) Run() error {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		return err
	}

	log.Println("Listening on", listener.Addr())

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			log.Println("New connection from", conn.RemoteAddr())
			var data []byte
			buffer := make([]byte, 16)

			for {
				n, err := conn.Read(buffer)
				//log.Println("Read", n, err)
				data = append(data, buffer[:n]...)
				buffer = buffer[:n]
				//log.Println(string(data))

				if err != nil {
					if err.Error() == "EOF" {
						log.Println("EOF")
						break
					}
					log.Println("Error reading:", err)
					break
				}

				if len(data) > MessageMaxSize {
					log.Println("message too large")
					break
				}

			}
			var msg Message
			log.Println("reader to: ", string(data))
			log.Println("data size: ", len(data))

			var byteBuffer bytes.Buffer

			dec := gob.NewDecoder(&byteBuffer)
			err = dec.Decode(&msg)

			//err = binary.Read(reader, binary.BigEndian, &msg)
			//if err != nil {
			//	log.Println("decode: ", err)
			//	continue
			//}

			log.Println("Message Size:", msg.Size)
			log.Println("Message RoomId:", msg.RoomId)
			log.Println("Message Content:", msg.Content)

			_, err = s.GetRoom(msg.RoomId)
			if err != nil {
				log.Println(err)
			}
			//msg.Content
			//room.Broadcast()

		}
	}()
	return nil
}
