package web

import "github.com/gorilla/sessions"

var cookieStore = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {
	cookieStore.Options = &sessions.Options{
		MaxAge: 86400,
	}
}
