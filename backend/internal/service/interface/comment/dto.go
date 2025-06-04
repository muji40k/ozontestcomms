package comment

type CommentOrder uint

const (
	COMMENT_ORDER_DATE_DESC CommentOrder = iota
	COMMENT_ORDER_DATE_ASC
)

type CommentForm struct {
	Content string
}

