# Vertex AI Generative AI Demos

Collection of demonstratinos of apps made with Vertex AI
Generative AI APIs.

## Text generation

`cmd/textbison` is a simple CLI application that calls the
`text-bison` model with a text prompt and returns the full
response data as JSON. This can then be parsed, like with
`jq`.

Installing:

    go install github.com/ronoaldo/genai-demos/cmd/textbison@latest

Sample:

    textbison describe generative ai

You can see all available options for that can be passed with
`textbison --help`. The program will use the Google Default
Application credentials algorithm to authenticate.