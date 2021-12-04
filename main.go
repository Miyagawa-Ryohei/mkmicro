package main

import "github.com/Miyagawa-Ryohei/mkmicro/adapter/cmd"

func main () {
	Initialize()
	if err := cmd.RootCmd.Execute(); err !=nil {
		panic(err)
	}
}