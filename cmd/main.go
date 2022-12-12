package main

import (
	_ "embed"

	"gitlab.com/germandv/sermon"
)

//go:embed services.toml
var configFileContent string

func main() {
	err := sermon.Run(configFileContent)
	if err != nil {
		panic(err)
	}
}
