package wechatwork

type MessageType string

const (
	MessageTypeText MessageType = "text"
)

const allUser = "@all"

type Message struct {
	ToUser  string      `json:"touser"`
	AgentId int         `json:"agentid"`
	Type    MessageType `json:"msgtype"`
	Text    *text       `json:"text,omitempty"`
}

type text struct {
	Content string `json:"content"`
}

func NewTextMessage(content string) *Message {
	return &Message{
		ToUser: allUser,
		Type:   MessageTypeText,
		Text:   &text{Content: content},
	}
}

func (msg *Message) Content() string {
	switch msg.Type {
	case MessageTypeText:
		return msg.Text.Content
	default:
		return "nothing"
	}
}
