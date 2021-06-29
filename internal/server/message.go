package server

type OperationId int

const (
	AddNewMember OperationId = iota
	PutObject
	GetObject
	DeleteObject
	CreateStoreInCluster
	DeleteStoreFromCluster
)

type Request struct {
	Id        OperationId `json:"id"`
	Arguments interface{} `json:"arguments"`
}

type Response struct {
	Id   string      `json:"id"`
	Resp interface{} `json:"response"`
}

type JoinMember struct {
	Port        int    `json:"Port"`
	Address     string `json:"Address"`
	ClusterName string `json:"ClusterName"`
	Name        string `json:"Name"`
}

type JoinMemberResp struct {
	Port        int    `json:"Port"`
	Address     string `json:"Address"`
	ClusterName string `json:"ClusterName"`
	Name        string `json:"Name"`
}
