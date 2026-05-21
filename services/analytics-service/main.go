package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Analytics service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("ANALYTICS_SERVICE_PORT")
}
