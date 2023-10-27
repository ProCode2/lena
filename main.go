package main

import (
	"fmt"
	"strings"
)

func parseMetaSection(metaString string) {
	// Split the meta section into lines
	lines := strings.Split(metaString, "\n")

	// Create a map to store the key-value pairs
	metaInfo := make(map[string]string)

	for _, line := range lines {
		// Split each line by ":" to separate key and value
		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Store the key-value pair in the map
			metaInfo[key] = value
		}
	}
	fmt.Println(metaInfo)
}

func main() {
	input := `title: HelloBrother whe are u doing
		description: demo
		how: you`
	parseMetaSection(input)
}
