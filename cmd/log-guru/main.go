package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ronoaldo/genai-demos/pkg/text"
)

var projectID string
var verbose bool
var promptContext = `
Você resume e interpreta a saída de logs estruturados do Google Cloud Logging.
A resposta deve ser curta e objetiva.

Explique em Português o que está acontecendo com base no log em JSON abaixo:

`

func init() {
	flag.StringVar(&projectID, "project",
		os.Getenv("GOOGLE_CLOUD_PROJECT"), "The Google `PROJECT_ID` to be used.")
	flag.BoolVar(&verbose, "v", false, "If the output should be more verbose.")
}

func main() {
	// Parse command line options
	flag.Parse()

	// Parse stdin as the prompt
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	jsonlog := string(b)
	params := text.DefaultParameters
	ctx := context.Background()
	model := text.NewClient(projectID)

	// Call the model to generate text
	if verbose {
		log.Printf("Analisando log: %v", jsonlog)
		model.Debug(true)
	}
	resp, err := model.GenerateText(ctx, promptContext, jsonlog, params)
	if err != nil {
		log.Fatalf("Erro: model.GenerateText: %v", err.Error())
	}

	// Print the full response as JSON to standard output
	generated := resp.Predictions[0]

	if generated.SafetyAttributes.Blocked {
		log.Printf("Detalhes: %#v", generated.SafetyAttributes)
		log.Fatal("Esta resposta foi bloqueada.")
	}

	fmt.Println(generated.Content)
	if len(generated.CitationMetadata.Citations) > 0 {
		fmt.Println("\nReferences:")
		for _, citation := range resp.Predictions[0].CitationMetadata.Citations {
			fmt.Println(citation.Title, " ", citation.URL)
		}
	}
}
