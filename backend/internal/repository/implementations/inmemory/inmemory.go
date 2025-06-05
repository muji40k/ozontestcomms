package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type Target struct {
	Comment uuid.NullUUID
	Post    uuid.NullUUID
}

type Repository struct {
	users    map[uuid.UUID]models.User
	comments map[uuid.UUID]models.Comment
	posts    map[uuid.UUID]models.Post
	targets  map[uuid.UUID]Target
	mutex    sync.Mutex
}

func postOrder(order post.PostOrder) func(*time.Time, *time.Time) int {
	switch order {
	case post.POST_ORDER_DATE_ASC:
		return func(a *time.Time, b *time.Time) int {
			return a.Compare(*b)
		}
	case post.POST_ORDER_DATE_DESC:
		return func(a *time.Time, b *time.Time) int {
			return b.Compare(*a)
		}
	default:
		panic("Unknown order")
	}
}

func commentOrder(order comment.CommentOrder) func(*time.Time, *time.Time) int {
	switch order {
	case comment.COMMENT_ORDER_DATE_ASC:
		return func(a *time.Time, b *time.Time) int {
			return a.Compare(*b)
		}
	case comment.COMMENT_ORDER_DATE_DESC:
		return func(a *time.Time, b *time.Time) int {
			return b.Compare(*a)
		}
	default:
		panic("Unknown order")
	}
}

func find[T any](m map[uuid.UUID]T, id uuid.UUID, what string) (T, error) {
	if v, found := m[id]; found {
		return v, nil
	} else {
		return v, repoerrors.NotFound(what)
	}
}

func findFreeUUID[T any](m map[uuid.UUID]T) (uuid.UUID, error) {
	id, err := uuid.NewRandom()

	for found := false; nil == err && !found; {
		if _, f := m[id]; !f {
			found = true
		} else {
			id, err = uuid.NewRandom()
		}
	}

	return id, err
}

type Comment struct {
	Value  models.Comment
	Target Target
}

func New(init func(func(models.User), func(Comment), func(models.Post))) *Repository {
	users := make(map[uuid.UUID]models.User)
	comments := make(map[uuid.UUID]models.Comment)
	posts := make(map[uuid.UUID]models.Post)
	targets := make(map[uuid.UUID]Target)

	if nil != init {
		init(
			func(v models.User) {
				users[v.Id] = v
			},
			func(v Comment) {
				if v.Target.Comment.Valid && v.Target.Post.Valid ||
					!v.Target.Comment.Valid && !v.Target.Post.Valid {
					return
				}

				id, err := findFreeUUID(targets)

				if nil != err {
					return
				}

				v.Value.TargetId = id
				comments[v.Value.Id] = v.Value
				targets[id] = v.Target
			},
			func(v models.Post) {
				posts[v.Id] = v
			},
		)
	}

	return &Repository{users, comments, posts, targets, sync.Mutex{}}
}

func (self *Repository) createComment(
	comment models.Comment,
	finder func(*uuid.UUID) error,
	selector func(*Target) uuid.NullUUID,
) (models.Comment, error) {
	var locked = false
	var targetId uuid.UUID
	var currentId uuid.UUID
	var id uuid.UUID

	defer func() {
		if locked {
			self.mutex.Unlock()
		}
	}()

	err := finder(&comment.TargetId)

	if nil == err {
		_, err = find(self.users, comment.AuthorId, "comment author")
	}

	if nil == err {
		locked = true
		self.mutex.Lock()
	}

	if nil == err {
		found := false

		for id, v := range self.targets {
			t := selector(&v)
			if t.Valid && t.UUID == comment.TargetId {
				targetId = id
				found = true
				break
			}
		}

		if !found {
			err = repoerrors.NotFound("target id")
		}
	}

	if nil == err {
		id, err = findFreeUUID(self.comments)
	}

	if nil == err {
		currentId, err = findFreeUUID(self.targets)
	}

	if nil == err {
		comment.Id = id
		comment.TargetId = targetId
		self.comments[id] = comment
		self.targets[currentId] = Target{
			Comment: uuid.NullUUID{
				UUID:  id,
				Valid: true,
			},
		}
	}

	return comment, err
}

