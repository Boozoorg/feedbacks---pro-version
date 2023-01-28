package models

type Message struct {
	ID          int64  `json:"id"`
	Sender      string `json:"sender"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	Time        string `json:"time"`
}

type ReceiveMessage struct {
	ID          int64
	Role        int64
	Sender      string
	Room        string `json:"room"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type User struct {
	ID             string
	Room           string
	Name           string
	Role           int64
	ImportantLevel int64
}

type Rooms struct {
	Max       int64       `json:"max"`
	RoomsConf []RoomsConf `json:"room_conf"`
}

type RoomsConf struct {
	RoomID         string `json:"room_id"`
	Name           string `json:"name"`
	Time           string `json:"time"`
	ImportantLevel int64  `json:"important_level"`
	Counter        int    `json:"counter"`
}

type MessageToBot struct {
	ProductID int64
	RoomID    int64
	Fio       string
	Message   string
	Time      string
}
