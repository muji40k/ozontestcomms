package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/muji40k/ozontestcomms/internal/domain/models"
	"github.com/muji40k/ozontestcomms/internal/repository/collection"
	repoerrors "github.com/muji40k/ozontestcomms/internal/repository/errors"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/comment"
	"github.com/muji40k/ozontestcomms/internal/repository/interface/post"
	"github.com/muji40k/ozontestcomms/misc/nullable"
	"github.com/muji40k/ozontestcomms/misc/result"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func generateId(
	ctx context.Context,
	db sqlx.PreparerContext,
	where string,
) (uuid.UUID, error) {
	var id uuid.UUID
	stmt, err := db.PrepareContext(ctx, fmt.Sprintf(`
        select exists(
            select * from %v where id = $1
        )
    `, where))

	if nil == err {
		id, err = uuid.NewRandom()

		for found := false; nil == err && !found; {
			row := stmt.QueryRowContext(ctx, id)
			err = row.Err()

			if nil == err {
				err = row.Scan(&found)
			}

			if nil == err {
				found = !found

				if !found {
					id, err = uuid.NewRandom()
				}
			}
		}
	}

	return id, err
}

func checkExists(
	ctx context.Context,
	db sqlx.QueryerContext,
	what string,
	where string,
	id uuid.UUID,
) error {
	var found bool

	err := sqlx.GetContext(ctx, db, &found,
		fmt.Sprintf("select exists(select * from %v where id = $1)", where),
		id,
	)

	if nil == err && !found {
		err = repoerrors.NotFound(what)
	}

	return err
}

func get[T any](
	ctx context.Context,
	db sqlx.QueryerContext,
	where string,
	id uuid.UUID,
) (T, error) {
	var out T

	err := sqlx.GetContext(ctx, db, &out,
		fmt.Sprintf("select * from %v where id = $1", where),
		id,
	)

	return out, err
}

func getPost(
	ctx context.Context,
	db sqlx.QueryerContext,
	id uuid.UUID,
) (Post, error) {
	var out Post

	err := sqlx.GetContext(ctx, db, &out, `
        select filtered.*, commentables.comments_allowed
        from (
            select * from posts.posts where id = $1
        ) as filtered
        join commentables.commentables
            on filtered.commentable_id = commentables.id
    `, id)

	return out, err
}

func generateOrder(start int, ids []uuid.UUID) string {
	order := make([]string, len(ids))

	for i, id := range ids {
		order[i] = fmt.Sprintf("('%v'::uuid, %v)", id, i)
	}

	return strings.Join(order, ", ")
}

func createCommentable(
	ctx context.Context,
	db interface {
		sqlx.ExtContext
		sqlx.PreparerContext
	},
	comments sql.NullBool,
) (Commentable, error) {
	out := Commentable{
		CommentsAllowed: comments,
	}
	var err error

	out.Id, err = generateId(ctx, db, "commentables.commentables")

	if nil == err {
		_, err = sqlx.NamedExecContext(ctx, db, `
            insert into commentables.commentables (
                id, comments_allowed
            ) values (
                :id, :comments_allowed
            )
        `, out)
	}

	return out, err
}

func createComment(
	ctx context.Context,
	db sqlx.ExtContext,
	comment Comment,
) error {
	_, err := sqlx.NamedExecContext(ctx, db, `
        insert into comments.comments (
            id, author_id, commentable_id, target_id, content,
            creation_date
        ) values (
            :id, :author_id, :commentable_id, :target_id, :content,
            :creation_date
        )
    `, comment)

	return err
}

func (self *Repository) CreatePostComment(
	ctx context.Context,
	comment models.Comment,
) (models.Comment, error) {
	var post Post
	lcomment := unmapComment(&comment)
	tx, err := self.db.Beginx()

	if nil == err {
		post, err = getPost(ctx, tx, comment.TargetId)
	}

	if nil == err {
		lcomment.TargetId = post.CommentableId
		lcomment.Id, err = generateId(ctx, tx, "comments.comments")
	}

	if nil == err {
		var comm Commentable
		comm, err = createCommentable(ctx, tx, sql.NullBool{})
		lcomment.CommentableId = comm.Id
	}

	if nil == err {
		err = createComment(ctx, tx, lcomment)
	}

	if nil == err {
		err = tx.Commit()
	}

	if nil == err {
		comment.Id = lcomment.Id
		comment.TargetId = lcomment.TargetId

		return comment, nil
	} else {
		if nil != tx {
			tx.Rollback()
		}

		return comment, err
	}
}