func (self *Repository) CreatePostComment(
	ctx context.Context,
	comment models.Comment,
) (models.Comment, error) {
	return self.createComment(
		comment,
		func(id *uuid.UUID) error {
			_, err := find(self.posts, *id, "comments post")
			return err
		},
		func(t *Target) uuid.NullUUID {
			return t.Post
		},
	)
}

func (self *Repository) CreateCommentComment(
	ctx context.Context,
	comment models.Comment,
) (models.Comment, error) {
	return self.createComment(
		comment,
		func(id *uuid.UUID) error {
			_, err := find(self.comments, *id, "comments root comment")
			return err
		},
		func(t *Target) uuid.NullUUID {
			return t.Comment
		},
	)
}

func (self *Repository) GetCommentsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Comment]], error) {
	return newPeekCollection(&self.comments, ids), nil
}

func (self *Repository) GetCommentsByPostId(
	ctx context.Context,
	postId uuid.UUID,
	order comment.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	cmp := commentOrder(order)
	var targetId uuid.UUID
	found := false

	for id, v := range self.targets {
		if v.Post.Valid && v.Post.UUID == postId {
			targetId = id
			found = true
			break
		}
	}

	if !found {
		return nil, repoerrors.NotFound("comments post id")
	} else {
		return collection.Map(newCollection(
			&self.comments,
			func(v *models.Comment) bool {
				return v.TargetId == targetId
			},
			func(a, b models.Comment) int {
				return cmp(&a.CreationDate, &b.CreationDate)
			},
		), func(v *models.Comment) result.Result[models.Comment] {
			return result.Ok(*v)
		}), nil
	}
}

func (self *Repository) GetCommentsByCommentId(
	ctx context.Context,
	commentId uuid.UUID,
	order comment.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	cmp := commentOrder(order)
	var targetId uuid.UUID
	found := false

	for id, v := range self.targets {
		if v.Comment.Valid && v.Comment.UUID == commentId {
			targetId = id
			found = true
			break
		}
	}

	if !found {
		return nil, repoerrors.NotFound("comments root comment id")
	} else {
		return collection.Map(newCollection(
			&self.comments,
			func(v *models.Comment) bool {
				return v.TargetId == targetId
			},
			func(a, b models.Comment) int {
				return cmp(&a.CreationDate, &b.CreationDate)
			},
		), func(v *models.Comment) result.Result[models.Comment] {
			return result.Ok(*v)
		}), nil
	}
}

func (self *Repository) CreatePost(
	ctx context.Context,
	post models.Post,
) (models.Post, error) {
	var locked = false
	var targetID uuid.UUID
	var id uuid.UUID

	defer func() {
		if locked {
			self.mutex.Unlock()
		}
	}()

	_, err := find(self.users, post.AuthorId, "post author")

	if nil == err {
		locked = true
		self.mutex.Lock()
	}

	if nil == err {
		id, err = findFreeUUID(self.posts)
	}

	if nil == err {
		targetID, err = findFreeUUID(self.targets)
	}

	if nil == err {
		post.Id = id
		self.posts[id] = post
		self.targets[targetID] = Target{
			Post: uuid.NullUUID{
				UUID:  id,
				Valid: true,
			},
		}
	}

	return post, err
}

func (self *Repository) GetPosts(
	ctx context.Context,
	order post.PostOrder,
) (collection.Collection[result.Result[models.Post]], error) {
	cmp := postOrder(order)

	return collection.Map(
		newCollection(&self.posts, nil, func(a, b models.Post) int {
			return cmp(&a.CreationDate, &b.CreationDate)
		}),
		func(v *models.Post) result.Result[models.Post] {
			return result.Ok(*v)
		},
	), nil
}

func (self *Repository) GetPostsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Post]], error) {
	return newPeekCollection(&self.posts, ids), nil
}

func (self *Repository) UpdatePost(
	ctx context.Context,
	post models.Post,
) (models.Post, error) {
	var locked = false

	defer func() {
		if locked {
			self.mutex.Unlock()
		}
	}()

	_, err := find(self.posts, post.Id, "post")

	if nil == err {
		locked = true
		self.mutex.Lock()
	}

	if nil == err {
		self.posts[post.Id] = post
	}

	return post, err
}

func (self *Repository) GetUsersById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.User]], error) {
	return newPeekCollection(&self.users, ids), nil
}

