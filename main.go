package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	url := flag.String("url", "https://jsonplaceholder.typicode.com/todos/1", "the url to be fetched from")
	method := flag.String("method", "GET", "the HTTP method to be used")
	body := flag.String("body", "", "the request body (for POST, PUT, etc.)")
	headers := flag.String("headers", "", "the request headers (format: key1:value1,key2:value2)")

	flag.Parse()

	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Method: %s\n", *method)
	fmt.Printf("Body: %s\n", *body)

	bodyMap := make(map[string]string)
	if *body != "" {
		trimmedBody := strings.Trim(*body, "{}")
		pairs := strings.Split(trimmedBody, ",")
		for _, pair := range pairs {
			keyValue := strings.Split(pair, ":")
			if len(keyValue) == 2 {
				bodyMap[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
			} else {
				fmt.Printf("Invalid key-value pair: %v\n", pair)
				return
			}
		}
	}

	var reqBody []byte
	var err error
	if len(bodyMap) > 0 {
		reqBody, err = json.Marshal(bodyMap)
		if err != nil {
			fmt.Printf("Error marshaling request body: %v\n", err)
			return
		}
	}

	client := &http.Client{}
	var req *http.Request

	switch strings.ToUpper(*method) {
	case "GET":
		req, err = http.NewRequest(http.MethodGet, *url, nil)
	case "POST":
		req, err = http.NewRequest(http.MethodPost, *url, strings.NewReader(string(reqBody)))
	case "PUT":
		req, err = http.NewRequest(http.MethodPut, *url, strings.NewReader(string(reqBody)))
	case "DELETE":
		req, err = http.NewRequest(http.MethodDelete, *url, nil)
	default:
		fmt.Printf("Unsupported HTTP method: %v\n", method)
		return
	}

	if *headers != "" {
		headersMap := make(map[string]string)
		trimmedHeaders := strings.Trim(*headers, "{}")
		headerPairs := strings.Split(trimmedHeaders, ",")
		for _, pair := range headerPairs {
			keyValue := strings.Split(pair, ":")
			if len(keyValue) == 2 {
				headersMap[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
			} else {
				fmt.Printf("Invalid header key-value pair: %v\n", pair)
				return
			}
		}
		for key, value := range headersMap {
			req.Header.Add(key, value)
		}
	}
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Body:")
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading response: %v\n", err)
	}

}
