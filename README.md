<a href="https://github.com/ztkent/moki/tags"><img src="https://img.shields.io/github/v/tag/ztkent/moki.svg" alt="Latest Release"></a>
<a href="https://github.com/ztkent/moki/actions"><img src="https://github.com/ztkent/moki/actions/workflows/build.yml/badge.svg?branch=main" alt="Build Status"></a>

# <img width="40" alt="logo_moki" src="https://github.com/Ztkent/moki/assets/7357311/f1dfb864-3c20-4384-898b-1acc4bb7c92f"> Moki

An AI assistant for the command line.  

Tuned to assist with developer tasks like finding files, installing packages, and git.   
Conversation mode can explain code snippets, generate unit tests, and scaffold new projects.

## Usage
- Install moki:  
  ```bash
  go install github.com/Ztkent/moki/cmd/moki@latest
  ```
  
- Set your API key as an environment variable:
  ```bash
  export OPENAI_API_KEY=<your key>
  export ANYSCALE_API_KEY=<your key>
  ```

- Run the assistant:
  ```bash
  # Ask the assistant a question
  moki [your message]
  cat moki.go | moki [tell me about this code]

  # Provide additional context
  moki [tell me about this code]    -file:moki.go
  moki [tell me about this project] -url:https://github.com/Ztkent/moki

  # Start a conversation with the assistant
  moki -c
  moki -m=turbo -c -max-tokens=100000 -t=0.5
  ```

## Example
https://github.com/Ztkent/moki/assets/7357311/52cb7637-39b8-4b49-8bf3-3875ab124c56


## Configuration
- There are two options for the API provider:  
  - OpenAI (https://platform.openai.com/docs/overview)  
  - Anyscale (https://www.anyscale.com/endpoints)  
```
Flags:
  -c:                        Start a conversation with Moki
  -llm [openai, anyscale]:   Set the LLM Provider
  -m:                        Set the model to use for the LLM response
  -max-tokens:               Set the maximum number of tokens to generate
  -t [0.0-1.0]:              Set the temperature for the LLM response
  -d:                        Show debug logging

Model Options:
  - OpenAI:
    - gpt-3.5-turbo, aka: turbo35
    - gpt-4-turbo-preview, aka: turbopreview
    - gpt-4-turbo, aka: turbo
  - Anyscale:
    - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
    - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
    - codellama/CodeLlama-70b-Instruct-hf, aka: cl70b
```

#### Conversation
The assistant can be used in conversation mode.  
This allows the assistant to generate more in-depth responses.
```bash
moki -c
```

#### API Provider
By default the assistant will use OpenAI. To use Anyscale, run the assistant with a flag. 
```bash
moki -llm=openai
moki -llm=anyscale 
```

#### Model
Depending on the LLM Provider selected, different models are available.  
By default the OpenAI API uses `gpt-4-turbo`, and OpenAI uses `Mistral-8x7b`.
```bash
moki -m=m8x7b
```

#### Token Limit
Tokens cost money.   
By default the assistant will limit any conversation to 100k tokens.
```bash
moki -max-tokens=100000
```

#### Temperature
The temperature of an LLM response is a measure of randomness.   
The value float between 0 and 1. By default the temperature is 0.2
```bash
moki -t=0.5
```
