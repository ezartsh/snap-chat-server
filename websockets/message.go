package websockets

const (
	MT_PrivateChat = "mt_private_chat" // message type for act send rivate chat message
	MT_GroupChat   = "mt_group_chat"   // message type for act send broadcast chat message
)

type Message struct {
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
}
