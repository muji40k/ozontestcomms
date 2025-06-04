package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/graphql/graph/mappers"
	"github.com/muji40k/ozontestcomms/graphql/graph/model"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection/iterator"
	usrsrv "github.com/muji40k/ozontestcomms/internal/service/interface/user"
	"github.com/muji40k/ozontestcomms/misc/result"
)

func New(user usrsrv.Service) func(
	ctx context.Context,
	ids []uuid.UUID,
) ([]*model.User, []error) {
	return func(ctx context.Context, ids []uuid.UUID) ([]*model.User, []error) {
		col, err := user.GetUsersById(ctx, ids...)
		var iter iterator.Iterator[result.Result[models.User]]
		var out []*model.User
		var errs []error

		if nil == err {
			iter, err = col.Get()
		}

		if nil == err {
			out = make([]*model.User, len(ids))
			errs = make([]error, len(ids))
			i := 0

			for res := range iterator.Values(iter) {
				if v, err := res.Unwrap(); nil == err {
					out[i] = mappers.MapUser(&v)
					errs[i] = nil
				} else {
					out[i] = nil
					errs[i] = err
				}
			}
		}

		if nil != err {
			errs = []error{err}
		}

		return out, errs
	}
}

