package pkg

import "testing"

func TestNewChatRoom(t *testing.T) {
	c := NewChatRoom()
	if c.GetClients() == nil {
		t.Error("NewChatRoom returned wrong client list")
	}
}

func TestChatRoom_Broadcast(t *testing.T) {
	c := NewChatRoom()
	c.Broadcast("hello world")
}

func TestChatRoom_GetClients(t *testing.T) {
	c := NewChatRoom()
	clients := c.GetClients()
	if clients == nil {
		t.Error("GetClients returned nil")
	}
}

func TestChatRoom_RemoveClient(t *testing.T) {
	c := NewChatRoom()
	err := c.Join("Client1")
	if err != nil {
		t.Error("ChatRoom failed to join client")
		return
	}
	clients := c.GetClients()

	if v, ok := clients["Client1"]; !ok || v != true {
		t.Error("ChatRoom failed to add client")
	}

	err = c.RemoveClient("Client1")
	if err != nil {
		t.Error("ChatRoom failed to remove client")
		return
	}
	clients = c.GetClients()
	if clients["Client1"] == true {
		t.Error("ChatRoom failed to remove client")
	}

}
