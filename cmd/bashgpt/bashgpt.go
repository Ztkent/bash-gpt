package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Ztkent/bash-gpt/internal/prompts"
	"github.com/Ztkent/bash-gpt/internal/tools"
	"github.com/Ztkent/bash-gpt/pkg/aiclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*
Command-line interface for a BashGPT conversation.

Usage:
  bashgpt [your question]

Flags:
  -h:                        Show this message
  -c:                        Start a conversation with BashGPT
  -ai [openai, anyscale]:    Set the LLM Provider
  -m [string]:               Set the model to use for the LLM response
  -max-messages [int]:       Set the maximum conversation context length
  -max-tokens [int]:         Set the maximum number of tokens to generate per response
  -t [0.0-1.0]:              Set the temperature for the LLM response
  -d:                        Show debug logging

  Model Options:
    -openai:
	  - gpt-3.5-turbo, aka: turbo
	-anyscale:
	  - mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
	  - mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
	  - meta-llama/Llama-2-7b-chat-hf, aka: l7b
	  - meta-llama/Llama-2-13b-chat-hf, aka: l13b
	  - meta-llama/Llama-2-70b-chat-hf, aka: l70b
	  - codellama/CodeLlama-34b-Instruct-hf, aka: cl34b
.*/

func main() {
	// Define the flags
	helpFlag := flag.Bool("h", false, "Show this message")
	debugFlag := flag.Bool("d", false, "Show debug logs")
	convFlag := flag.Bool("c", false, "Start a conversation with BashGPT")
	aiFlag := flag.String("llm", "anyscale", "Selct the LLM provider, either OpenAI or Anyscale")
	modelFlag := flag.String("m", "cl34b", "Set the model to use for the LLM response")
	temperatureFlag := flag.Float64("t", 0.2, "Set the temperature for the LLM response")
	maxMessagesFlag := flag.Int("max-messages", 0, "Set the maximum conversation context length")
	maxTokensFlag := flag.Int("max-tokens", 1000, "Set the maximum number of tokens to generate per response")

	// Parse the flags
	flag.Parse()

	// Set Logging level
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Show the help message
	if *helpFlag {
		fmt.Println(
			`
Usage:
	bashgpt [your question]

Flags:
	-h:                        Show this message
	-c:                        Start a conversation with BashGPT
	-llm [openai, anyscale]:   Set the LLM Provider
	-m [string]:               Set the model to use for the LLM response
	-max-messages [int]:       Set the maximum conversation context length
	-max-tokens [int]:         Set the maximum number of tokens to generate per response
	-t [0.0-1.0]:              Set the temperature for the LLM response
	-d:                        Show debug logging

API Keys:
	Set your API keys as environment variables:
		- export OPENAI_API_KEY=<your key>
		- export ANYSCALE_API_KEY=<your key>

Model Options:
	- OpenAI:
		- gpt-3.5-turbo, aka: turbo
	- Anyscale:
		- mistralai/Mistral-7B-Instruct-v0.1, aka: m7b
		- mistralai/Mixtral-8x7B-Instruct-v0.1, aka: m8x7b
		- meta-llama/Llama-2-7b-chat-hf, aka: l7b
		- meta-llama/Llama-2-13b-chat-hf, aka: l13b
		- meta-llama/Llama-2-70b-chat-hf, aka: l70b
		- codellama/CodeLlama-34b-Instruct-hf, aka: cl34b
			`)
		return
	}

	var client *aiclient.Client
	if *aiFlag == "openai" {
		err := aiclient.MustLoadAPIKey(true, false)
		if err != nil {
			fmt.Printf("Failed to load OpenAI API key: %s\n", err)
			return
		}

		//  Connect to the OpenAI Client with the given model
		if model, ok := aiclient.IsOpenAIModel(*modelFlag); ok {
			log.Debug().Msg(fmt.Sprintf("Starting client with OpenAI-%s\n", model))
			client = aiclient.MustConnectOpenAI(model, float32(*temperatureFlag))
		} else {
			// Default to GPT-3.5 Turbo
			log.Debug().Msg(fmt.Sprintf("Starting client with OpenAI-%s\n", aiclient.GPT35Turbo))
			client = aiclient.MustConnectOpenAI(aiclient.GPT35Turbo, float32(*temperatureFlag))
		}
	} else if *aiFlag == "anyscale" {
		err := aiclient.MustLoadAPIKey(false, true)
		if err != nil {
			log.Error().AnErr("Failed to load Anyscale API key", err)
			return
		}

		//  Connect to the Anyscale Client with the given model
		if model, ok := aiclient.IsAnyscaleModel(*modelFlag); ok {
			log.Debug().Msg(fmt.Sprintf("Starting client with Anyscale-%s\n", model))
			client = aiclient.MustConnectAnyscale(model, float32(*temperatureFlag))
		} else {
			// Default to Mistral 7B Instruct
			log.Debug().Msg(fmt.Sprintf("Starting client with Anyscale-%s\n", aiclient.Mistral7BInstruct))
			client = aiclient.MustConnectAnyscale(aiclient.Mistral7BInstruct, float32(*temperatureFlag))
		}
	} else {
		fmt.Println(fmt.Sprintf("Invalid AI provider: %s provided, select either anyscale or openai", *aiFlag))
		return
	}

	conv := aiclient.NewConversation(prompts.BashGPTPrompt, *maxMessagesFlag, *maxTokensFlag)
	if *convFlag {
		// Start a conversation with the BashGPT prompt
		err := tools.StartConversationCLI(client, conv)
		if err != nil {
			fmt.Printf("Failed to start conversation: %s\n", err)
		}
		return
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a question to ask BashGPT")
		return
	}

	// Send a single request to the LLM, return it to the user.
	err := tools.LogNewChatStream(client, conv, strings.Join(flag.Args(), " "))
	if err != nil {
		fmt.Printf("Failed to log new chat stream: %s\n", err)
	}
}
