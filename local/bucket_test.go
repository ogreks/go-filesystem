package local

import (
	"testing"
)

func TestCreateFile(t *testing.T) {
	storage := Bucket{
		filepath: "./aa/test.txt",
	}
	storage.CreateFile()
}
