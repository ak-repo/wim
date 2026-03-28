package utils

import "database/sql"

func StringNil(s sql.NullString) *string {
	if s.Valid && s.String != "" {
		return &s.String
	}
	return nil
}
