package main

import (
	"flag"
	"io/ioutil"

	app "github.com/taufikardiyan28/chat/app"
	"github.com/taufikardiyan28/chat/helper"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config/development.yaml", "config file")
}

func ReadConfig(filname string) (*helper.Configuration, error) {
	dat, err := ioutil.ReadFile(filname)
	if err != nil {
		return nil, err
	}

	Config := &helper.Configuration{}
	err = yaml.Unmarshal(dat, Config)

	return Config, err
}

func main() {
	flag.Parse()
	cfg, err := ReadConfig(configFile)

	if err != nil {
		panic(err)
	}

	server := app.Server{Config: cfg}
	server.Start()
}
