package voicebr

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	RecDir  string = "recs"
	BaseDir string = "data"
)

type Store struct {
	RootDir string
}

func ensurePresent(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func (s *Store) PutRec(r io.Reader, name string) error {
	path := s.RecsPath()
	if err := ensurePresent(path); err != nil {
		return fmt.Errorf("PutRec: unable to obtain directory: %v", err)
	}

	path = filepath.Join(path, name)
	dest, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("PutRec: unable to create destination: %v", err)
	}

	n, err := io.Copy(dest, r)
	if err != nil {
		return fmt.Errorf("PutRec: unable to copy data: %v", err)
	}

	log.Printf("PutRec: %d bytes written into %s", n, path)
	return nil
}

func (s *Store) RecsPath() string {
	return filepath.Join(s.RootDir, BaseDir, RecDir)
}
