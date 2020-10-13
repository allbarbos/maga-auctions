package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// EnvVars is the configuration of environment
var EnvVars Config

func init() {
	readFile(&EnvVars)
	readEnv()
	fmt.Printf("%+v", EnvVars)
}

// Config contains the mapping of environment variables
type Config struct {
	API struct {
		Env  string `yaml:"env", envconfig:"API_ENV"`
		Port string `yaml:"port", envconfig:"API_PORT"`
	} `yaml:"api"`

	Legacy struct {
		URI string `yaml:"uri", envconfig:"LEGACY_URI"`
	} `yaml:"legacy"`
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

type API struct {
	Env  string
	Port string
}

func readEnv() {
	err := envconfig.Process("", &EnvVars)
	if err != nil {
		processError(err)
	}
}
