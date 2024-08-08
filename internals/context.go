package internals

import (
	"context"
	"net/http"
	"quiz-go/internals/quiz"
)

type AppContext struct {
	Store *quiz.SQLiteStore
}

func (ctx *AppContext) WithContext(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "appCtx", ctx))
}

func GetAppContext(r *http.Request) *AppContext {
	return r.Context().Value("appCtx").(*AppContext)
}

func StoreMiddleware(store *quiz.SQLiteStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appCtx := &AppContext{
				Store: store,
			}
			next.ServeHTTP(w, appCtx.WithContext(r))
		})
	}
}
