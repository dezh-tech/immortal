package repository

import (
	"github.com/dezh-tech/immortal/infrastructure/database"
	grpcclient "github.com/dezh-tech/immortal/infrastructure/grpc_client"
	"github.com/dezh-tech/immortal/types"
)

type Handler struct {
	db     *database.Database
	grpc   grpcclient.IClient
	config Config
}

func New(cfg Config, db *database.Database, grpc grpcclient.IClient) *Handler {
	return &Handler{
		db:     db,
		config: cfg,
		grpc:   grpc,
	}
}

func getCollectionName(k types.Kind) string {
	collName, ok := types.KindToName[k]
	if ok {
		return collName
	}

	if k >= 9000 && k <= 9030 {
		return "groups"
	}

	if k >= 1630 && k <= 1633 {
		return "status"
	}

	if k >= 39000 && k <= 39009 {
		return "groups_metadata"
	}

	if k >= 5000 && k <= 5999 || k >= 6000 && k <= 6999 || k == 7000 {
		return "dvm"
	}

	return "unknown"
}
