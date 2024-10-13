package main

import (
	"flag"
	"fmt"
	"gonumb/internal" // Importing the internal package
	"os"
)

func main() {
	// Print usage information if no flags are provided
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Define flags for ASN file, IP block file, and output directory
	asnFile := flag.String("asn", "", "Path to the input ASN file")
	ipBlockFile := flag.String("ipblock", "", "Path to the input IP block file")
	outputDir := flag.String("outdir", "", "Directory to save output files")
	flag.Parse()

	// Validate flags
	if *asnFile == "" && *ipBlockFile == "" {
		fmt.Println("Usage: gonumb --asn <asn_file> --ipblock <ip_block_file> --outdir <output_directory>")
		return
	}
	if *outputDir == "" {
		fmt.Println("Output directory is required.")
		return
	}

	// Ensure the output directory exists
	if err := os.MkdirAll(*outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	if *asnFile != "" {
		// Read ASNs from the file
		asns, err := internal.ReadASNs(*asnFile)
		if err != nil {
			fmt.Printf("Error reading ASNs: %v\n", err)
			return
		}

		// Process each ASN
		for _, asn := range asns {
			ipBlocks := internal.RunNmap(asn)
			internal.SaveIPBlocks(*outputDir, asn, ipBlocks)

			hostnames := internal.ReverseLookup(ipBlocks)
			for _, ipBlock := range ipBlocks {
				internal.SaveHostnames(*outputDir, ipBlock, hostnames)
			}
		}
	}

	if *ipBlockFile != "" {
		// Read IP blocks from the file
		ipBlocks, err := internal.ReadIPBlocks(*ipBlockFile)
		if err != nil {
			fmt.Printf("Error reading IP blocks: %v\n", err)
			return
		}

		// Perform reverse DNS lookup for IP blocks
		for _, ipBlock := range ipBlocks {
			hostnames := internal.ReverseLookupSingle(ipBlock)
			internal.SaveHostnames(*outputDir, ipBlock, hostnames)
		}

		// Scan open ports with Masscan
		internal.Scan(ipBlocks, *outputDir)
	}

	// Run httpx for all masscan output files
	internal.RunHttpxForMasscanOutputs(*outputDir)

	fmt.Println("Processing complete.")
}

// printUsage prints the usage information for the program.
func printUsage() {
	fmt.Println("gonumb: A tool for network reconnaissance. This program takes an ASN or IPblock for input, then gets DNS names from rapiddns, does port enumeration with masscan and runs httpx against the open ports")
	fmt.Println("\nUsage:")
	fmt.Println("  gonumb --asn <asn_file> --outdir <output_directory>")
	fmt.Println("  gonumb --ipblock <ip_block_file> --outdir <output_directory>")
	fmt.Println("\nFlags:")
	fmt.Println("  --asn      Path to the input ASN file")
	fmt.Println("  --ipblock  Path to the input IP block file (CIDR format)")
	fmt.Println("  --outdir   Directory to save output files")
	fmt.Println("\nExamples:")
	fmt.Println("  gonumb --asn asn_list.txt --outdir output")
	fmt.Println("  gonumb --ipblock cidr_list.txt --outdir output")
}
