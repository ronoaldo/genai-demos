package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ronoaldo/genai-demos/pkg/text"
)

var projectID string

func init() {
	flag.StringVar(&projectID, "project",
		os.Getenv("GOOGLE_CLOUD_PROJECT"), "The Google `PROJECT_ID` to be used.")
}

func main() {
	// Parse command line options
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalf("Please provide a prompt in the command line.")
	}
	prompt := strings.Join(flag.Args(), " ")
	params := text.DefaultParameters

	// Print the request attributes used
	log.Printf("Prompt: %#v", prompt)
	log.Printf("Params: %#v", params)
	ctx := context.Background()

	// Call the model to generate text
	model := text.NewClient(projectID)
	resp, err := model.GenerateText(ctx, "%s", prompt, params)
	if err != nil {
		log.Fatalf("error invoking model.GenerateText: %v", err.Error())
	}

	// Print the full response as JSON to standard output
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err = enc.Encode(resp); err != nil {
		log.Fatalf("error formatting the output: %v", err.Error())
	}
}
