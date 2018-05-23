package scm

type ProjectInterface interface {
	ListUserProjects (baseUrl, token, name, namespace string) (int, error)
}


type Project struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	ID        interface{}
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}