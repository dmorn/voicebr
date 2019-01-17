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
	"encoding/csv"
	"errors"
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

type Contact struct {
	Name   string `json:"-"`
	Type   string `json:"type"`
	Number string `json:"number"`
}

func NewContact(num, name string) *Contact {
	return &Contact{
		Type:   "phone",
		Number: num,
		Name:   name,
	}
}

func (s *Store) ContactsPath() string {
	return filepath.Join(s.RootDir, BaseDir, "contacts.csv")
}

var ErrCorruptedContacts = errors.New("contacts file read contains corrupted data, thus the results could be partial")

func (s *Store) Contacts() ([]*Contact, error) {
	file, err := os.Open(s.ContactsPath())
	if err != nil {
		return nil, fmt.Errorf("Contacts: unable to open file: %v", err)
	}
	defer file.Close()

	recs, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Contacts: unable to read file: %v", err)
	}

	acc := make([]*Contact, 0, len(recs))
	for _, rec := range recs {
		if len(rec) < 2 {
			// discard record
			continue
		}
		acc = append(acc, NewContact(rec[0], rec[1]))
	}
	if len(acc) != len(recs) {
		return acc, ErrCorruptedContacts
	}
	return acc, nil
}
