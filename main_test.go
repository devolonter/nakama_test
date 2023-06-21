package main

import (
	"fmt"
	"os"
	"testing"
)

var readFile = os.Open

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestEmptyInPayload(t *testing.T) {
	payload := `{}`
	result, err := bindPayload(payload)

	if err != nil {
		t.Errorf("bindPayload returned an error: %v", err)
	}

	if result.Type != DefaultContentType {
		t.Errorf("bindPayload returned an unexpected type: %v", result.Type)
	}

	if result.Version != DefaultContentVersion {
		t.Errorf("bindPayload returned an unexpected version: %v", result.Version)
	}

	if result.Hash != nil {
		t.Errorf("bindPayload returned an unexpected hash: %v", result.Hash)
	}
}

func TestNonEmptyPayload(t *testing.T) {
	const (
		ContentType    = "text"
		ContentVersion = "2.0.0"
		ContentHash    = "1234567890"
	)

	payload := fmt.Sprintf(`{"type":"%s","version":"%s","hash":"%s"}`, ContentType, ContentVersion, ContentHash)
	result, err := bindPayload(payload)
	if err != nil {
		t.Errorf("bindPayload returned an error: %v", err)
	}

	if result.Type != ContentType {
		t.Errorf("bindPayload returned an unexpected type: %v", result.Type)
	}

	if result.Version != ContentVersion {
		t.Errorf("bindPayload returned an unexpected version: %v", result.Version)
	}

	if result.Hash == nil {
		t.Errorf("bindPayload returned an unexpected hash: %v", result.Hash)
	} else if *result.Hash != ContentHash {
		t.Errorf("bindPayload returned an unexpected hash: %v", result.Hash)
	}
}

func TestGetContentWithoutHash(t *testing.T) {
	payload := GetContentInPayload{
		Type:    DefaultContentType,
		Version: DefaultContentVersion,
		Hash:    nil,
	}

	result, err := getContent(payload, readFile)
	if err != nil {
		t.Errorf("getContent returned an error: %v", err)
	}

	if result.Content != nil {
		t.Errorf("getContent returned an unexpected content: %v", result.Content)
	}
}

func TestGetContentWithInvalidHash(t *testing.T) {
	hash := "1234567890"
	payload := GetContentInPayload{
		Type:    DefaultContentType,
		Version: DefaultContentVersion,
		Hash:    &hash,
	}

	result, err := getContent(payload, readFile)
	if err != nil {
		t.Errorf("getContent returned an error: %v", err)
	}

	if result.Content != nil {
		t.Errorf("getContent returned an unexpected content: %v", result.Content)
	}
}

func TestGetContentNotExists(t *testing.T) {
	hash := "14a8517ea074b2906ebbdcc426acdbd2bf24b92bb27aa14ffaf00527c21c68ac"
	payload := GetContentInPayload{
		Type:    DefaultContentType,
		Version: "2.0.0",
		Hash:    &hash,
	}

	_, err := getContent(payload, readFile)
	if err == nil {
		t.Errorf("getContent did not return an error")
	}
}

func TestGetContentWithValidHash(t *testing.T) {
	hash := "14a8517ea074b2906ebbdcc426acdbd2bf24b92bb27aa14ffaf00527c21c68ac"
	payload := GetContentInPayload{
		Type:    DefaultContentType,
		Version: DefaultContentVersion,
		Hash:    &hash,
	}

	result, err := getContent(payload, readFile)
	if err != nil {
		t.Errorf("getContent returned an error: %v", err)
	}

	if result.Content == nil {
		t.Errorf("getContent returned an unexpected content: %v", result.Content)
	} else if *result.Content != "Hello World!" {
		t.Errorf("getContent returned an unexpected content: %v", result.Content)
	}
}
