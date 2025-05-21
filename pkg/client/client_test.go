package client

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientCreation(t *testing.T) {
	p := t.TempDir()
	_, e := NewKVClient(p)
	require.NoError(t, e, "error while creating client with a valid file path")
	_, e = NewKVClient("")
	require.Error(t, e, "no error while creating client with an invalid file path")
}

func Test_Client_Set(t *testing.T) {

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "valid set operation",
			key:     "testKey",
			value:   "testValue",
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     "",
			value:   "testValue",
			wantErr: true,
		},
		{
			name:    "empty value",
			key:     "testKey",
			value:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := t.TempDir()
			c, e := NewKVClient(p)
			require.NoError(t, e, "error while creating client with a valid file path")
			e = c.Set(tt.key, tt.value)
			if tt.wantErr {
				require.Error(t, e, "error while setting a key value pair")
			} else {
				require.NoError(t, e, "error while setting a key value pair")
			}
		})
	}
}

func TestConcurrentSetAndGet(t *testing.T) {
	dir := t.TempDir()
	kvClient, err := NewKVClient(dir)
	if err != nil {
		t.Fatalf("failed to create KV client: %v", err)
	}

	const goroutines = 100
	var wg sync.WaitGroup

	// Launch concurrent Set
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			val := "value" + strconv.Itoa(i)
			if err := kvClient.Set(key, val); err != nil {
				t.Errorf("Set failed: %v", err)
			}
		}(i)
	}

	// Launch concurrent Get (may not find all if Set didn't finish)
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			_, _ = kvClient.Get(key) // ignore not-found errors
		}(i)
	}

	wg.Wait()
}
