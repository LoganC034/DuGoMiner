package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	GetPoolUrl string `json:"GetPoolURL"`
	UserName   string `json:"UserName"`
	Difficulty string `json:"Difficulty"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) GetConfig() {
	configExists := c.CheckForConfig()
	if configExists {
		file, err := os.Open("config.json")
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err = json.Unmarshal(bytes, &c); err != nil { // Parse []byte to go struct pointer
		}

	} else {
		log.Fatalln("No Config File found. 'config.json' should be located at the application root")
	}
}

// CheckForConfig TODO Check Errors
func (c *Config) CheckForConfig() bool {
	if _, err := os.Stat("config.json"); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}
}
