package titan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type ErrorCode int

// Request defines a JSON RPC request from the spec
// http://www.jsonrpc.org/specification#request_object
type Request struct {
	Jsonrpc string            `json:"jsonrpc"`
	ID      interface{}       `json:"id,omitempty"`
	Method  string            `json:"method"`
	Params  json.RawMessage   `json:"params"`
	Meta    map[string]string `json:"meta,omitempty"`
}

// Response defines a JSON RPC response from the spec
// http://www.jsonrpc.org/specification#response_object
type Response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	ID      interface{} `json:"id"`
	Error   *respError  `json:"error,omitempty"`
}

type respError struct {
	Code    ErrorCode       `json:"code"`
	Message string          `json:"message"`
	Meta    json.RawMessage `json:"meta,omitempty"`
}

func (e *respError) Error() string {
	if e.Code >= -32768 && e.Code <= -32000 {
		return fmt.Sprintf("RPC error (%d): %s", e.Code, e.Message)
	}
	return e.Message
}

func getRpcURL(baseURL string) string {
	if strings.Contains(baseURL, "/rpc") {
		return baseURL
	}
	return fmt.Sprintf("%s/rpc/v0", baseURL)
}

func handleRequest(client *http.Client, url string, req Request) (*Response, error) {
	b, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Errorf("marshalling request: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, getRpcURL(url), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out Response
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	if out.Error != nil {
		return nil, out.Error
	}

	return &out, nil
}

func decodeResult(data interface{}, result interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, result)
}
