package config

type Config struct {
	Nodes []Node `yaml:"nodes"`
}

type Node struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
