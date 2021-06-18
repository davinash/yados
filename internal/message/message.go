package message

type Operation struct {
	Id        string      `json:"id"`
	Type      string      `json:"type"`
	Arguments interface{} `json:"arguments"`
}
