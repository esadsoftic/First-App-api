package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Config contains database server name and database name
type Config struct {
	Server   string
	Database string
}

// Read reads and parses the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config/config.toml", &c); err != nil {
		log.Fatal(err)
	}
	//log.Println(c)
}
