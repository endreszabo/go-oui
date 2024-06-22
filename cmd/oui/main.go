package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dmowcomber/oui"
)

func main() {
	var macAddr string
	flag.StringVar(&macAddr, "m", "", "The Mac Address")
	flag.Parse()
	if macAddr == "" {
		flag.Usage()
		log.Fatal("must set a mac address")
	}

	o, err := oui.New()
	if err != nil {
		log.Fatalf("failed to initialize OUI data: %s", err)
	}

	org, err := o.Lookup(macAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(org)
}