func (self *Repository) CreateCommentComment(
	ctx context.Context,
	comment models.Comment,
) (models.Comment, error) {
	var root Comment
	lcomment := unmapComment(&comment)
	tx, err := self.db.Beginx()

	if nil == err {
		root, err = get[Comment](ctx, tx, "comments.comments", comment.TargetId)
	}

	if nil == err {
		lcomment.TargetId = root.CommentableId
		lcomment.Id, err = generateId(ctx, tx, "comments.comments")
	}

	if nil == err {
		var comm Commentable
		comm, err = createCommentable(ctx, tx, sql.NullBool{})
		lcomment.CommentableId = comm.Id
	}

	if nil == err {
		err = createComment(ctx, tx, lcomment)
	}

	if nil == err {
		err = tx.Commit()
	}

	if nil == err {
		comment.Id = lcomment.Id
		comment.TargetId = lcomment.TargetId

		return comment, nil
	} else {
		if nil != tx {
			tx.Rollback()
		}

		return comment, err
	}
}

func (self *Repository) GetCommentsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Comment]], error) {
	if 0 == len(ids) {
		return collection.EmptyCollection[result.Result[models.Comment]](), nil
	}

	return collection.Map(newPeekCollection[qComment](ids, func(ids []uuid.UUID) (*sqlx.Rows, error) {
		var rows *sqlx.Rows

		stmt, err := self.db.PreparexContext(ctx, `
            select comments.*, orderer.ord
            from comments.comments
            right outer join (values `+generateOrder(1, ids)+`) as orderer (id, ord)
                on comments.id = orderer.id
            order by orderer.ord
        `)

		if nil == err {
			rows, err = stmt.QueryxContext(ctx)
		}

		return rows, err
	}), result.OkMapper(mapQComment)), nil
}

func (self *Repository) getCommentsByTarget(
	ctx context.Context,
	targetId uuid.UUID,
	order comment.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	sort, rel := mapCommentOrder(order)

	return collection.Map(newCollection[Comment](
		func(
			after *nullable.Nullable[uuid.UUID],
			limit *nullable.Nullable[uint],
		) (*sqlx.Rows, error) {
			builder := strings.Builder{}
			args := make([]any, 1, 4)
			cnt := 2

			fmt.Fprint(&builder,
				"select * from comments.comments where comments.target_id = $1",
			)
			args[0] = targetId

			nullable.IfSome(after, func(id *uuid.UUID) {
				fmt.Fprintf(&builder, `
                        and comments.creation_date %v (
                            select creation_date
                            from comments.comments
                            where comments.id = $%v
                        ) and comments.id != $%v`, rel, cnt, cnt+1)
				cnt += 2
				args = append(args, *id, *id)
			})

			fmt.Fprintf(&builder,
				" order by comments.creation_date %v", sort,
			)

			nullable.IfSome(limit, func(sz *uint) {
				fmt.Fprintf(&builder, " limit $%v", cnt)
				args = append(args, *sz)
			})

			var rows *sqlx.Rows
			stmt, err := self.db.PreparexContext(ctx, builder.String())

			if nil == err {
				rows, err = stmt.QueryxContext(ctx, args...)
			}

			return rows, err
		},
		func(id uuid.UUID) error {
			var found bool

			err := self.db.GetContext(ctx, &found, `
                    select exists(
                        select *
                        from comments.comments
                        where comments.id = $1 and comments.target_id = $2
                    )
                `, id, targetId)

			if nil == err && !found {
				err = repoerrors.NotFound("comment")
			}

			return err
		},
	), result.OkMapper(mapComment)), nil
}

func (self *Repository) GetCommentsByPostId(
	ctx context.Context,
	postId uuid.UUID,
	order comment.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	post, err := getPost(ctx, self.db, postId)

	if nil != err {
		return nil, err
	} else {
		return self.getCommentsByTarget(ctx, post.CommentableId, order)
	}
}

func (self *Repository) GetCommentsByCommentId(
	ctx context.Context,
	commentId uuid.UUID,
	order comment.CommentOrder,
) (collection.Collection[result.Result[models.Comment]], error) {
	comment, err := get[Comment](ctx, self.db, "comments.comments", commentId)

	if nil != err {
		return nil, err
	} else {
		return self.getCommentsByTarget(ctx, comment.CommentableId, order)
	}
}

func (self *Repository) GetUsersById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.User]], error) {
	if 0 == len(ids) {
		return collection.EmptyCollection[result.Result[models.User]](), nil
	}

	return collection.Map(newPeekCollection[qUser](ids, func(ids []uuid.UUID) (*sqlx.Rows, error) {
		var rows *sqlx.Rows

		stmt, err := self.db.PreparexContext(ctx, `
            select users.*, orderer.ord
            from users.users
            right outer join (values `+generateOrder(1, ids)+`) as orderer (id, ord)
                on users.id = orderer.id
            order by orderer.ord
        `)

		if nil == err {
			rows, err = stmt.QueryxContext(ctx)
		}

		return rows, err
	}), result.OkMapper(mapQUser)), nil
}

