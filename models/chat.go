package models

type User struct {
	ID      int64
	Room_id int64
	Name    string
	Role    int64
}

type SendMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Message struct {
	ID      int64  `json:"id"`
	Role    int64  `json:"role"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type ReceiveMessage struct {
	User        User
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type Rooms struct {
	Max        int64       `json:"max"`
	Rooms_conf []RoomsConf `json:"room_conf"`
}

type RoomsConf struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Time           string `json:"time"`
	ImportantLevel int64  `json:"important_level"`
}

type MessageToBot struct {
	ProductID int64
	RoomID    int64
	Fio       string
	Message   string
	Time      string
}
