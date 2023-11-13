package filesystem

type Operator interface {
	Reader
	Writer
}