func (self *Repository) CreatePost(
	ctx context.Context,
	post models.Post,
) (models.Post, error) {
	lpost := unmapPost(&post)
	tx, err := self.db.Beginx()

	if nil == err {
		lpost.Id, err = generateId(ctx, tx, "posts.posts")
	}

	if nil == err {
		var comm Commentable
		comm, err = createCommentable(ctx, tx, lpost.CommentsAllowed)
		lpost.CommentableId = comm.Id
	}

	if nil == err {
		_, err = tx.NamedExecContext(ctx, `
            insert into posts.posts (
                id, author_id, commentable_id, title, content,
                creation_date
            ) values (
                :id, :author_id, :commentable_id, :title, :content,
                :creation_date
            )
        `, lpost)
	}

	if nil == err {
		err = tx.Commit()
	}

	if nil == err {
		post.Id = lpost.Id
		return post, nil
	} else {
		if nil != tx {
			tx.Rollback()
		}

		return post, err
	}
}

func (self *Repository) GetPosts(
	ctx context.Context,
	order post.PostOrder,
) (collection.Collection[result.Result[models.Post]], error) {
	sort, rel := mapPostOrder(order)

	return collection.Map(newCollection[Post](
		func(
			after *nullable.Nullable[uuid.UUID],
			limit *nullable.Nullable[uint],
		) (*sqlx.Rows, error) {
			builder := strings.Builder{}
			args := make([]any, 0, 3)
			cnt := 1

			fmt.Fprint(&builder, `
                select posts.*, commentables.comments_allowed
                from posts.posts
                join commentables.commentables
                    on posts.commentable_id = commentables.id
            `)

			nullable.IfSome(after, func(id *uuid.UUID) {
				fmt.Fprintf(&builder, `
                        where posts.creation_date %v (
                            select creation_date
                            from posts.posts
                            where posts.id = $%v
                        ) and posts.id != $%v`, rel, cnt, cnt+1)
				cnt += 2
				args = append(args, *id, *id)
			})

			fmt.Fprintf(&builder,
				" order by posts.creation_date %v", sort,
			)

			nullable.IfSome(limit, func(sz *uint) {
				fmt.Fprintf(&builder, " limit $%v", cnt)
				args = append(args, *sz)
			})

			var rows *sqlx.Rows
			stmt, err := self.db.PreparexContext(ctx, builder.String())

			if nil == err {
				rows, err = stmt.QueryxContext(ctx, args...)
			}

			return rows, err
		},
		func(id uuid.UUID) error {
			return checkExists(ctx, self.db, "post", "posts.posts", id)
		},
	), result.OkMapper(mapPost)), nil
}

func (self *Repository) GetPostsById(
	ctx context.Context,
	ids ...uuid.UUID,
) (collection.Collection[result.Result[models.Post]], error) {
	if 0 == len(ids) {
		return collection.EmptyCollection[result.Result[models.Post]](), nil
	}

	return collection.Map(newPeekCollection[qPost](ids, func(ids []uuid.UUID) (*sqlx.Rows, error) {
		var rows *sqlx.Rows

		stmt, err := self.db.PreparexContext(ctx, `
            select filtered.*, commentables.comments_allowed
            from (
                select posts.*, orderer.ord
                from posts.posts
                right outer join (values `+generateOrder(1, ids)+`) as orderer (id, ord)
                    on posts.id = orderer.id
            ) as filtered
            left outer join commentables.commentables
                on filtered.commentable_id = commentables.id
            order by filtered.ord
        `)

		if nil == err {
			rows, err = stmt.QueryxContext(ctx)
		}

		return rows, err
	}), result.OkMapper(mapQPost)), nil
}

func (self *Repository) UpdatePost(
	ctx context.Context,
	post models.Post,
) (models.Post, error) {
	lpost := unmapPost(&post)
	var tx *sqlx.Tx
	cpost, err := getPost(ctx, self.db, lpost.Id)

	if nil == err {
		lpost.CommentableId = cpost.CommentableId
	}

	if nil == err {
		tx, err = self.db.Beginx()
	}

	if nil == err {
		_, err = tx.NamedExec(`
            update posts.posts
            set author_id = :author_id,
                commentable_id = :commentable_id,
                title = :title,
                content = :content,
                creation_date = :creation_date
            where id = :id
        `, lpost)
	}

	if nil == err {
		_, err = tx.NamedExec(`
            update commentables.commentables
            set comments_allowed = :comments_allowed
            where id = :commentable_id
        `, lpost)
	}

	if nil == err {
		err = tx.Commit()
	}

	if nil == err {
		return post, nil
	} else {
		if nil != tx {
			tx.Rollback()
		}

		return post, err
	}
}

