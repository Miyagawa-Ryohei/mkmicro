package entity

import "github.com/spf13/cobra"

type WorkerCommand struct {
	Command *cobra.Command
	Handler func(Message) []byte
}

func(c *WorkerCommand) SetHandler(handler func(Message) []byte) {
	c.Handler = handler
}