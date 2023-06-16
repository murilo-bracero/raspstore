package middleware

import (
	"net/http"
)

func Authentication(requiredPermissions ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requesterRoles := r.Context().Value(UserRolesKey).([]string)

			if requesterRoles == nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			for rp := range requiredPermissions {
				for p := range requesterRoles {
					if rp == p {
						h.ServeHTTP(w, r)
						return
					}
				}
			}
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	}
}
