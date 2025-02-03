package aws

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getAWSProfile(t *testing.T) {
	tests := []struct {
		name       string
		envVar     string
		want       string
		wantEnvVar bool
	}{
		{
			name:       "AWS_PROFILE is set",
			envVar:     "my-profile",
			want:       "my-profile",
			wantEnvVar: true,
		},
		{
			name:       "AWS_PROFILE is not set",
			envVar:     "",
			want:       "default",
			wantEnvVar: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantEnvVar {
				t.Setenv("AWS_PROFILE", tt.envVar)
			}
			defer os.Unsetenv("AWS_PROFILE")

			got := getAWSProfile()
			assert.Equal(t, tt.want, got)
		})
	}
}
