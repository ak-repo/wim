package httpx

import (
	"fmt"
	"net/http"

	"github.com/ak-repo/wim/internal/errs"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				var err error
				if e, ok := recovered.(*errs.Error); ok {
					err = e
				} else {
					err = errs.E("middleware/Recover", errs.Unanticipated, fmt.Errorf("%v", recovered))
				}
				WriteError(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
