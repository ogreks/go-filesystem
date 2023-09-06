package local

import (
	"github.com/noOvertimeGroup/go-filesystem/local"
	"testing"
)

func TestCreateFile(t *testing.T) {
	storage := local.Bucket{
		Filepath: "./aa/test.txt",
	}
	storage.CreateFile()
}
