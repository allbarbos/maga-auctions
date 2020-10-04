package env

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Vars is the configuration of environment
var Vars Config

func init() {
	readFile(&Vars)
	readEnv(&Vars)
	fmt.Printf("%+v", Vars)
}

// Config contains the mapping of environment variables
type Config struct {
	API struct {
		Env  string `yaml:"env", envconfig:"API_ENV"`
		Port string `yaml:"port", envconfig:"API_PORT"`
	} `yaml:"api"`

	APILegacy struct {
		URI string `yaml:"uri", envconfig:"API_LEGACY"`
	} `yaml:"api-legacy"`
}

func processError(err error) {
	log.Print(err)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}

	err = f.Close()

	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
