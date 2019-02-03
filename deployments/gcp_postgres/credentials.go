package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConnectionConfiguration is a struct type for loading in our db connection information
type ConnectionConfiguration struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// GetConnectionConfiguration returns the crentials config for db connection
func GetConnectionConfiguration() ConnectionConfiguration {
	file, _ := os.Open("connection.config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := ConnectionConfiguration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
