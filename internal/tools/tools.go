package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Ztkent/bash-gpt/internal/prompts"
	aiclient "github.com/Ztkent/go-openai-extended"
	"github.com/rs/zerolog/log"
)

func StartConversationCLI(client *aiclient.Client, conv *aiclient.Conversation) error {
	var exitCommands = []string{"exit", "quit", "bye", ":q", "end", "q"}
	var helpCommands = []string{"help", "?"}

	// This is the maximum conversation time
	thirtyMin, cancel0 := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel0()

	oneMin, cancel := context.WithTimeout(thirtyMin, time.Minute*1)
	defer cancel()

	// Start the chat with a fresh conversation, and get the system greeting
	introChat, err := client.SendCompletionRequest(oneMin, aiclient.NewConversation(prompts.BashGPTPrompt, 0, 0), "We're starting a conversation. Introduce yourself.")
	if err != nil {
		return err
	}
	fmt.Println("BashGPT: " + introChat)

	// Lets start a conversation with the user via CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Request: ")
		// Ask for the user's input
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Check if the user wants to exit
		if strings.Contains(strings.Join(exitCommands, "|"), strings.ToLower(userInput)) {
			break
		} else if strings.Contains(strings.Join(helpCommands, "|"), strings.ToLower(userInput)) {
			fmt.Println("--------------------------------------------------")
			fmt.Println("bashgpt: ")
			fmt.Println("    Type 'exit', 'quit', or 'bye' to end the conversation.")
			fmt.Println("    Type your message to continue the conversation.")
			continue
		}

		// Check if the user provided a message
		if len(userInput) == 0 {
			fmt.Println("Please provide a message to continue the conversation.")
			continue
		}

		// Send the user's input to the LLM 🤖, wait at most 1 minute.
		oneMin, cancel = context.WithTimeout(thirtyMin, time.Minute*1)
		defer cancel()
		responseChan, errChan := make(chan string), make(chan error)
		go client.SendStreamRequest(oneMin, conv, userInput, responseChan, errChan)
		fmt.Print("BashGPT: ")

		// Read the response from the channel as it is streamed
		done := false
		for !done {
			select {
			case response, ok := <-responseChan:
				if !ok {
					// Request channel closed
					done = true
					break
				}
				fmt.Print(response)
			case err := <-errChan:
				if err != nil {
					return err
				}
			}
		}
		fmt.Println()
	}
	return nil
}

// Log the results of a fresh chat stream
func LogNewChatStream(client *aiclient.Client, conv *aiclient.Conversation, chatPrompt string) error {
	oneMin, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// Start the chat with a fresh conversation, and the users prompt
	responseChan, errChan := make(chan string), make(chan error)
	log.Debug().Msg(fmt.Sprintf("prompt: " + chatPrompt))
	go client.SendStreamRequest(oneMin, conv, chatPrompt, responseChan, errChan)
	// Read the response from the channel as it is streamed
	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// Request channel closed
				fmt.Println()
				return nil
			}
			fmt.Print(response)
		case err := <-errChan:
			fmt.Println()
			return err
		}
	}
}
