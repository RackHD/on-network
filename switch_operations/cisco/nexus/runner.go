package nexus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Runner struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Params struct {
	Command string `json:"cmd"`
	Version int    `json:"version"`
}

type CommandRunnerBody struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type CopyCommand struct {
	Src string
	Dst string
}

func (nr *Runner) Run(command string, timeout time.Duration) (string, error) {
	endpoint := fmt.Sprintf("http://%s/ins", nr.IP)

	commandParam := Params{command, 1}
	postBody := CommandRunnerBody{"2.0", "cli", commandParam, 1}
	bodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling command: %+v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error making request: %+v", err)
	}

	req.Header.Set("Content-Type", "application/json-rpc")
	req.SetBasicAuth(nr.Username, nr.Password)

	fmt.Println("making request to nxos...")

	client := &http.Client{}
	if timeout != 0 {
		client.Timeout = timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %+v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body), fmt.Errorf("error reading response body: %+v", err)
	}

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("failed to get expected string. status code: %d\nbody: %+v", resp.StatusCode, string(body))
		return string(body), errors.New(errMsg)
	}

	return string(body), nil
}
