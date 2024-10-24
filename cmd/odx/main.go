package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	core "github.com/3zcurdia/odx/core"
	parsers "github.com/3zcurdia/odx/parsers"

	_ "github.com/mattn/go-sqlite3"
)

const VERSION = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "convert":
		handleConvert(os.Args[2:])
	case "version":
		fmt.Printf("odx version %s\n", VERSION)
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func handleConvert(args []string) {
	// Setup flags for convert command
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	outputFlag := convertCmd.String("o", "", "Output file path (optional)")

	err := convertCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	// Need at least one input file
	if convertCmd.NArg() < 1 {
		fmt.Println("Error: No input file specified")
		fmt.Println("Usage: odx convert [options] <input-file>")
		os.Exit(1)
	}

	inputFile := convertCmd.Arg(0)
	outputFile := *outputFlag

	// If no output file specified, use input file name with .odx extension
	if outputFile == "" {
		outputFile = replaceExtension(inputFile, ".odx")
	}

	// Perform conversion
	if err := convertFile(inputFile, outputFile); err != nil {
		fmt.Printf("Error converting file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func convertFile(input, output string) error {
	if _, err := os.Stat(input); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", input)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		return fmt.Errorf("output file already exists: %s", output)
	}

	db, err := core.Init(output)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	defer db.Close()

	mesh, err := loadFile(input)
	if err != nil {
		return fmt.Errorf("error loading file: %v", err)
	}

	if err := core.Insert(db, mesh); err != nil {
		return fmt.Errorf("error inserting data: %v", err)
	}

	return nil
}

func loadFile(input string) (*core.Mesh, error) {
	ext := strings.ToLower(filepath.Ext(input))
	switch ext {
	case ".ply":
		return parsers.LoadPLY(input)
		// Implement other file types here
		// case ".obj":
		// 	return parsers.LoadOBJ(input)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func replaceExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	return filename[:len(filename)-len(ext)] + newExt
}

func printUsage() {
	fmt.Println("Usage: odx <command> [options] [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  convert [options] <input-file>  Convert 3D model file to ODX format")
	fmt.Println("  version                         Show version information")
	fmt.Println("  help                            Show this help message")
	fmt.Println("\nConvert options:")
	fmt.Println("  -o string                       Output file path (default: input file path with .odx extension)")
}
