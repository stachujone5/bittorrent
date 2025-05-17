package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: bittorrent <torrent-file> <output-directory>")
		os.Exit(1)
	}

	torrentFilePath := os.Args[1]

	content, err := os.ReadFile(torrentFilePath)
	if err != nil {
		fmt.Println("Error reading torrent file:", err)
		os.Exit(1)
	}

	outputDir := os.Args[2]

	_, err = os.Stat(outputDir)
	if err != nil {
		fmt.Println("Error checking output directory:", err)
		os.Exit(1)
	}

	fmt.Println(string(content))

	fmt.Println(
		"BitTorrent client initialized. Download would start here in a complete implementation.",
	)
}
