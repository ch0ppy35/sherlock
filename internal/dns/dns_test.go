package dns

import (
	"testing"
)

func TestCompareRecords(t *testing.T) {
	type args struct {
		expected []string
		actual   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All records match exactly",
			args: args{
				expected: []string{"example.com", "test.com"},
				actual:   []string{"example.com", "test.com"},
			},
			wantErr: false,
		},
		{
			name: "Some records are missing",
			args: args{
				expected: []string{"example.com", "test.com", "foo.com"},
				actual:   []string{"example.com", "test.com"},
			},
			wantErr: true,
		},
		{
			name: "Some unexpected records are present",
			args: args{
				expected: []string{"example.com", "test.com"},
				actual:   []string{"example.com", "test.com", "unexpected.com"},
			},
			wantErr: true,
		},
		{
			name: "Both missing and unexpected records are present",
			args: args{
				expected: []string{"example.com", "test.com", "foo.com"},
				actual:   []string{"example.com", "unexpected.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty input cases",
			args: args{
				expected: []string{},
				actual:   []string{},
			},
			wantErr: false,
		},
		{
			name: "Empty expected records",
			args: args{
				expected: []string{},
				actual:   []string{"example.com"},
			},
			wantErr: true,
		},
		{
			name: "Empty actual records",
			args: args{
				expected: []string{"example.com"},
				actual:   []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareRecords(tt.args.expected, tt.args.actual); (err != nil) != tt.wantErr {
				t.Errorf("CompareRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
