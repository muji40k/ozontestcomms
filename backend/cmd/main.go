package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/builders/applications/graphql"
	psqlbuilder "github.com/muji40k/ozontestcomms/builders/repositories/psql"
	"github.com/muji40k/ozontestcomms/builders/services/domain"
	"github.com/muji40k/ozontestcomms/internal/application"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/inmemory"
	"github.com/muji40k/ozontestcomms/internal/repository/implementations/psql"
	commrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	postrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	usrrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/user"
	commsrv "github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	postsrv "github.com/muji40k/ozontestcomms/internal/service/interface/post"
	usrsrv "github.com/muji40k/ozontestcomms/internal/service/interface/user"
)

func getenvOr(key, def string) string {
	if v := os.Getenv(key); "" == v {
		return def
	} else {
		return v
	}
}

const (
	ENV_REPOSITORY_TYPE  string = "POSTER_REPOSITORY_TYPE"
	ENV_SERVICE_TYPE     string = "POSTER_SERVICE_TYPE"
	ENV_APPLICATION_TYPE string = "POSTER_APPLICATION_TYPE"
)

type RepositoryContext struct {
	Comment commrepo.Repository
	Post    postrepo.Repository
	User    usrrepo.Repository
}

type ServiceContext struct {
	Comment commsrv.Service
	Post    postsrv.Service
	User    usrsrv.Service
}

type Clearable interface {
	Clear()
}

type FCleaner func()

func (self FCleaner) Clear() {
	self()
}

type Cleaner []Clearable

func NewCleaner() Cleaner {
	return Cleaner(make([]Clearable, 0))
}

func (self *Cleaner) Push(v Clearable) {
	*self = append(*self, v)
}

func (self *Cleaner) Clear() {
	for i := len(*self) - 1; 0 <= i; i-- {
		(*self)[i].Clear()
	}
}

func InMemoryRepositoryConstructor() (RepositoryContext, Clearable, error) {
	repo := inmemory.New(
		func(adduser func(models.User), _ func(inmemory.Comment), _ func(models.Post)) {
			adduser(models.User{
				Id:       uuid.MustParse("9c3d7dba-d1b2-42de-b708-158e32f11623"),
				Email:    "aboba@mail.com",
				Password: "asdf",
			})
		},
	)

	return RepositoryContext{repo, repo, repo}, nil, nil
}

type PSQLRepositoryConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

const (
	ENV_PSQL_REPO_HOST     string = "POSTER_PSQL_HOST"
	ENV_PSQL_REPO_PORT     string = "POSTER_PSQL_PORT"
	ENV_PSQL_REPO_DB       string = "POSTER_PSQL_DBNAME"
	ENV_PSQL_REPO_USER     string = "POSTER_PSQL_USER"
	ENV_PSQL_REPO_PASSWORD string = "POSTER_PSQL_PASSWORD"
)

func PSQLRepositoryConfigEnvParser() (PSQLRepositoryConfig, error) {
	return PSQLRepositoryConfig{
		Host:     getenvOr(ENV_PSQL_REPO_HOST, "127.0.0.1"),
		Port:     getenvOr(ENV_PSQL_REPO_PORT, "5432"),
		DBName:   getenvOr(ENV_PSQL_REPO_DB, "poster"),
		User:     getenvOr(ENV_PSQL_REPO_USER, "postgres"),
		Password: getenvOr(ENV_PSQL_REPO_PASSWORD, "postgres"),
	}, nil
}

func PSQLRepositoryConstructor(parser func() (PSQLRepositoryConfig, error)) func() (RepositoryContext, Clearable, error) {
	return func() (RepositoryContext, Clearable, error) {
		var repo *psql.Repository
		var clr func()
		cfg, err := parser()

		if nil == err {
			repo, clr, err = psqlbuilder.NewRepositoryBuilder().
				WithHost(cfg.Host).
				WithPort(cfg.Port).
				WithDbname(cfg.DBName).
				WithUser(cfg.User).
				WithPassword(cfg.Password).
				Build()
		}

		if nil == err {
			return RepositoryContext{repo, repo, repo}, FCleaner(clr), nil
		} else {
			return RepositoryContext{}, nil, err
		}
	}
}

