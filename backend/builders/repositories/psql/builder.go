package psql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/muji40k/ozontestcomms/builders/errors"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/psql"
	"github.com/muji40k/ozontestcomms/misc/nullable"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RepositoryBuilder struct {
	host     *nullable.Nullable[string]
	port     *nullable.Nullable[string]
	dbname   *nullable.Nullable[string]
	user     *nullable.Nullable[string]
	password *nullable.Nullable[string]
}

func NewRepositoryBuilder() *RepositoryBuilder {
	return &RepositoryBuilder{
		host:     nullable.None[string](),
		port:     nullable.None[string](),
		dbname:   nullable.None[string](),
		user:     nullable.None[string](),
		password: nullable.None[string](),
	}
}

func (self *RepositoryBuilder) WithHost(value string) *RepositoryBuilder {
	self.host = nullable.Some(value)
	return self
}

func (self *RepositoryBuilder) WithPort(value string) *RepositoryBuilder {
	self.port = nullable.Some(value)
	return self
}

func (self *RepositoryBuilder) WithDbname(value string) *RepositoryBuilder {
	self.dbname = nullable.Some(value)
	return self
}

func (self *RepositoryBuilder) WithUser(value string) *RepositoryBuilder {
	self.user = nullable.Some(value)
	return self
}

func (self *RepositoryBuilder) WithPassword(value string) *RepositoryBuilder {
	self.password = nullable.Some(value)
	return self
}

func (self *RepositoryBuilder) getConnString() (string, error) {
	if nullable.IsNone(self.host) || nullable.IsNone(self.port) ||
		nullable.IsNone(self.dbname) || nullable.IsNone(self.user) ||
		nullable.IsNone(self.password) {
		return "", errors.NotReady("psql.Repository")
	} else {
		return fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v",
			nullable.Unwrap(self.user),
			nullable.Unwrap(self.password),
			nullable.Unwrap(self.host),
			nullable.Unwrap(self.port),
			nullable.Unwrap(self.dbname),
		), nil
	}

}

func (self *RepositoryBuilder) Build() (*psql.Repository, func(), error) {
	var db *sqlx.DB
	cstr, err := self.getConnString()

	if nil == err {
		db, err = sqlx.Connect("pgx", cstr)
	}

	if nil == err {
		return psql.NewRepository(db), func() { db.Close() }, nil
	} else {
		return nil, nil, err
	}
}

