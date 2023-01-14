package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/canack/sentient/pkg/tools"
)

func main() {
	var davinci = tools.Davinci{
		MaxTokens:   300,
		Temperature: 0.05,
	}

	chat := tools.NewChatBot(&davinci)
	chat.Setup(getTokenFromEnv())

	log.Printf("Testing connection to OpenAI API...")
	if err := chat.TestConnection(); err != nil {
		log.Printf("chatbot connection failed: %v", err)
		os.Exit(1)
	}
	log.Printf("Connection successful!")
	fmt.Println("----------------------------------------")

	// To speak with the bot, type in a message and press enter.
	for {
		fmt.Printf("Write something: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)
		message := davinci.Query(text)
		fmt.Println(message.Pretty())
		fmt.Println("------")
	}

}

func getTokenFromEnv() string {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		fmt.Printf("OPENAI_API_KEY not set in environment\n")
		os.Exit(1)
	}
	return key
}
