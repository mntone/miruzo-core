package shared

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestDatabaseAdminOptionsResolveCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		options      DatabaseAdminOptions
		baseUser     string
		basePassword string
		wantUser     string
		wantPassword string
	}{
		{
			name: "override_both",
			options: DatabaseAdminOptions{
				UserName: "admin",
				Password: "secret",
			},
			baseUser:     "app",
			basePassword: "app-secret",
			wantUser:     "admin",
			wantPassword: "secret",
		},
		{
			name: "override_user_only",
			options: DatabaseAdminOptions{
				UserName: "admin",
			},
			baseUser:     "app",
			basePassword: "app-secret",
			wantUser:     "admin",
			wantPassword: "app-secret",
		},
		{
			name: "override_password_only",
			options: DatabaseAdminOptions{
				Password: "secret",
			},
			baseUser:     "app",
			basePassword: "app-secret",
			wantUser:     "app",
			wantPassword: "secret",
		},
		{
			name:         "use_base_values",
			options:      DatabaseAdminOptions{},
			baseUser:     "app",
			basePassword: "app-secret",
			wantUser:     "app",
			wantPassword: "app-secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotPassword := tt.options.ResolveCredentials(
				tt.baseUser,
				tt.basePassword,
			)
			assert.Equal(t, "ResolveCredentials() user", gotUser, tt.wantUser)
			assert.Equal(
				t,
				"ResolveCredentials() password",
				gotPassword,
				tt.wantPassword,
			)
		})
	}
}
