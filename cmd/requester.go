package main

import (
	"github.com/armadanet/requester"
	"log"
)

func main() {
	err := requester.Run("spinner:5912")
	if err != nil {log.Fatalln(err)}
}