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
	maxAttempts   = 3               // Maximum number of retry attempts
	retryInterval = 2 * time.Second // Sleep interval between retries
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
	for attempts < maxAttempts {
		resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil && attempts == maxAttempts-1 {
			return fmt.Errorf("error making the RPC request: %w", err)
		} else if err == nil {
			defer resp.Body.Close()
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

		attempts++
		log.Println("Retry Request", attempts)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("max retries reached")
}
