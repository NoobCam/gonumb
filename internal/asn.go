package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ReadASNs reads ASNs from a file.
func ReadASNs(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var asns []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			var asn int
			if _, err := fmt.Sscanf(line, "%d", &asn); err == nil {
				asns = append(asns, asn)
			}
		}
	}

	return asns, scanner.Err()
}

// RunNmap runs the Nmap command for a given ASN.
func RunNmap(asn int) []string {
	cmd := exec.Command("nmap", "--script", "targets-asn", "--script-args", fmt.Sprintf("targets-asn.asn=%d", asn))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running Nmap for ASN %d: %v\n", asn)
		return nil
	}
	return extractIPBlocks(string(output))
}

// SaveIPBlocks saves the IP blocks to a file.
func SaveIPBlocks(outputDir string, asn int, ipBlocks []string) {
	ipBlockFilePath := filepath.Join(outputDir, fmt.Sprintf("ASN_%d_ipblocks.txt", asn))
	file, err := os.Create(ipBlockFilePath)
	if err != nil {
		fmt.Printf("Error saving IP blocks for ASN %d: %v\n", asn)
		return
	}
	defer file.Close()

	for _, block := range ipBlocks {
		file.WriteString(block + "\n")
	}
}

// extractIPBlocks parses the output of Nmap to find IP blocks.
func extractIPBlocks(output string) []string {
	var ipBlocks []string
	re := regexp.MustCompile(`\b\d{1,3}(\.\d{1,3}){3}/\d{1,2}\b`)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindString(line)
		if match != "" {
			ipBlocks = append(ipBlocks, match)
		}
	}
	return ipBlocks
}
