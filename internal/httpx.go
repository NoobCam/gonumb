package internal

import (
	"fmt"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// RunHttpxForMasscanOutputs runs httpx for all masscan output files
func RunHttpxForMasscanOutputs(outputDir string) {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error getting current user: %v\n", err)
		return
	}
	httpxPath := filepath.Join(currentUser.HomeDir, "go/bin/httpx")

	// Debugging output to verify the httpx path
	fmt.Printf("Using httpx tool at: %s\n", httpxPath)

	// Find all masscan output files
	masscanFiles, err := filepath.Glob(filepath.Join(outputDir, "masscan_*.txt"))
	if err != nil {
		fmt.Printf("Error finding masscan output files: %v\n", err)
		return
	}

	for _, masscanFile := range masscanFiles {
		httpxFilePath := filepath.Join(outputDir, fmt.Sprintf("httpx_%s.txt", strings.ReplaceAll(filepath.Base(masscanFile), ".txt", "")))

		// Prepare command to run httpx with the necessary filters
		cmd := exec.Command("sh", "-c", fmt.Sprintf("grep tcp %s | awk '{print $4, \":\", $3}' | tr -d ' ' | %s -title -sc -cl -tech-detect -no-color > %s", masscanFile, httpxPath, httpxFilePath))

		// Run the command and capture output
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running httpx scan for %s: %v\n", masscanFile, err)
			return
		}

		fmt.Printf("httpx completed for %s. Results saved to %s\n", masscanFile, httpxFilePath)
	}
}
