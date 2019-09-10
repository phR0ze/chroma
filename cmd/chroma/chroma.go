package main

import (
	"fmt"

	"github.com/phR0ze/chroma/internal/chroma"
)

func main() {
	c := chroma.New()
	if err := c.Execute(); err != nil {
		logFatal(c, err)
	}
}

func logFatal(c *chroma.Chroma, err error) {
	if argsErr, ok := err.(chroma.ArgsError); ok {
		argsErr.Command.Help()
		fmt.Println()
		c.LogError("Argument failure: the '%s' command's arguments were not satisfied", argsErr.Command.Name())
	}
	c.LogFatal(err)
}
