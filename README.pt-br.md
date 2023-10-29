# Demonstrações de IA generativa da Vertex AI

Coleção de apps de demonstração feitos com as APIs da Vertex AI
para IA generativas.

## Geração de texto

Alguns exemplos de aplicativos CLI que usam os modelos de geração de texto
para mostrar aplicações de IA generativa.

Você precisa de um projeto do Google Cloud configurado previamente.
O projeto precisa ter o faturamento ativo e ter as APIs necessárias
ativas. Depois de criar o projeto e ativar as APIs, instale a SDK
do Google Cloud par facilitar o processo de autenticação:

    gcloud auth application-default login
    gcloud auth application-default set-quota-project *your-project-id*

### cmd/textbison

`cmd/textbison` é um aplicativo CLI simples que chama o
modelo `text-bison@001` com um prompt de texto e retorna
os dados de resposta como JSON. Isso pode então ser analisado,
como com o programa `jq`.

Instalando:

    go install github.com/ronoaldo/genai-demos/cmd/textbison@latest

Exemplo:

    textbison "descrever IA generativa"

Você pode ver todas as opções disponíveis para que possam ser passadas com
`textbison --help`. O programa utilizará as configurações padrão de
autenticação do Google (Google Default Application Credentials).

### cmd/linux-guru

`cmd/linux-guru` é uma ferramenta CLI em português brasileiro para ajudá-lo
a aprender mais sobre Linux.

Instalando:

    go install github.com/ronoaldo/genai-demos/cmd/linux-guru@latest

Exemplos:

    linux-guru quem criou o Linux?
    linux-guru como fazer backup compactado da minha pasta pessoal?

Você pode ver todas as opções disponíveis para que possam ser passadas com
`linux-guru --help`. O programa utilizará as configurações padrão de
autenticação do Google (Google Default Application Credentials).

### cmd/log-guru

`cmd/log-guru` é uma ferramenta CLI em português brasileiro para ajudá-lo
entender o significado dos registros estruturados do Google Cloud.

Instalando:

    go install github.com/ronoaldo/genai-demos/cmd/log-guru@latest

Esta ferramenta espera que os dados JSON de registro sejam enviados para
a entrada padrão e imprimirá a resposta na saída padrão.

Exemplo de uso para explicar os cinco registros mais recentes do Google Cloud:

    gcloud logging read "resource.type=audited_resource" --limit 5 --format=json | log-guru

Você pode ver todas as opções disponíveis para que possam ser passadas com
`log-guru --help`. O programa utilizará as configurações padrão de
autenticação do Google (Google Default Application Credentials).