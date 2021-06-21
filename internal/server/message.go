package server

type Operation struct {
	Id        string      `json:"id"`
	Name      string      `json:"Name"`
	Arguments interface{} `json:"arguments"`
}
