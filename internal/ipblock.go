package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ReadIPBlocks reads IP blocks from a file.
func ReadIPBlocks(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ipBlocks []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			ipBlocks = append(ipBlocks, line)
		}
	}

	return ipBlocks, scanner.Err()
}

// ReverseLookup performs a reverse DNS lookup for multiple IP blocks.
func ReverseLookup(ipBlocks []string) []string {
	var hostnames []string

	for _, ipBlock := range ipBlocks {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -s 'https://rapiddns.io/sameip/%s?full=1#result' | grep 'target=\"' -B1 | grep -E -v '(--|) ' | rev | cut -c 6- | rev | cut -c 5- | sort -u", ipBlock))

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error running reverse lookup for IP block %s: %v\n", ipBlock)
			continue
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		hostnames = append(hostnames, lines...)
	}

	return hostnames
}

// ReverseLookupSingle performs a reverse DNS lookup for a single IP block.
func ReverseLookupSingle(ipBlock string) []string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -s 'https://rapiddns.io/sameip/%s?full=1#result' | grep 'target=\"' -B1 | grep -E -v '(--|) ' | rev | cut -c 6- | rev | cut -c 5- | sort -u", ipBlock))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running reverse lookup for IP block %s: %v\n", ipBlock)
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines
}

// SaveHostnames saves hostnames to a specified file.
func SaveHostnames(outputDir string, identifier string, hostnames []string) {
	// Replace the `/` with `_` in CIDR for the filename
	cleanIdentifier := strings.ReplaceAll(identifier, "/", "_")
	filename := filepath.Join(outputDir, fmt.Sprintf("%s_hostnames.txt", cleanIdentifier))

	// Create the directory for saving hostnames if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory for hostnames: %v\n", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error saving hostnames for %s: %v\n", identifier, err)
		return
	}
	defer file.Close()

	for _, hostname := range hostnames {
		file.WriteString(hostname + "\n")
	}
}
