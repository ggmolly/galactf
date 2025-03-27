package types

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
)

type Message struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type   string `json:"type,omitempty"`
	Text   *Text   `json:"text,omitempty"`
	Fields []Text `json:"fields,omitempty"`
}

type Text struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}


func SendSlackWebhook(uri string, message *Message) {
	messageJson, err := sonic.Marshal(message)
    if err != nil {
        log.Printf("[!] Failed to marshal slack message: %v", err)
        return
    }
	req, err := http.NewRequest("POST",
        os.Getenv("SLACK_WEBHOOK_URI"),
        bytes.NewReader(messageJson))
    if err != nil {
        log.Printf("[!] Failed to create HTTP request: %v", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("[!] Slack webhook failed to post message: %v", err)
        return
    }

    log.Println(string(messageJson))
    defer res.Body.Close()
}