package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/canack/sentient/pkg/models"
)

type Davinci struct {
	MaxTokens   int
	Temperature float64

	token  string
	apiURL string

	client davinciClient
	data   davinciData
}

type davinciClient struct {
	request     *models.DavinciRequest
	response    *models.DavinciResponse
	httpClient  *http.Client
	httpRequest *http.Request
}

type davinciData struct {
	messages *[]string // Q: and A:
}

func (d *Davinci) Setup(token string) error {
	*d = Davinci{
		token:  token,
		apiURL: "https://api.openai.com/v1/completions",
		client: davinciClient{
			httpClient: &http.Client{
				Timeout: time.Second * 15,
			},
			request: &models.DavinciRequest{
				Model:            "text-davinci-003",
				Prompt:           "",
				MaxTokens:        d.MaxTokens,
				Temperature:      d.Temperature,
				TopP:             1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
			},
		},
		data: davinciData{
			messages: &[]string{},
		},
	}
	return nil
}

func (d *Davinci) TestConnection() error {
	oldPromt := d.client.request.Prompt
	d.client.request.Prompt = "Ping"
	defer func() {
		d.client.request.Prompt = oldPromt
	}()

	_, err := d.do()
	if err != nil {
		return err
	}

	return nil
}

func (d *Davinci) Query(message string) string {
	d.setNewPrompt(message, true)
	d.client.request.Prompt = d.getPrompt()

	raw, err := d.do()
	if err != nil {
		log.Printf("Error executing request: %v\n", err)
		return ""
	}

	var custom models.DavinciResponse
	if err := json.Unmarshal(raw, &custom); err != nil {
		log.Printf("Error unmarshalling response: %v\n", err)
		return ""
	}

	answer := custom.Choices[0].Text

	d.setNewPrompt(answer, false)

	return answer
}

func (d *Davinci) setNewPrompt(message string, fromUser bool) {
	var messages = d.data.messages
	var msg string

	if fromUser {
		msg = "Q: " + message + "\n"
	} else {
		trimAnswer := strings.TrimLeft(message, "\n")
		msg = trimAnswer + "\n"
	}

	*messages = append(*messages, msg)
}

func (d *Davinci) getPrompt() string {
	var prompt string

	for _, msg := range *d.data.messages {
		prompt += msg + "\n"
	}

	return prompt
}

func (d *Davinci) setRequest() error {
	jsonBody, err := json.Marshal(d.client.request)
	if err != nil {
		log.Printf("Error marshalling request body: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", d.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating new request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.token)

	d.client.httpRequest = req
	return nil
}

func (d *Davinci) do() ([]byte, error) {
	if err := d.setRequest(); err != nil {
		log.Printf("Error setting request: %v\n", err)
		return nil, err
	}

	response, err := d.client.httpClient.Do(d.client.httpRequest)
	if err != nil {
		log.Printf("Error executing request: %v\n", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %v\n", response.Status)
		return nil, errors.New("an error occurred while executing the request")
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	return bytes, nil
}
