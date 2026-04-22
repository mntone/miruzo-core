package shared

type DatabaseAdminOptions struct {
	DatabaseName string
	UserName     string
	Password     string
}

func (o DatabaseAdminOptions) ResolveDatabaseName(
	configDatabaseName string,
	defaultDatabaseName string,
) string {
	if o.DatabaseName != "" {
		return o.DatabaseName
	}
	if configDatabaseName != "" {
		return configDatabaseName
	}
	return defaultDatabaseName
}

func (o DatabaseAdminOptions) ResolveCredentials(
	baseUserName string,
	basePassword string,
) (string, string) {
	userName := baseUserName
	if o.UserName != "" {
		userName = o.UserName
	}

	password := basePassword
	if o.Password != "" {
		password = o.Password
	}

	return userName, password
}
