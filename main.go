package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/sudipbhandari126/amuse/utils"
)

var (
	url         = "https://github.com/Amuze-Me/knowledge-base/raw/main/kb.csv"
	destination = "/tmp/kb_local.csv"
	etagPath    = "/tmp/kb_etag"
)

var rootCmd = &cobra.Command{
	Use:   "amuse",
	Short: "A simple CLI program",
}

var techCmd = &cobra.Command{
	Use:   "tech",
	Short: "Print a technical message",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a valid tag (tech/debug/arch/trivia)")
			return
		}

		tag := args[0]
		if tag != "tech" && tag != "debug" && tag != "arch" && tag != "trivia" {
			fmt.Println("Invalid tag. Use one of: tech, debug, arch, trivia")
			return
		}

		// Fetch a random record with the specified tag from the CSV file
		record, err := getRandomRecordWithTag(tag)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Print the fetched record
		fmt.Printf("Format: %s\nLink: %s\nTags: %s\n", record[0], record[1], record[2])
	},
}

func init() {
	// Add subcommands to the root command
	rootCmd.AddCommand(techCmd)
	rand.Seed(time.Now().UnixNano())
	must(utils.DownloadIfChanged(destination, etagPath, url))
}

func must(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func getRandomRecordWithTag(tag string) ([]string, error) {
	file, err := os.Open("/tmp/kb_local.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	// Filter records based on the specified tag
	var filteredRecords [][]string
	for _, record := range records {
		if len(record) >= 3 && record[2] == tag {
			filteredRecords = append(filteredRecords, record)
		}
	}

	if len(filteredRecords) == 0 {
		return nil, fmt.Errorf("no records found with tag: %s", tag)
	}

	// Fetch a random record from the filtered records
	randomIndex := rand.Intn(len(filteredRecords))
	return filteredRecords[randomIndex], nil
}
