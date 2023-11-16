package server

import (
	"github.com/go-micro/plugins/v4/config/encoder/yaml"
	"os"
)

type Server struct {
	Name string `yaml:"Name"`
	Port string `yaml:"Port"`
}

var (
	Config = &Server{}
)

func (s Server) GetFilePath() string {
	return "./env/config.yaml"
}

func (s Server) Init() error {
	body, err := s.read(s.GetFilePath())
	if err != nil {
		return err
	}
	return yaml.NewEncoder().Decode(body, Config)
}

func (s Server) Handle() error {
	return nil
}

func (s Server) read(name string) ([]byte, error) {
	return os.ReadFile(name)
}
