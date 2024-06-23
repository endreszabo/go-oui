package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dmowcomber/oui"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("must set a mac address")
	}
	macAddr := args[0]

	o, err := oui.New()
	if err != nil {
		log.Fatalf("failed to initialize OUI data: %s", err)
	}

	org, err := o.Lookup(macAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(org.Name)
}
