package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	CodeOk              = 0
	CodeInvalidArgument = 3
)

const (
	DefaultContentType    = "core"
	DefaultContentVersion = "1.0.0"
)

var (
	ErrJsonMarshal   = runtime.NewError("failed to marshal json", CodeInvalidArgument)
	ErrJsonUnmarshal = runtime.NewError("failed to unmarshal json", CodeInvalidArgument)
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
	initializer.RegisterRpc("get_content", GetContentRpc)
	return nil
}

func GetContentRpc(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, payload string) (string, error) {
	in, err := bindPayload(payload)
	if err != nil {
		logger.Error("Failed to bind payload: %s", err.Error())
		return "", ErrJsonUnmarshal
	}

	out := GetContentOutPayload{
		GetContentInPayload: in,
	}

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
