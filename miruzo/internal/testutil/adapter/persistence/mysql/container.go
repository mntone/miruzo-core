package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mysqlImageName = "mysql:9-oracle"
	mysqlDatabase  = "miruzo_test"
	mysqlUsername  = "m"
	mysqlPassword  = "miruzo1234"
	mysqlCharset   = "utf8mb4"
	mysqlCollation = "utf8mb4_0900_bin"
)

func startMySQLContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	container, err := mysql.Run(
		ctx,
		mysqlImageName,
		testcontainers.WithEnv(map[string]string{
			"MYSQL_DATABASE":      mysqlDatabase,
			"MYSQL_USER":          mysqlUsername,
			"MYSQL_ROOT_PASSWORD": mysqlPassword,
			"MYSQL_PASSWORD":      mysqlPassword,
		}),
		testcontainers.WithCmd(
			fmt.Sprintf("--character-set-server=%s", mysqlCharset),
			fmt.Sprintf("--collation-server=%s", mysqlCollation),
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("ready for connections").WithOccurrence(1).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		if container != nil {
			_ = container.Terminate(ctx)
		}
		return nil, fmt.Errorf("run mysql container: %w", err)
	}

	return container, nil
}
