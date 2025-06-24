package pgerrors

import (
	"github.com/jackc/pgx/v5/pgconn"
	"errors"
	"net/http"
)

var errorHTTPStatusMap = map[string]int{
	"00000": http.StatusOK,
	"01000": http.StatusOK,
	"01P01": http.StatusOK,
	"02000": http.StatusNoContent,
	"02001": http.StatusNoContent,
	"03000": http.StatusBadRequest,
	"08000": http.StatusServiceUnavailable,
	"08003": http.StatusServiceUnavailable,
	"08006": http.StatusServiceUnavailable,
	"08001": http.StatusServiceUnavailable,
	"08004": http.StatusForbidden,
	"08007": http.StatusInternalServerError,
	"08P01": http.StatusBadRequest,
	"09000": http.StatusInternalServerError,
	"0A000": http.StatusNotImplemented,
	"0B000": http.StatusBadRequest,
	"0F000": http.StatusInternalServerError,
	"0F001": http.StatusBadRequest,
	"0L000": http.StatusForbidden,
	"0LP01": http.StatusForbidden,
	"0P000": http.StatusForbidden,
	"0Z000": http.StatusInternalServerError,
	"0Z002": http.StatusInternalServerError,
	"20000": http.StatusNotFound,
	"21000": http.StatusConflict,
	"22000": http.StatusBadRequest,
	"22012": http.StatusBadRequest,
	"22003": http.StatusBadRequest,
	"22P02": http.StatusBadRequest,
	"22001": http.StatusBadRequest,
	"23000": http.StatusConflict,
	"23502": http.StatusBadRequest,
	"23503": http.StatusConflict,
	"23505": http.StatusConflict,
	"23514": http.StatusBadRequest,
	"23P01": http.StatusConflict,
	"24000": http.StatusBadRequest,
	"25000": http.StatusConflict,
	"25P01": http.StatusConflict,
	"25P02": http.StatusConflict,
	"26000": http.StatusBadRequest,
	"27000": http.StatusForbidden,
	"28000": http.StatusUnauthorized,
	"28P01": http.StatusUnauthorized,
	"2B000": http.StatusConflict,
	"2D000": http.StatusBadRequest,
	"2F000": http.StatusBadRequest,
	"34000": http.StatusBadRequest,
	"3D000": http.StatusBadRequest,
	"3F000": http.StatusBadRequest,
	"40000": http.StatusConflict,
	"40001": http.StatusConflict,
	"40002": http.StatusConflict,
	"40P01": http.StatusConflict,
	"42000": http.StatusBadRequest,
	"42601": http.StatusBadRequest,
	"42501": http.StatusForbidden,
	"42804": http.StatusBadRequest,
	"42P01": http.StatusNotFound,
	"42703": http.StatusBadRequest,
	"44000": http.StatusForbidden,
	"53000": http.StatusServiceUnavailable,
	"53100": http.StatusInsufficientStorage,
	"53200": http.StatusServiceUnavailable,
	"53300": http.StatusServiceUnavailable,
	"53400": http.StatusServiceUnavailable,
	"54000": http.StatusInternalServerError,
	"55000": http.StatusConflict,
	"55P03": http.StatusConflict,
	"57000": http.StatusServiceUnavailable,
	"57014": http.StatusRequestTimeout,
	"57P01": http.StatusServiceUnavailable,
	"58000": http.StatusInternalServerError,
	"58P01": http.StatusInternalServerError,
	"F0000": http.StatusInternalServerError,
	"HV000": http.StatusInternalServerError,
	"P0000": http.StatusInternalServerError,
	"P0001": http.StatusBadRequest,
	"XX000": http.StatusInternalServerError,
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
