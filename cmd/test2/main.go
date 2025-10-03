package main

import (
	"github.com/nexgenomics/go-nexgenomics"
	"log"
)

func main() {
	log.Printf("<<<%s>>>", nexgenomics.Ping("zzz"))
}
