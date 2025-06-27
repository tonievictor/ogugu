package pgerrors

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

var errorHTTPStatusMap = map[string]int{
	"00000": http.StatusOK,
	"01000": http.StatusOK,
	"01P01": http.StatusOK,

	"02000": http.StatusNoContent,
	"02001": http.StatusNoContent,

	"22000": http.StatusBadRequest,
	"22001": http.StatusBadRequest,
	"22003": http.StatusBadRequest,
	"22012": http.StatusBadRequest,
	"22P02": http.StatusBadRequest,
	"23000": http.StatusConflict,
	"23502": http.StatusBadRequest,
	"23514": http.StatusBadRequest,
	"24000": http.StatusBadRequest,
	"26000": http.StatusBadRequest,
	"2D000": http.StatusBadRequest,
	"2F000": http.StatusBadRequest,
	"34000": http.StatusBadRequest,
	"3D000": http.StatusBadRequest,
	"3F000": http.StatusBadRequest,
	"42000": http.StatusBadRequest,
	"42601": http.StatusBadRequest,
	"42804": http.StatusBadRequest,
	"42703": http.StatusBadRequest,
	"P0001": http.StatusBadRequest,
	"08P01": http.StatusBadRequest,
	"0B000": http.StatusBadRequest,
	"0F001": http.StatusBadRequest,

	"27000": http.StatusForbidden,
	"42501": http.StatusForbidden,
	"44000": http.StatusForbidden,
	"0L000": http.StatusForbidden,
	"0LP01": http.StatusForbidden,
	"0P000": http.StatusForbidden,
	"08004": http.StatusForbidden,

	"28000": http.StatusUnauthorized,
	"28P01": http.StatusUnauthorized,

	"53100": http.StatusInsufficientStorage,

	"23503": http.StatusConflict,
	"23505": http.StatusConflict,
	"23P01": http.StatusConflict,
	"25000": http.StatusConflict,
	"25P01": http.StatusConflict,
	"25P02": http.StatusConflict,
	"2B000": http.StatusConflict,
	"40000": http.StatusConflict,
	"40001": http.StatusConflict,
	"40002": http.StatusConflict,
	"40P01": http.StatusConflict,
	"55000": http.StatusConflict,
	"55P03": http.StatusConflict,
	"21000": http.StatusConflict,

	"57014": http.StatusRequestTimeout,

	"0A000": http.StatusNotImplemented,

	"42P01": http.StatusNotFound,
	"20000": http.StatusNotFound,

	"08000": http.StatusServiceUnavailable,
	"08001": http.StatusServiceUnavailable,
	"08003": http.StatusServiceUnavailable,
	"08006": http.StatusServiceUnavailable,
	"53000": http.StatusServiceUnavailable,
	"53200": http.StatusServiceUnavailable,
	"53300": http.StatusServiceUnavailable,
	"53400": http.StatusServiceUnavailable,
	"57000": http.StatusServiceUnavailable,
	"57P01": http.StatusServiceUnavailable,

	"03000": http.StatusBadRequest,
}

func getStatusFromCode(code string) int {
	status := errorHTTPStatusMap[code]
	if status == 0 {
		status = http.StatusInternalServerError
	}

	return status
}

func Details(err error) (int, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		status := getStatusFromCode(pgErr.Code)
		return status, pgErr.Detail
	}

	return vanillaErrors(err)
}

func vanillaErrors(err error) (int, string) {
	if err == nil {
		return http.StatusOK, "No error"
	}

	switch err.Error() {
	case "sql: no rows in result set":
		return http.StatusNotFound, "Resource does not exist"
	default:
		return http.StatusInternalServerError, "An unexpected error occurred"
	}
}
