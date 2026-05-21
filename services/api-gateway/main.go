package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("API Gateway starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("API_GATEWAY_PORT")
}
