package text

import (
	"context"
	"flag"
	"strings"
	"testing"
)

func init() {
	flag.StringVar(&projectID, "project", "", "The Google Project ID to be used")
}

var projectID string

var promptContext = `Context: Only answers questions about Information Technology Google Cloud Platform.
For any other questions, answer:  I don't know about this topic.

Question: %s
Answer: `

func TestTextBison(t *testing.T) {
	type args struct {
		prompt string
	}
	tests := []struct {
		name              string
		args              args
		wantGenerated     string
		wantErr           bool
		wantBillableChars int
	}{
		{
			"test off-topic prompt",
			args{"In one word, what is the color of the sky?"},
			"I don't know about this topic.",
			false,
			176 + 25, // input + output
		},
		{
			"test app engine prompt",
			args{"When was Google App Engine launched?"},
			"May 2008",
			false,
			174 + 7,
		},
		{
			"test gsutil prompt",
			args{"What is the command to copy a file to a bucket?"},
			"gsutil cp file.txt gs://bucket/",
			false,
			180 + 28,
		},
	}

	ctx := context.Background()
	textClient := NewClient(projectID)

	l := strings.ToLower
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := textClient.GenerateText(ctx, promptContext, tc.args.prompt, MoreDeterministic)
			if (err != nil) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			t.Logf("\nPrompt: %v\nAnswer: %v", tc.args.prompt, gen.Predictions[0].Content)
			if l(gen.Predictions[0].Content) != l(tc.wantGenerated) {
				t.Errorf("got %#v as answer, want %#v", gen.Predictions[0].Content, tc.wantGenerated)
			}
			billableChars := gen.Metadata.InputTokenCount.TotalBillableCharacters +
				gen.Metadata.OutputTokenCount.TotalBillableCharacters
			if billableChars != tc.wantBillableChars {
				t.Errorf("got %v billable characters, want %v", billableChars, tc.wantBillableChars)
			}
			t.Logf("\nFull response: %s", gen)
		})
	}
}
