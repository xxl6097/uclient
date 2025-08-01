package u

type NtfyInfo struct {
	Address  string `json:"address"`
	Topic    string `json:"topic"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//	type NtfyEventData struct {
//		Id    string `json:"id,omitempty"`
//		Time  int    `json:"time,omitempty"`
//		Event string `json:"event,omitempty"`
//		Topic string `json:"topic,omitempty"`
//	}

const (
	OPEN      = "open"      // 0
	KEEPALIVE = "keepalive" // 1
	MESSAGE   = "message"   // 1
)

type NtfyEventData struct {
	Id       string `json:"id,omitempty"`
	Time     int64  `json:"time,omitempty"`
	Expires  int64  `json:"expires,omitempty"`
	Event    string `json:"event,omitempty"`
	Topic    string `json:"topic,omitempty"`
	Title    string `json:"title,omitempty"`
	Message  string `json:"message,omitempty"`
	Markdown bool   `json:"markdown,omitempty"`
}
