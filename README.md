# Vertex AI Generative AI Demos

Collection of demonstratinos of apps made with Vertex AI
Generative AI APIs.

## Text generation

Some sample CLI apps that use the text generation models
to showcase applications of generative AI.

You need a Google Cloud project to be configured previously.
The project need to have billing enabled and to have the required
APIs active. After you have created the project and enabled the APIs,
install the Google Cloud SDK to make easier to authenticate:

    gcloud auth application-default login
    gcloud auth application-default set-quota-project *your-project-id*

### cmd/textbison

`cmd/textbison` is a simple CLI application that calls the
`text-bison` model with a text prompt and returns the full
response data as JSON. This can then be parsed, like with
`jq`.

Installing:

    go install github.com/ronoaldo/genai-demos/cmd/textbison@latest

Sample:

    textbison "describe generative ai"

You can see all available options for that can be passed with
`textbison --help`. The program will use the Google Default
Application credentials algorithm to authenticate.

### cmd/linux-guru

`cmd/linux-guru` is a Brazilian Portuguese CLI tool to help you
learn more about Linux.

Installing:

    go install github.com/ronoaldo/genai-demos/cmd/linux-guru@latest

Samples:

    linux-guru quem criou o Linux?
    linux-guru como fazer backup compactado da minha pasta pessoal?

You can see all available options for that can be passed with
`linux-guru --help`. The program will use the Google Default
Application credentials algorithm to authenticate.

### cmd/log-guru

`cmd/log-guru` is a Brazilian Portuguese CLI tool to help you
understand the meaning of Google Cloud structured logs.

Installing:

    go install github.com/ronoaldo/genai-demos/cmd/log-guru@latest

This tool expects the logging JSON data to be passed from standard
input and will print the response to standard output.
Sample usage to explain the latest 5 Google Cloud logs:

    gcloud logging read "resource.type=audited_resource" --limit 5 --format=json | log-guru

You can see all available options for that can be passed with
`log-guru --help`. The program will use the Google Default
Application credentials algorithm to authenticate.