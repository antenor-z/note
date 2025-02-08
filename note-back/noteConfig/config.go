package noteConfig

import (
	"os"

	"github.com/BurntSushi/toml"
)

var config Config

func ConfigInit() {
	dat, err := os.ReadFile("config.toml")
	if err != nil {
		panic("Error while opening config file")
	}
	_, err2 := toml.Decode(string(dat), &config)
	if err2 != nil {
		panic("Error while reading toml")
	}
}

func GetDomain() string {
	return config.Domain
}

type Config struct {
	Domain string
}
