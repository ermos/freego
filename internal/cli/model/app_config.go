package model

type AppConfig struct {
	Link    string            `yaml:"-"`
	Domains map[string]Domain `yaml:"domains"`
}

type Domain struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
