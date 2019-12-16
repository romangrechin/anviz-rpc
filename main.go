package main

import (
	"encoding/json"
	"fmt"
	"github.com/romangrechin/anviz-rpc/api"
	"io/ioutil"
	"log"
	"os"
)

var cfg *config

type config struct {
	Host  string `json:"host"`
	Token string `json:"api-key"`
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-c" {
		fmt.Println("Usage: anviz-rpc -c [path to config.json]")
		os.Exit(1)
	}

	err := parseConfig(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	if cfg != nil {
		fmt.Println("Server started on: ", cfg.Host)
		api.RunServer(cfg.Host, cfg.Token)
		os.Exit(0)
	}
	os.Exit(1)
}

func parseConfig(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	cfg = &config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
