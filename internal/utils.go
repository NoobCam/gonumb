package internal

import (
	"bufio"
	"os"
	"strings"
)

// ReadIPs reads a list of IP addresses from a given file.
func ReadIPs(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			ips = append(ips, line)
		}
	}

	return ips, scanner.Err()
}
