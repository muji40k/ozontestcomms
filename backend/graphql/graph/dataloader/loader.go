package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	usrloader "github.com/muji40k/ozontestcomms/graphql/graph/dataloader/user"
	"github.com/muji40k/ozontestcomms/graphql/graph/model"
	usrsrv "github.com/muji40k/ozontestcomms/internal/service/interface/user"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	User *dataloadgen.Loader[uuid.UUID, *model.User]
}

func NewLoaders(user usrsrv.Service, d time.Duration) *Loaders {
	return &Loaders{
		User: dataloadgen.NewLoader(
			usrloader.New(user),
			dataloadgen.WithWait(d),
		),
	}
}

func Middleware(ctor func() *Loaders, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, ctor()))
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) (*Loaders, bool) {
	v, ok := ctx.Value(loadersKey).(*Loaders)
	return v, ok
}

