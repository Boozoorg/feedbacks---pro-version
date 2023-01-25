package models

const (
	Text = 1
	File = 2
)

type Message struct {
	Sender      string `json:"sender"`
	Message     string `json:"message"`
	MesaageType int64  `json:"message_type"`
	Time        string `json:"time"`
}

type Receiver struct {
	Name        string
	Room_id     int64
	Message     string
	MessageType int64
}

type Rooms struct {
	Max     int64   `json:"max"`
	RoomsConf []string `json:"rooms"`
}

type MessageToBot struct {
	ProductID int64
	RoomID    int64
	Fio       string
	Message   string
	Time      string
}
