/// Broadcast voice messages to a set of recipients.
/// Copyright (C) 2019 Daniel Morandini (jecoz)
///
/// This program is free software: you can redistribute it and/or modify
/// it under the terms of the GNU General Public License as published by
/// the Free Software Foundation, either version 3 of the License, or
/// (at your option) any later version.
///
/// This program is distributed in the hope that it will be useful,
/// but WITHOUT ANY WARRANTY; without even the implied warranty of
/// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
/// GNU General Public License for more details.
///
/// You should have received a copy of the GNU General Public License
/// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
