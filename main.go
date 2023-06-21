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

type ReadFile func(filename string) (*os.File, error)

const (
	CodeInvalidArgument = 3
	CodeNotFound        = 5
	CodeInternalError   = 13
)

const (
	DefaultContentType    = "core"
	DefaultContentVersion = "1.0.0"
)

const DataDir = "data"

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

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS system_logs (
			id INT PRIMARY KEY DEFAULT unique_rowid(),
			user_id UUID,
			payload STRING,
			time TIMESTAMP DEFAULT current_timestamp()
		)`)
	if err != nil {
		return err
	}

	return nil
}

func GetContentRpc(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, payload string) (string, error) {
	_, err := db.ExecContext(ctx, "INSERT INTO system_logs (user_id, payload) VALUES ($1, $2)",
		ctx.Value(runtime.RUNTIME_CTX_USER_ID), payload)

	if err != nil {
		logger.Error("Failed to insert into system_logs: %s", err.Error())
	}

	in, err := bindPayload(payload)
	if err != nil {
		logger.Error("Failed to bind payload: %s", err.Error())
		return "", ErrJsonUnmarshal
	}

	out, err := getContent(in, nk.ReadFile)
	if err != nil {
		logger.Error("Failed to read file: %s", err.Error())
		return "", ErrorFileNotFound
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

func getContent(in GetContentInPayload, readFile ReadFile) (GetContentOutPayload, error) {
	out := GetContentOutPayload{
		GetContentInPayload: in,
	}

	file, err := readFile(filepath.Join(DataDir, in.Type, in.Version+".json"))
	if err != nil {
		return out, err
	}

	defer file.Close()

	reader := io.Reader(file)
	data, err := io.ReadAll(reader)

	if err != nil {
		return out, nil
	}

	hashBytes := sha256.Sum256(data)
	hash := hex.EncodeToString(hashBytes[:])

	if in.Hash != nil && *in.Hash == hash {
		var fileJson map[string]string
		err := json.Unmarshal(data, &fileJson)

		if err != nil {
			return out, nil
		}

		if content, ok := fileJson["content"]; ok {
			out.Content = &content
		}
	}

	return out, nil
}
