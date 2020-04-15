package callrelay

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadPid(r io.Reader) (int, error) {
	s, err := bufio.NewReader(r).ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("read pid: %w", err)
	}
	pid, err := strconv.Atoi(strings.TrimRight(s, "\n"))
	if err != nil {
		return 0, fmt.Errorf("read pid: %w", err)
	}
	return pid, nil
}

func LoadPid(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("load pid: %w", err)
	}
	defer f.Close()
	return ReadPid(f)
}

func WritePid(w io.Writer, pid int) error {
	_, err := w.Write([]byte(fmt.Sprintf("%d\n", pid)))
	if err != nil {
		return fmt.Errorf("write pid: %w", err)
	}
	return nil
}

func SavePid(filename string, pid int) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return fmt.Errorf("save pid: %w", err)
	}
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save pid: %w", err)
	}
	defer f.Close()
	return WritePid(f, pid)
}
