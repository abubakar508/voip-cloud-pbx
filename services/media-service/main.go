package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Media service starting (Phase 1.2 minimal entrypoint)")
	_ = os.Getenv("MEDIA_RTP_BASE_PORT")
}
