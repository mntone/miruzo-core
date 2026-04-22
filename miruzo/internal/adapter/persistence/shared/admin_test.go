package shared

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/testutil/assert"
)

func TestDatabaseAdminOptionsResolveAdminDatabaseName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		options    DatabaseAdminOptions
		configName string
		fallback   string
		want       string
	}{
		{
			name: "options_override",
			options: DatabaseAdminOptions{
				DatabaseName: "opt_admin",
			},
			configName: "cfg_admin",
			fallback:   "fallback_admin",
			want:       "opt_admin",
		},
		{
			name:       "config_when_options_empty",
			options:    DatabaseAdminOptions{},
			configName: "cfg_admin",
			fallback:   "fallback_admin",
			want:       "cfg_admin",
		},
		{
			name:       "fallback_when_all_empty",
			options:    DatabaseAdminOptions{},
			configName: "",
			fallback:   "fallback_admin",
			want:       "fallback_admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.options.ResolveDatabaseName(tt.configName, tt.fallback)
			assert.Equal(
				t,
				"ResolveAdminDatabaseName()",
				got,
				tt.want,
			)
		})
	}
}

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
