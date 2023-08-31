package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChan <- APIResponse{URL: url, Error: fmt.Sprintf("Error reading the response body: %v", err)}
		return
	}

	elapsedTime := time.Since(startTime)
	resultChan <- APIResponse{URL: url, Content: json.RawMessage(body), Time: elapsedTime.String()}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scan(&cep)

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
