package email

import (
	"strings"
	"testing"
)

func TestDecodeMessage_InvalidPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]any
		wantErr string
	}{
		{name: "nil payload", wantErr: "empty payload"},
		{name: "unsupported JSON value", payload: map[string]any{"to": make(chan int)}, wantErr: "unsupported type"},
		{name: "wrong recipient type", payload: map[string]any{"to": 42}, wantErr: "cannot unmarshal number"},
		{name: "missing recipient", payload: map[string]any{"invoice_number": "INV-1"}, wantErr: "missing 'to'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeMessage(TypeInvoice, tt.payload)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("decodeMessage() error = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
