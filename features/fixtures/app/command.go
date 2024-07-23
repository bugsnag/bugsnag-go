package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const DEFAULT_MAZE_ADDRESS = "http://localhost:9339"

type Command struct {
	Action           string `json:"action,omitempty"`
	ScenarioName     string `json:"scenario_name,omitempty"`
	APIKey           string `json:"api_key,omitempty"`
	NotifyEndpoint   string `json:"notify_endpoint,omitempty"`
	SessionsEndpoint string `json:"sessions_endpoint,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	RunUUID          string `json:"run_uuid,omitempty"`
}

func GetCommand(mazeAddress string) Command {
	var command Command
	mazeURL := fmt.Sprintf("%+v/command", mazeAddress)
	client := http.Client{Timeout: 2 * time.Second}
	res, err := client.Get(mazeURL)
	if err != nil {
		fmt.Printf("[Bugsnag] Error while receiving command: %+v\n", err)
		return command
	}

	if res != nil {
		err = json.NewDecoder(res.Body).Decode(&command)
		res.Body.Close()
		if err != nil {
			fmt.Printf("[Bugsnag] Error while decoding command: %+v\n", err)
			return command
		}
	}

	return command
}
