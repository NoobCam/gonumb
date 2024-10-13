package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// Scan function to perform masscan on each IP block
func Scan(ipBlocks []string, outputDir string) {
	// Prompt for sudo password
	fmt.Print("Enter your password to run masscan with sudo: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	fmt.Println() // Move to the next line after password input

	for _, ipBlock := range ipBlocks {
		masscanIPBlock(ipBlock, outputDir, string(password)) // Pass password to masscan
	}
}

func masscanIPBlock(ipBlock string, outputDir string, password string) {
	masscanFilePath := filepath.Join(outputDir, fmt.Sprintf("masscan_%s.txt", strings.ReplaceAll(ipBlock, "/", "_")))

	cmd := exec.Command("sudo", "masscan", "--open-only", ipBlock, "-p1-65535,U:1-65535",
		"--rate=10000", "--http-user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0", "-oL", masscanFilePath)

	// Create a pipe to capture stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe for Masscan: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting Masscan for IP block %s: %v\n", ipBlock)
		return
	}

	// Read the output line by line
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				fmt.Println(line) // Print the line to show progress
			}
		}
	}()

	// Wait for the command to finish and capture any errors
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error running Masscan for IP block %s: %v\n", ipBlock)
		return
	}

	// Notify the user of completion
	fmt.Printf("Masscan completed for IP block %s. Results saved to %s\n", ipBlock, masscanFilePath)
}
