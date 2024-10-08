package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	maxAttempts   = 5               // Maximum number of retry attempts
	retryInterval = 5 * time.Second // Sleep interval between retries
)

func parseHexToInt(hexString string) int {
	var result int
	fmt.Sscanf(hexString, "0x%x", &result)
	return result
}

type RPCResponse interface {
	Unmarshal([]byte) error
}

func sendRequest(rpcURL string, requestData interface{}, responseType RPCResponse) error {
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error encoding request data: %w", err)
	}

	attempts := 0
	for {
		resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			if attempts == maxAttempts-1 {
				return fmt.Errorf("error making the RPC request: %w", err)
			}
			attempts++
			log.Println("Retry Request", attempts)
			time.Sleep(retryInterval)
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			if attempts == maxAttempts-1 {
				return fmt.Errorf("error getting the RPC response, status code %d", resp.StatusCode)
			}
			attempts++
			log.Println("Retry Request", attempts)
			time.Sleep(retryInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}

		err = responseType.Unmarshal(body)
		if err != nil {
			return fmt.Errorf("error decoding RPC response: %w", err)
		}
		return nil
	}
}
