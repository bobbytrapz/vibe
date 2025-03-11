package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

type ClocResult struct {
	Header struct {
		ClocURL        string  `json:"cloc_url"`
		ClocVersion    string  `json:"cloc_version"`
		ElapsedSeconds float64 `json:"elapsed_seconds"`
		NFiles         int     `json:"n_files"`
		NLines         int     `json:"n_lines"`
	} `json:"header"`
	Languages map[string]ClocLanguage `json:"languages"`
	Total     ClocLanguage            `json:"SUM"`
}

type ClocLanguage struct {
	Files   int `json:"nFiles"`
	Blank   int `json:"blank"`
	Comment int `json:"comment"`
	Code    int `json:"code"`
}

func main() {
	// Define command line flags
	shouldShowRawOutput := flag.Bool("raw", false, "Show all raw command output instead of a summary")
	flag.Parse()

	// Check if repository URL is provided
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [--raw] <github-repo-url>\n", os.Args[0])
		os.Exit(1)
	}

	repoURL := flag.Arg(0)
	if !strings.HasPrefix(repoURL, "https://github.com/") {
		fmt.Println("Error: Please provide a valid GitHub repository URL")
		os.Exit(1)
	}

	// Check dependencies
	if err := checkDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create a temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "gloc-repo-*")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	fmt.Printf("Cloning repository to %s...\n", tempDir)
	if err := cloneRepo(repoURL, tempDir); err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	// Run cloc on the repository with JSON output
	fmt.Println("\nAnalyzing code...")
	output, err := runCloc(tempDir)
	if err != nil {
		fmt.Printf("Error running cloc: %v\n", err)
		os.Exit(1)
	}

	// Process and display results
	if *shouldShowRawOutput {
		fmt.Println(string(output))
	} else {
		result, err := parseClocOutput(output)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Create tabwriter
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer w.Flush()

		printResults(w, result)
	}
}

func checkDependencies() error {
	if _, err := exec.LookPath("cloc"); err != nil {
		return fmt.Errorf("cloc is not installed: %w", err)
	}
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed: %w", err)
	}
	return nil
}

func cloneRepo(repoURL, tempDir string) error {
	cloneCmd := exec.Command("git", "clone", repoURL, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	return cloneCmd.Run()
}

func runCloc(tempDir string) ([]byte, error) {
	clocCmd := exec.Command("cloc", "--json", tempDir)
	return clocCmd.CombinedOutput()
}

func parseClocOutput(output []byte) (*ClocResult, error) {
	var result map[string]any
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse cloc output: %w", err)
	}

	clocResult := &ClocResult{
		Languages: make(map[string]ClocLanguage),
	}

	// Parse header
	if header, ok := result["header"].(map[string]any); ok {
		clocResult.Header.ClocURL = header["cloc_url"].(string)
		clocResult.Header.ClocVersion = header["cloc_version"].(string)
		clocResult.Header.ElapsedSeconds = header["elapsed_seconds"].(float64)
		clocResult.Header.NFiles = int(header["n_files"].(float64))
		clocResult.Header.NLines = int(header["n_lines"].(float64))
	}

	// Parse languages and total
	for k, v := range result {
		if k == "header" {
			continue
		}

		// Convert the map to JSON bytes
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			continue
		}

		// Unmarshal into ClocLanguage
		var lang ClocLanguage
		if err := json.Unmarshal(jsonBytes, &lang); err != nil {
			continue
		}

		if k == "SUM" {
			clocResult.Total = lang
		} else {
			clocResult.Languages[k] = lang
		}
	}

	return clocResult, nil
}

func printResults(w *tabwriter.Writer, result *ClocResult) {
	// Print header
	fmt.Fprintln(w, "Language\tFiles\tCode")

	// Print language breakdown
	for lang, stats := range result.Languages {
		fmt.Fprintf(w, "%s\t%d\t%d\n", lang, stats.Files, stats.Code)
	}

	// Print separator
	fmt.Fprintln(w, "---\t---\t---")

	// Print total
	fmt.Fprintf(w, "Total\t%d\t%d\n", result.Total.Files, result.Total.Code)
}
