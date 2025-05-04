package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type KV interface {
	Set(string, string) error
	Get(string) (string, error)
}

type kv struct {
	fPath string
}

func NewKVClient(fPath string) (KV, error) {
	id := uuid.New()
	if _, e := os.Stat(fPath); os.IsNotExist(e) {
		return nil, fmt.Errorf("error creating file '%s' with error '%w'", fPath, e)
	}
	return &kv{
		fPath: filepath.Join(fPath, id.String()),
	}, nil
}

func (c *kv) Get(k string) (string, error) {
	file, e := os.Open(c.fPath)
	if e != nil {
		return "", fmt.Errorf("error opening file '%s' with error '%w'", c.fPath, e)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning lines from file '%s' with error '%w'", c.fPath, e)
	}

	// Iterate backwards through the lines
	for i := len(lines) - 1; i >= 0; i-- {
		parts := strings.SplitN(lines[i], ",", 2)
		if len(parts) < 2 {
			continue // skip malformed lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == k {
			return value, nil
		}
	}

	return "", fmt.Errorf("key '%s' not found", k)
}

func (c *kv) Set(k string, v string) error {

	if k == "" {
		return fmt.Errorf("key cannot be empty")
	}

	file, e := os.OpenFile(c.fPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		return fmt.Errorf("error opening file '%s' with error '%w'", c.fPath, e)
	}
	defer file.Close()
	_, e = file.WriteString(strings.Join([]string{k, v}, ","))
	if e != nil {
		return fmt.Errorf("error writing to file '%s' with error '%w'", c.fPath, e)
	}
	return nil
}
