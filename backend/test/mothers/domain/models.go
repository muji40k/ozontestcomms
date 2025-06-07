package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muji40k/ozontestcomms/builders/domain"
	"github.com/muji40k/ozontestcomms/misc/nullable"
)

func UserRandom() *domain.UserBuilder {
	id := uuid.Must(uuid.NewRandom())

	return domain.NewUserBuilder().
		WithId(id).
		WithEmail(fmt.Sprintf("user%v@mail.ru", id)).
		WithPassword("correct")
}

func PostRandomId() *domain.PostBuilder {
	return domain.NewPostBuilder().
		WithId(uuid.Must(uuid.NewRandom()))
}

func PostDefault(
	userId uuid.UUID,
	comments *nullable.Nullable[bool],
	prefix *nullable.Nullable[string],
	create *nullable.Nullable[time.Time],
) *domain.PostBuilder {
	return PostRandomId().
		WithAuthorId(userId).
		WithCommentsAllowed(nullable.GetOr(comments, true)).
		WithTitle("Post title" + nullable.GetOr(prefix, "")).
		WithContent("Post content" + nullable.GetOr(prefix, "")).
		WithCreationDate(nullable.GetOrFunc(create, time.Now))
}

func CommentRandomId() *domain.CommentBuilder {
	return domain.NewCommentBuilder().
		WithId(uuid.Must(uuid.NewRandom()))
}

func CommentDefault(
	userId uuid.UUID,
	targetId uuid.UUID,
	prefix *nullable.Nullable[string],
	create *nullable.Nullable[time.Time],
) *domain.CommentBuilder {
	return CommentRandomId().
		WithAuthorId(userId).
		WithTargetId(targetId).
		WithContent("Post content" + nullable.GetOr(prefix, "")).
		WithCreationDate(nullable.GetOrFunc(create, time.Now))
}

