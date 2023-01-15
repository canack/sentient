# Sentient

#### Warning: This library is still in development and is not yet ready for production use.

---

### Mini library which works with OpenAI's models

### Please check the mini [example](cmd/sentient/main.go)
````shell
git clone https://github.com/canack/sentient.git

cd sentient/cmd/sentient

OPENAI_API_KEY="API_KEY" go run .
````


### Usage/Demo:
````shell
go get github.com/canack/sentient
````

````go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/canack/sentient/pkg/tools"
)

func main() {
	var davinci = tools.Davinci{
		MaxTokens:   300,
		Temperature: 0.05,
	}

	chat := tools.NewChatBot(&davinci)
	chat.Setup("OPENAI_API_KEY")

	log.Println("Testing connection to OpenAI API...")
	if err := chat.TestConnection(); err != nil {
		log.Printf("chatbot connection failed: %v\n", err)
		os.Exit(1)
	}
	log.Println("Connection successful!")

	message := davinci.Query("Hi, how are you?")
	fmt.Println(message.Pretty())
}
````

### License: MIT
