package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"github.com/heroiclabs/nakama-common/runtime"
	"io"
	"os"
	"path/filepath"
)

const (
	CodeInvalidArgument = 3
	CodeNotFound        = 5
	CodeInternalError   = 13
)

const (
	DefaultContentType    = "core"
	DefaultContentVersion = "1.0.0"
)

const DataDir = "data/"

var (
	ErrJsonMarshal    = runtime.NewError("failed to marshal json", CodeInvalidArgument)
	ErrJsonUnmarshal  = runtime.NewError("failed to unmarshal json", CodeInternalError)
	ErrorFileNotFound = runtime.NewError("file not found", CodeNotFound)
)

type GetContentInPayload struct {
	Type    string  `json:"type"`
	Version string  `json:"version"`
	Hash    *string `json:"hash"`
}

type GetContentOutPayload struct {
	GetContentInPayload
	Content *string `json:"content"`
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, initializer runtime.Initializer) error {
	err := initializer.RegisterRpc("get_content", GetContentRpc)
	if err != nil {
		return err
	}

	return nil
}

func GetContentRpc(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, payload string) (string, error) {
	in, err := bindPayload(payload)
	if err != nil {
		logger.Error("Failed to bind payload: %s", err.Error())
		return "", ErrJsonUnmarshal
	}

	file, err := nk.ReadFile(filepath.Join(DataDir, in.Type, in.Version+".json"))

	if err != nil {
		logger.Error("Failed to read file: %s", err.Error())
		return "", ErrorFileNotFound
	}

	defer file.Close()
	out := getContent(in, file)

	result, err := json.Marshal(out)
	if err != nil {
		logger.Error("Failed to marshal json: %s", err.Error())
		return "", ErrJsonMarshal
	}

	return string(result), nil
}

func bindPayload(payload string) (GetContentInPayload, error) {
	in := GetContentInPayload{
		Type:    DefaultContentType,
		Version: DefaultContentVersion,
	}

	err := json.Unmarshal([]byte(payload), &in)
	return in, err
}

func getContent(in GetContentInPayload, file *os.File) GetContentOutPayload {
	out := GetContentOutPayload{
		GetContentInPayload: in,
	}

	reader := io.Reader(file)
	data, err := io.ReadAll(reader)

	if err != nil {
		return out
	}

	hashBytes := sha256.Sum256(data)
	hash := hex.EncodeToString(hashBytes[:])

	if in.Hash != nil && *in.Hash == hash {
		var fileJson map[string]string
		err := json.Unmarshal(data, &fileJson)

		if err != nil {
			return out
		}

		if content, ok := fileJson["content"]; ok {
			out.Content = &content
		}
	}

	return out
}
