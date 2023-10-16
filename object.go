package filesystem

import (
	"strings"

	"github.com/noOvertimeGroup/go-filesystem/internal/errs"
)

type Object struct {
	Bucket string
	Target string
}

func NewObject(s string) (o Object, err error) {
	index := strings.Index(s, "/")
	if index == -1 {
		return Object{}, errs.ErrRelativePath
	}

	o.Bucket = s[:index]
	o.Target = s[index+1:]

	if o.Bucket == "" || o.Target == "" {
		return Object{}, errs.ErrNotFoundBucket
	}

	return
}
