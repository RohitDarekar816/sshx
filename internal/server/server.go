package server

type Server struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}
