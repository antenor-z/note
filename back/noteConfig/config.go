package noteConfig

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Auth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Totp     string `toml:"totp"`
}

type Server struct {
	Domain    string `toml:"domain"`
	DebugMode bool   `toml:"debugmode"`
}
type Config struct {
	Server Server
	Auth   Auth
}

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
	return config.Server.Domain
}

func IsDebug() bool {
	return config.Server.DebugMode
}

func GetAuthSecret() Auth {
	return config.Auth
}

func GetVersion() string {
	dat, err := os.ReadFile("VERSION")
	if err != nil {
		dat, err := os.ReadFile("../VERSION")
		if err != nil {
			panic("Error while reading version file")
		}
		return string(dat)
	}
	return string(dat)
}
