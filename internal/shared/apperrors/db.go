package apperrors

import (
	"errors"
	"net"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDBConnection = errors.New("database connection failed")
)

func IsDBConnectionError(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return strings.HasPrefix(pgErr.Code, "08") // 08 - ошибки коннекта к БД postgres
	}

	var netErr net.Error
	return errors.As(err, &netErr)
}
