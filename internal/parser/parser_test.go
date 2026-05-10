package parser

import (
	"bytes"
	"godis/internal/resp"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]byte
		wantErr bool
	}{
		{
			name:    "Valid simple command",
			input:   "*1\r\n$4\r\nPING\r\n",
			want:    [][]byte{[]byte("PING")},
			wantErr: false,
		},
		{
			name:    "Valid multi-argument command",
			input:   "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			want:    [][]byte{[]byte("SET"), []byte("key"), []byte("value")},
			wantErr: false,
		},
		{
			name:    "Empty array",
			input:   "*0\r\n",
			want:    [][]byte{},
			wantErr: false,
		},
		{
			name:    "Array with empty bulk string",
			input:   "*1\r\n$0\r\n\r\n",
			want:    [][]byte{[]byte("")},
			wantErr: false,
		},
		{
			name:    "Bulk string containing CRLF inside",
			input:   "*1\r\n$4\r\nA\r\nB\r\n",
			want:    [][]byte{[]byte("A\r\nB")},
			wantErr: false,
		},
		{
			name:    "Multiple bulk strings with varying lengths",
			input:   "*2\r\n$1\r\nA\r\n$5\r\nABCDE\r\n",
			want:    [][]byte{[]byte("A"), []byte("ABCDE")},
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid starting character (Simple String instead of Array)",
			input:   "+OK\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Non-numeric array size",
			input:   "*abc\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing bulk string prefix $",
			input:   "*1\r\n4\r\nPING\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Mismatched bulk string length (declared 3, got 4)",
			input:   "*1\r\n$3\r\nPING\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Mismatched bulk string length (declared 5, got 4)",
			input:   "*1\r\n$5\r\nPING\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Missing final CRLF for bulk string",
			input:   "*1\r\n$4\r\nPING",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Array size declared 2, but only 1 provided",
			input:   "*2\r\n$4\r\nPING\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Negative array size (invalid for request)",
			input:   "*-2\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Bulk string length too short for data",
			input:   "*1\r\n$2\r\nHELLO\r\n",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRequest([]byte(tt.input))
			_, isErr := got.(*resp.RespSimpleError)
			if isErr != tt.wantErr {
				t.Errorf("ParseArray(%q) isErr = %v, wantErr %v", tt.input, isErr, tt.wantErr)
				return
			}
			if !tt.wantErr {
				respArray, ok := got.(*resp.RespArray)
				if !ok {
					t.Errorf("ParseArray(%q) got type %T, want *resp.RespArray", tt.input, got)
					return
				}

				if len(respArray.Value) != len(tt.want) {
					t.Errorf("ParseArray(%q) length got = %d, want %d", tt.input, len(respArray.Value), len(tt.want))
					return
				}
				for i := range respArray.Value {
					bulk, ok := respArray.Value[i].(*resp.RespBulkString)
					if !ok {
						t.Errorf("ParseArray(%q) element %d is type %T, want *resp.RespBulkString", tt.input, i, respArray.Value[i])
						continue
					}
					if !bytes.Equal(bulk.Value, tt.want[i]) {
						t.Errorf("ParseArray(%q) at index %d: got = %q, want %q", tt.input, i, string(bulk.Value), string(tt.want[i]))
					}
				}
			}
		})
	}
}
