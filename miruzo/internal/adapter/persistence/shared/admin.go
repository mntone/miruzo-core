package shared

type DatabaseAdminOptions struct {
	Database string
	UserName string
	Password string
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
