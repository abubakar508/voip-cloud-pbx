package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("AI service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("AI_SERVICE_PORT")
}
