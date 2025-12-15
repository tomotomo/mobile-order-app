package utils

import (
	"log"
	"os/exec"
	"runtime"
	"time"
)

func OpenBrowser(url string) {
	// Wait a bit for server to start
	time.Sleep(1 * time.Second)
	
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		// Unsupported platform
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
