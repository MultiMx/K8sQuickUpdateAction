package kube

type Config struct {
	Backend     string `yaml:"backend"`
	Cluster     string `yaml:"cluster"`
	Container   uint   `yaml:"container"`
	BearerToken string `yaml:"token"`
}

type Request struct {
	Url   string
	Query map[string]interface{}
	Body  interface{}
}
