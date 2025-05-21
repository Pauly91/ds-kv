package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type KV interface {
	Set(string, string) error
	Get(string) (string, error)
}

const fName = "fdf9a89c-332d-4984-964b-94f6169be9db"

type kv struct {
	fPath string
	mu    sync.Mutex
}

func NewKVClient(fPath string) (KV, error) {
	if fPath == "" {
		pwd, e := os.Getwd()
		fPath = pwd
		if e != nil {
			return nil, fmt.Errorf("error creating file: '%w'", e)
		}
	}
	if _, e := os.Stat(fPath); os.IsNotExist(e) {
		return nil, fmt.Errorf("error creating file '%s' with error '%w'", fPath, e)
	}
	return &kv{
		fPath: filepath.Join(fPath, fName),
	}, nil
}

func (c *kv) Get(k string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	if e := scanner.Err(); e != nil {
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
	c.mu.Lock()
	defer c.mu.Unlock()

	tempFilePath := c.fPath + ".tmp"

	// Recover from leftover temp file if exists
	if _, e := os.Stat(tempFilePath); e == nil {
		e = os.Rename(tempFilePath, c.fPath)
		if e != nil {
			return fmt.Errorf("failed to recover from temp file: %w", e)
		}
	}

	// Read existing data
	existingData, e := os.ReadFile(c.fPath)
	if e != nil && !os.IsNotExist(e) {
		return fmt.Errorf("failed to read existing file '%s': %w", c.fPath, e)
	}

	// Create temp file
	tempFile, e := os.Create(tempFilePath)
	if e != nil {
		return fmt.Errorf("failed to create temp file '%s': %w", tempFilePath, e)
	}
	defer tempFile.Close()

	// Write existing data + new data
	if len(existingData) > 0 {
		if _, e := tempFile.Write(existingData); e != nil {
			return fmt.Errorf("failed to write existing data: %w", e)
		}
	}
	if _, e := tempFile.WriteString(fmt.Sprintf("%s,%s\n", k, v)); e != nil {
		return fmt.Errorf("failed to write new line: %w", e)
	}

	// Atomically replace file
	if e := os.Rename(tempFilePath, c.fPath); e != nil {
		return fmt.Errorf("failed to rename temp file: %w", e)
	}

	return nil
}
