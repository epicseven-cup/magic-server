package pkg

import "github.com/google/uuid"

type Message struct {
	Type    uint16     // 4
	RoomId  uuid.UUID  // 16
	Content [2048]byte // 2048
	// totoal -> 4 + 16 + 2048 = 2068
}
