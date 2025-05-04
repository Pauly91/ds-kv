package client

import (
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

func Test_Client_Get(t *testing.T) {
	p := t.TempDir()
	c, e := NewKVClient(p)
	require.NoError(t, e, "error while creating client with a valid file path")

	c.Set("key1","value1")

	
}
