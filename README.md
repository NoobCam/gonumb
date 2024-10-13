
# Gonumb

Gonumb is a network reconnaissance tool designed to facilitate the discovery of IP addresses and associated information through ASN and IP block lookups, as well as port scanning and HTTP probing using tools like Nmap, Masscan, and httpx.

## Features

- **ASN Lookup**: Read a list of ASNs from a file and retrieve associated IP blocks using Nmap.
- **IP Block Lookup**: Read a list of CIDR IP blocks from a file and perform reverse DNS lookups.
- **Port Scanning**: Utilize Masscan to scan for open ports on the discovered IP addresses or blocks.
- **HTTP Probing**: Use httpx to gather additional information about the services running on the discovered IPs.

## Prerequisites

Before using Gonumb, ensure that you have the following installed on your system:

- **Go**: You need to have Go installed. You can download it from [golang.org](https://golang.org/dl/).
- **Nmap**: Install Nmap for network discovery and security auditing. Instructions can be found at [nmap.org](https://nmap.org/download.html).
- **Masscan**: Install Masscan for high-speed port scanning. Instructions can be found at [masscan.org](https://github.com/robertdavidgraham/masscan).
- **httpx**: Install httpx for probing URLs to discover alive hosts. Instructions can be found at [projectdiscovery.io](https://github.com/projectdiscovery/httpx).

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/<your_username>/gonumb.git
   cd gonumb
   ```

2. Install the required Go modules:

   ```bash
   go mod tidy
   ```

3. Ensure that the `httpx` binary is in your `$HOME/go/bin/` directory. You can install it using:

   ```bash
   go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
   ```

## Usage

Gonumb accepts several command-line flags to operate. Below are the flags and examples:

### Flags

- `--asn <asn_file>`: Path to a file containing a list of ASNs (one ASN per line).
- `--ipblock <ip_block_file>`: Path to a file containing CIDR-formatted IP blocks (one per line).
- `--outdir <output_directory>`: Directory where output files will be saved.

### Examples

1. **Using ASN File**

   To use an ASN file to discover associated IP blocks:

   ```bash
   go run cmd/main.go --asn asn_list.txt --outdir output
   ```

2. **Using IP Block File**

   To use an IP block file for reverse DNS lookups and mass scanning:

   ```bash
   go run cmd/main.go --ipblock cidr_list.txt --outdir output
   ```

## Output

The output will be saved in the specified directory. The following files will be generated:

- `ASN_<asn>_ipblocks.txt`: Contains the IP blocks discovered for each ASN.
- `<ip_block>_hostnames.txt`: Contains the hostnames discovered for each IP block.
- `masscan_<ip_block>.txt`: Contains the results of the masscan port scan.
- `httpx_<masscan_file>.txt`: Contains the results of the httpx probing based on the masscan output.

## Notes

- **Running with Root Privileges**: Masscan requires root privileges to run. You may need to execute your command with `sudo`.
- **Input Files**: Ensure that your input files (`asn_list.txt`, `cidr_list.txt`, etc.) are formatted correctly with one entry per line.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

If you would like to contribute to this project, please open an issue or submit a pull request.

## Contact

For any inquiries or feedback, please reach out to <your_email@example.com>.
