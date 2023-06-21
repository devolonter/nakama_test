package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initializer.RegisterRpc("get_content", GetContentRpc)
	return nil
}

func GetContentRpc(ctx context.Context, logger runtime.Logger, db *sql.DB,
	nk runtime.NakamaModule, payload string) (string, error) {
	return "{}", nil
}
