package main

import (
	"fmt"
	"log"
	"os"

	"github.com/endreszabo/go-oui"
)

func main() {
	args := os.Args[2:]
	if len(args) < 2 {
		log.Fatal("must set a mac address")
	}
	macAddr := args[0]

	o, err := oui.New(os.Args[1], true)
	if err != nil {
		log.Fatalf("failed to initialize OUI data: %s", err)
	}

	org, err, _ := o.Lookup(macAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(org.Name)
}
