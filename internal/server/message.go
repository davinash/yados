package server

type Request struct {
	Id        string      `json:"id"`
	Name      string      `json:"Name"`
	Arguments interface{} `json:"arguments"`
}

type Response struct {
	Id       string
	Response interface{}
	err      error
}
