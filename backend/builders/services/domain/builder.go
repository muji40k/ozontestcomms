package domain

import (
	"github.com/muji40k/ozontestcomms/builders/errors"
	"github.com/muji40k/ozontestcomms/internal/domain/logic"
	commrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	postrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	usrrepo "github.com/muji40k/ozontestcomms/internal/repository/interface/user"
)

type LogicBuilder struct {
	comment commrepo.Repository
	post    postrepo.Repository
	user    usrrepo.Repository
}

func NewLogicBuilder() *LogicBuilder {
	return &LogicBuilder{nil, nil, nil}
}

func (self *LogicBuilder) WithCommentRepository(repo commrepo.Repository) *LogicBuilder {
	self.comment = repo
	return self
}

func (self *LogicBuilder) WithPostRepository(repo postrepo.Repository) *LogicBuilder {
	self.post = repo
	return self
}

func (self *LogicBuilder) WithUserRepository(repo usrrepo.Repository) *LogicBuilder {
	self.user = repo
	return self
}

func (self *LogicBuilder) Build() (*logic.Logic, error) {
	if nil == self.comment || nil == self.post || nil == self.user {
		return nil, errors.NotReady("logic.Logic")
	}

	return logic.New(logic.Context{
		Comment: self.comment,
		Post:    self.post,
		User:    self.user,
	}), nil
}

