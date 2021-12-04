package main

import "mkmicro/adapter/cmd"

func main () {
	Initialize()
	if err := cmd.RootCmd.Execute(); err !=nil {
		panic(err)
	}
}