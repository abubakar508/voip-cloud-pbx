package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Websocket service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("WEBSOCKET_SERVICE_PORT")
}
