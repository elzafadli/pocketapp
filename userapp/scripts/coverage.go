package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run coverage.go <min_coverage>")
		os.Exit(1)
	}

	minCoverage, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid minimum coverage value: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "tool", "cover", "-func=cp.out")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running go tool cover: %v\n%s\n", err, string(output))
		os.Exit(1)
	}

	fmt.Println(string(output))
	fmt.Println("-----------------------")

	re := regexp.MustCompile(`total:\s+\(statements\)\s+([0-9.]+)%`)
	match := re.FindStringSubmatch(string(output))
	if len(match) < 2 {
		fmt.Println("Could not find total coverage in output")
		os.Exit(1)
	}

	coverage, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		fmt.Printf("Error parsing coverage value: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Coverage: %.1f%%\n", coverage)
	if coverage < minCoverage {
		fmt.Printf("Coverage %.1f%% is below minimum required %.1f%%\n", coverage, minCoverage)
		os.Exit(1)
	}
}
