package entity

type Selector struct {
	Id string `json:"id"`
	Nodes []Node `json:"nodes"`
}

type Node struct {
	Profile string `json:"profile"`
	Project string `json:"project"`
	NodeSelector string `json:"node_selector"`
}


