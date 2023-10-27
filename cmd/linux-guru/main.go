package main

import (
	"context"
	"flag"
	"fmt"
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

var promptContext = `Context: apenas responda a perguntas sobre Linux e GNU/Linux.
Para outras perguntas, responda: Não sei sobre este tema, tente outra pergunta.

Pergunta: %s
Resposta: `

var disclaimer = `
+--[Aviso]---------------------------------------+
| Este é um conteúdo gerado por AI.              |
| Revise quaisquer comandos antes de executá-los.|
+------------------------------------------------+

`

func main() {
	// Parse command line options
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalf("Erro: nenhuma pergunta informada na linha de comandos.")
	}
	prompt := strings.Join(flag.Args(), " ")
	params := text.DefaultParameters

	ctx := context.Background()

	// Call the model to generate text
	model := text.NewClient(projectID)
	resp, err := model.GenerateText(ctx, promptContext, prompt, params)
	if err != nil {
		log.Fatalf("Erro: model.GenerateText: %v", err.Error())
	}

	// Print the full response as JSON to standard output
	generated := resp.Predictions[0]

	if generated.SafetyAttributes.Blocked {
		log.Printf("Detalhes: %#v", generated.SafetyAttributes)
		log.Fatal("Esta resposta foi bloqueada.")
	}

	fmt.Println(disclaimer)

	fmt.Println(generated.Content)
	if len(generated.CitationMetadata.Citations) > 0 {
		fmt.Println("\nReferences:")
		for _, citation := range resp.Predictions[0].CitationMetadata.Citations {
			fmt.Println(citation.Title, " ", citation.URL)
		}
	}
}
