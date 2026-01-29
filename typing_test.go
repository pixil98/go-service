package service

import (
	"testing"
)

func TestTypeOf(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    string
		wantErr bool
	}{
		"valid type field": {
			input:   `{"type": "worker"}`,
			want:    "worker",
			wantErr: false,
		},
		"type with extra fields": {
			input:   `{"type": "server", "port": 8080}`,
			want:    "server",
			wantErr: false,
		},
		"missing type field": {
			input:   `{"name": "test"}`,
			want:    "",
			wantErr: true,
		},
		"empty type field": {
			input:   `{"type": ""}`,
			want:    "",
			wantErr: true,
		},
		"invalid json": {
			input:   `not json`,
			want:    "",
			wantErr: true,
		},
		"empty json object": {
			input:   `{}`,
			want:    "",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := TypeOf([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("TypeOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TypeOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
