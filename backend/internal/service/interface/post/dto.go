package post

type PostOrder uint

const (
	POST_ORDER_DATE_DESC PostOrder = iota
	POST_ORDER_DATE_ASC
)

type PostCreationForm struct {
	Title         string
	Content       string
	AllowComments bool
}

