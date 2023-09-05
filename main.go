package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIResponse struct {
	URL     string          `json:"url"`
	Content json.RawMessage `json:"content"`
	Time    string          `json:"time"`
	Error   string          `json:"error,omitempty"`
}

func fetchAPI(url string, resultChan chan<- APIResponse) {
	startTime := time.Now()

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		resultChan <- APIResponse{URL: url, Error: fmt.Sprintf("Error: %v", err)}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChan <- APIResponse{URL: url, Error: fmt.Sprintf("Error reading the response body: %v", err)}
		return
	}

	elapsedTime := time.Since(startTime)
	resultChan <- APIResponse{URL: url, Content: json.RawMessage(body), Time: elapsedTime.String()}
}

func main() {
	var cep string
	fmt.Print("Enter the ZIP code: ")
	fmt.Scan(&cep)

	// Format the ZIP code if necessary
	if len(cep) == 8 {
		cep = fmt.Sprintf("%s-%s", cep[:5], cep[5:])
	} else if len(cep) != 9 || cep[5] != '-' {
		fmt.Println("Invalid ZIP code format. Please use the format 00000-000.")
		return
	}

	apiURLs := []string{
		"https://cdn.apicep.com/file/apicep/" + cep + ".json",
		"http://viacep.com.br/ws/" + cep + "/json/",
	}

	resultChannel := make(chan APIResponse, len(apiURLs))

	for _, url := range apiURLs {
		go fetchAPI(url, resultChannel)
	}

	var res APIResponse
	select {
	case res = <-resultChannel:
	case <-time.After(1 * time.Second):
		res = APIResponse{URL: "", Error: "Timeout - Response time limit exceeded"}
	}

	result, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(result))
}
