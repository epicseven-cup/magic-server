package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var MessageMaxSize, _ = strconv.Atoi(os.Getenv("MESSAGE_MAX_SIZE"))

type Server struct {
	chatrooms map[string]*ChatRoom
}

func NewServer() *Server {
	s := &Server{
		chatrooms: make(map[string]*ChatRoom),
	}
	return s
}

func (s *Server) JoinChatroom(roomId string, clientAddress string) error {
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

func (s *Server) CreateRoom() (string, error) {
	roomId := uuid.New().String()
	_, ok := s.chatrooms[roomId]
	if ok {
		return "", errors.New("room already exists")
	}
	s.chatrooms[roomId] = NewChatRoom()
	return roomId, nil
}

func (s *Server) GetRoom(roomId string) (*ChatRoom, error) {
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
			}
			log.Println("New connection from", conn.RemoteAddr())
			buffer := make([]byte, 1024)
			data := make([]byte, MessageMaxSize)

			log.Println("buffer size: ", len(buffer))

			for {
				n, err := conn.Read(buffer)
				log.Println("amount read: ", n)

				data = append(data, buffer[:n]...)
				buffer = buffer[n:]
				log.Println(string(data))

				if err != nil && err != io.EOF {
					log.Println(err)
					break
				}
				if err == io.EOF {
					log.Println("EOF")
					break
				}

				if len(data) > MessageMaxSize {
					log.Println("message too large")
					break
				}

			}

			msg := &Message{}

			err = binary.Read(bytes.NewBuffer(data), binary.LittleEndian, msg)
			if err != nil {
				log.Println(err)
			}

			log.Println("Message received: ", msg)

			room, err := s.GetRoom(msg.RoomId)
			if err != nil {
				log.Println(err)
			}
			room.Broadcast(msg.Content)

		}
	}()
	return nil
}
