package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Recording service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("RECORDING_SERVICE_PORT")
}
