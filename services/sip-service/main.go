package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("SIP service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("SIP_UDP_PORT")
}
