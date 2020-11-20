package repositories

import (
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"regexp"
)

const PQUniqueViolation = "23505"

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// parseError take an error and passes it through the respective database error parser
// or returns the error itself if there is no matching parser.
// If there is a match, the resulting error may be a custom one.
func parseError(e error) error {
	log.Error("database error: ", e)

	pqErr, ok := e.(*pq.Error)
	if !ok {
		return e
	}

	switch pqErr.Code {
	case PQUniqueViolation:
		column, value := extractColumnValue(pqErr.Detail)
		msg := fmt.Sprintf("[%s] already exists with this value (%s)", column, value)

		return &ConflictError{
			Message: msg,
		}
	default:
		return e
	}
}

// extractColumnValue takes a string in the form of a sql error detail
// and returns the contained column and value
func extractColumnValue(detail string) (string, string) {
	var columnFinder = regexp.MustCompile(`Key \((.+)\)=`)
	var valueFinder = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)

	column := extractStringSubmatch(columnFinder, detail)
	value := extractStringSubmatch(valueFinder, detail)

	return column, value
}

// extractString takes a regex and a string and returns
// the matched string or an empty string if no match exists
func extractStringSubmatch(regex *regexp.Regexp, str string) string {
	results := regex.FindStringSubmatch(str)
	if len(results) < 2 {
		return ""
	}
	return results[1]
}
