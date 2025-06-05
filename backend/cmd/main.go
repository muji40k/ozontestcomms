package main

import (
	"fmt"
	"os"

	"github.com/muji40k/ozontestcomms/internal/application"
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

var repositoryConstructors = map[string]func() (RepositoryContext, Clearable, error){
	"in-memory": nil,
	"psql":      nil,
}
var serviceConstructors = map[string]func(*RepositoryContext) (ServiceContext, Clearable, error){
	"domain": nil,
}
var appConstructors = map[string]func(*ServiceContext) (application.Application, error){
	"graphql": nil,
}

func main() {
	// repo := inmemory.New(
	//     func(adduser func(models.User), _ func(inmemory.Comment), _ func(models.Post)) {
	//         adduser(models.User{
	//             Id:       uuid.MustParse("9c3d7dba-d1b2-42de-b708-158e32f11623"),
	//             Email:    "aboba@mail.com",
	//             Password: "asdf",
	//         })
	//     },
	// )
	// srv := logic.New(logic.Context{repo, repo, repo})
	// app := graphql.New("0.0.0.0", "8080", time.Millisecond, graphql.Context{
	//     srv, srv, srv,
	// })

	cleaner := NewCleaner()
	defer cleaner.Clear()

	var rcontext RepositoryContext
	var scontext ServiceContext
	var app application.Application
	var err error

	rtype := getenvOr(ENV_REPOSITORY_TYPE, "in-memory")
	stype := getenvOr(ENV_REPOSITORY_TYPE, "domain")
	atype := getenvOr(ENV_REPOSITORY_TYPE, "graphql")

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

