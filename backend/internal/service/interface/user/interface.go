package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	"github.com/muji40k/ozontestcomms/misc/result"
)

//go:generate mockgen -source=interface.go -destination=../../mock/user/service.go

type Service interface {
	GetUsersById(ctx context.Context, ids ...uuid.UUID) (collection.Collection[result.Result[models.User]], error)
}

