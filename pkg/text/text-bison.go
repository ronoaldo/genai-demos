package text

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

// ModelVersion is the currently used model version for text generation
// used by the prediction API calls.
const ModelVersion = "text-bison@001"

// Parameters are model parameters that can be used by the Generative AI
// models on Vertex AI.
type Parameters struct {
	Temperature    float64
	TopP           float64
	TopK           int
	MaxTokens      int
	CandidateCount int
}

// DefaultParameters are the parameters used by default in the Vertex
// Generative AI Studio.
var DefaultParameters = Parameters{
	Temperature:    0.2,
	TopP:           0.8,
	TopK:           40,
	MaxTokens:      1024,
	CandidateCount: 1,
}

// MoreDeterministic are suggested parameters to experiment with.
// According to the documentation, they may generate results
// that are more deterministic.
//
// Warning: these are just suggestions and not a specific recommendation
// from Google. Addapt these to your use case.
var MoreDeterministic = Parameters{
	Temperature:    0.0,
	TopK:           1,
	TopP:           0.8,
	MaxTokens:      1024,
	CandidateCount: 1,
}

// MoreCreative are suggested parameters to experiment with.
// According to the documentation, they may generate results
// that are more creative.
//
// Warning: these are just suggestions and not a specific recommendation
// from Google. Addapt these to your use case.
var MoreCreative = Parameters{
	Temperature:    1.0,
	TopK:           40,
	TopP:           1.0,
	MaxTokens:      1024,
	CandidateCount: 1,
}

// Citation describes a citation reference when the model detects that
// one is needed.
type Citation struct {
	StartIndex      int    `json:"startIndex,omitempty"`
	EndIndex        int    `json:"endIndex,omitempty"`
	URL             string `json:"url,omitempty"`
	Title           string `json:"title,omitempty"`
	License         string `json:"license,omitempty"`
	PublicationDate string `json:"publicationDate,omitempty"`
}

// CitationMetadata holds the list of citations if any.
type CitationMetadata struct {
	Citations []Citation `json:"citations,omitempty"`
}

// SafetyAttributes holds the detailed information of safety attributes
// returned by the model, as part of the Responsible AI actions taken
// by Google to signal if the output is harmful.
type SafetyAttributes struct {
	Blocked    bool      `json:"blocked,omitempty"`
	Categories []string  `json:"categories,omitempty"`
	Scores     []float64 `json:"scores,omitempty"`
}

// Prediction represents the returned prediction content and metadata
// from the Generative AI model on Vertex AI.
type Prediction struct {
	Content          string           `json:"content"`
	CitationMetadata CitationMetadata `json:"citationMetadata,omitempty"`
	SafetyAttributes SafetyAttributes `json:"safetyAttributes,omitempty"`
}

// TokenCountMetadata is a helper struct to encode the resulting
// metadata about billable characters and tokens in the API call.
type TokenCountMetadata struct {
	TotalBillableCharacters int `json:"totalBillableCharacters"`
	TotalTokens             int `json:"totalTokens"`
}

// TokenMetadata is a helper struct to encode the resulting
// metadata about the input/output tokens.
// {"tokenMetadata":{"inputTokenCount":{"totalBillableCharacters":15,"totalTokens":5},"outputTokenCount":{"totalBillableCharacters":191,"totalTokens":52}}}
type TokenMetadata struct {
	InputTokenCount  TokenCountMetadata `json:"inputTokenCount"`
	OutputTokenCount TokenCountMetadata `json:"outputTokenCount"`
}

// Response contains the returned response by the API call of a
// Generative AI model.
type Response struct {
	Predictions []Prediction  `json:"predictions,omitempty"`
	Metadata    TokenMetadata `json:"tokenMetadata,omitempty"`
}

func (r Response) String() string {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

// TextClient is a helper to setup a Generative AI client for
// text generation.
type TextClient struct {
	projectID string
	debugFlag bool
}

// NewClient initializes a new TextClient using the provided projectID.
func NewClient(projectID string) *TextClient {
	return &TextClient{projectID: projectID}
}

// GenerateText calls the Vertex AI text-bison model to generate a new text.
//
// `promptContext` is used as a template for fmt.Sprintf togheter with `prompt`,
// allowing one to define the structured prompt only once. It can be an empty
// string, and if it has no '%s' formatting, it will be a prefix added before prompt.
//
// `params` are the parameters that will be passed to the model, like the temperature,
// top-k or top-p. Note that the default struct values for them are not always what you
// want: it is a good idea to set all parameters explicitly.
//
// The returned Response will contain the list of predictions as well as any metadata
// returned by the call.
func (t *TextClient) GenerateText(ctx context.Context, promptContext, prompt string, params Parameters) (response *Response, err error) {
	// Compile the promptContext with the prompt, allowing for empty context and
	// no formatting strings to be properly used.
	compiledPrompt := fmt.Sprintf(promptContext, prompt)
	if !strings.Contains(promptContext, "%s") {
		compiledPrompt = promptContext + " " + prompt
	}
	// Preparing the request data, using the structpb.Value as a
	// conteiner for the input. This will use the gRPC APIs.
	instance, err := structpb.NewValue(map[string]interface{}{
		"prompt": compiledPrompt,
	})
	if err != nil {
		return nil, err
	}
	parameters, err := structpb.NewValue(map[string]interface{}{
		"temperature":     params.Temperature,
		"maxOutputTokens": params.MaxTokens,
		"topP":            params.TopP,
		"topK":            params.TopK,
		"candidateCount":  params.CandidateCount,
	})
	if err != nil {
		return nil, err
	}
	// Creating the protobuff request to send call the model prediction.
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/%s/models/%s", t.projectID, "us-central1", "google", ModelVersion)
	req := &aiplatformpb.PredictRequest{
		Endpoint:   endpoint,
		Instances:  []*structpb.Value{instance},
		Parameters: parameters,
	}
	t.debug("Sending request => %v", req)

	// Connecting to the desired server
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"))
	if err != nil {
		return nil, err
	}

	// Actually makes the call
	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, err
	}
	t.debug("Got Response => %v", resp)

	// Decoding the response with the help of encoding/json
	r := &Response{}
	b, err := resp.Metadata.MarshalJSON()
	if err != nil {
		return nil, err
	}
	t.debug("Parsing resp.Metadata => %v", string(b))
	if err = json.Unmarshal(b, r); err != nil {
		return nil, err
	}

	for i := range resp.Predictions {
		m := resp.Predictions[i].GetStructValue().AsMap()
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		t.debug("Parsing resp.Prediction[i] => %v", string(b))
		p := Prediction{}
		if err = json.Unmarshal(b, &p); err != nil {
			return nil, err
		}
		r.Predictions = append(r.Predictions, p)
	}
	return r, nil
}

// EnableDebug activates extra messages printed to stderr for debugging.
func (t *TextClient) Debug(enable bool) {
	t.debugFlag = enable
}

// debug is a helper function to print debug messages if debug flag is on.
func (t *TextClient) debug(msg string, v any) {
	if !t.debugFlag {
		return
	}
	log.Printf("DEBUG: "+msg, v)
}
