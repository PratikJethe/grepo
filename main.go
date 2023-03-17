package main

import (
	"github.com/pratikjethe/grepo/cmd"
	"github.com/pratikjethe/grepo/grepo"
)

func main() {
	parsedInput := cmd.GetParsedInput()
	grepo.GrepSearch(parsedInput)

}
