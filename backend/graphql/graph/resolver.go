package graph

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	commsrv "github.com/muji40k/ozontestcomms/internal/service/interface/comment"
	postsrv "github.com/muji40k/ozontestcomms/internal/service/interface/post"
	usrsrv "github.com/muji40k/ozontestcomms/internal/service/interface/user"
)

type services struct {
	user    usrsrv.Service
	comment commsrv.Service
	post    postsrv.Service
}

type Resolver struct {
	services services
}

func NewResolver(
	user usrsrv.Service,
	comment commsrv.Service,
	post postsrv.Service,
) Resolver {
	return Resolver{services{user, comment, post}}
}