func DomainServiceConstructor(rcontext *RepositoryContext) (ServiceContext, Clearable, error) {
	svc, err := domain.NewLogicBuilder().
		WithCommentRepository(rcontext.Comment).
		WithPostRepository(rcontext.Post).
		WithUserRepository(rcontext.User).
		Build()

	return ServiceContext{svc, svc, svc}, nil, err
}

type GraphqlAppConfig struct {
	Host           string
	Port           string
	LoaderDuration time.Duration
}

const (
	ENV_GRAPHQL_APP_HOST            string = "POSTER_GRAPHQL_HOST"
	ENV_GRAPHQL_APP_PORT            string = "POSTER_GRAPHQL_PORT"
	ENV_GRAPHQL_APP_LOADER_DURATION string = "POSTER_GRAPHQL_LOADER"
)

func GraphqlAppConfigEnvParser() (GraphqlAppConfig, error) {
	host := getenvOr(ENV_GRAPHQL_APP_HOST, "0.0.0.0")
	port := getenvOr(ENV_GRAPHQL_APP_PORT, "80")
	var duration time.Duration
	var err error
	sduration := os.Getenv(ENV_GRAPHQL_APP_LOADER_DURATION)

	if "" == sduration {
		duration = time.Millisecond
	} else {
		duration, err = time.ParseDuration(sduration)
	}

	if nil != err {
		return GraphqlAppConfig{}, err
	} else {
		return GraphqlAppConfig{host, port, duration}, nil
	}
}

func GraphqlAppConstructor(
	parser func() (GraphqlAppConfig, error),
) func(*ServiceContext) (application.Application, error) {
	return func(scontext *ServiceContext) (application.Application, error) {
		var app application.Application
		cfg, err := parser()

		if nil == err {
			app, err = graphql.NewServerBuilder().
				WithHost(cfg.Host).
				WithPort(cfg.Port).
				WithLoaderDuration(cfg.LoaderDuration).
				WithCommentService(scontext.Comment).
				WithPostService(scontext.Post).
				WithUserService(scontext.User).
				Build()
		}

		return app, err
	}
}

var repositoryConstructors = map[string]func() (RepositoryContext, Clearable, error){
	"in-memory": InMemoryRepositoryConstructor,
	"psql":      PSQLRepositoryConstructor(PSQLRepositoryConfigEnvParser),
}
var serviceConstructors = map[string]func(*RepositoryContext) (ServiceContext, Clearable, error){
	"domain": DomainServiceConstructor,
}
var appConstructors = map[string]func(*ServiceContext) (application.Application, error){
	"graphql": GraphqlAppConstructor(GraphqlAppConfigEnvParser),
}

func main() {
	cleaner := NewCleaner()
	defer cleaner.Clear()

	var rcontext RepositoryContext
	var scontext ServiceContext
	var app application.Application
	var err error

	rtype := getenvOr(ENV_REPOSITORY_TYPE, "psql")
	stype := getenvOr(ENV_SERVICE_TYPE, "domain")
	atype := getenvOr(ENV_APPLICATION_TYPE, "graphql")

	if rconstr, found := repositoryConstructors[rtype]; !found {
		err = fmt.Errorf("Unknown repository type: %v", rtype)
	} else {
		var clr Clearable
		rcontext, clr, err = rconstr()

		if nil != clr {
			cleaner.Push(clr)
		}
	}

	if nil == err {
		if sconstr, found := serviceConstructors[stype]; !found {
			err = fmt.Errorf("Unknown service type: %v", stype)
		} else {
			var clr Clearable
			scontext, clr, err = sconstr(&rcontext)

			if nil != clr {
				cleaner.Push(clr)
			}
		}
	}

	if nil == err {
		if aconstr, found := appConstructors[atype]; !found {
			err = fmt.Errorf("Unknown application type: %v", atype)
		} else {
			app, err = aconstr(&scontext)

			if nil != app {
				cleaner.Push(app)
			}
		}
	}

	if nil == err {
		app.Run()
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

