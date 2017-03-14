package web

import (
    "net/http"
	"fmt"
	
	"github.com/jackc/pgx"
	qv "github.com/jackc/quo_vadis"
	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/SCKelemen/Cassius/mail"
	"github.com/SCKelemen/Cassius/log"
	"github.com/SCKelemen/Cassius/common"
)

func NewAPIHandler(config common.AppConfig, pool *pgx.ConnPool, mailer mail.Mailer, logger log15.Logger) http.Handler {
	logger = logger.New("module", "api")
	log.SetFilterHandler("warn", logger, log15.StdoutHandler)


	router := qv.NewRouter()

	router.Post("/register", EnvHandler(pool, mailer, logger, RegisterHandler))
	router.Delete("/deregister", EnvHandler(pool, mailer, logger, AuthHandler(DeregisterHandler)))

	router.Post("/sessions", EnvHandler(pool, mailer, logger, CreateSessionHandler))
	router.Delete("/sessions", EnvHandler(pool, mailer, logger, AuthHandler(DeleteSessionHandler)))

	router.Get("/users/:uid", BaseHandler(pool, logger, UserHandler))
	router.Get("/users/:uid/image", BaseHandler(pool, logger, UserImageHandler))

	router.Post("/request_password_reset", EnvHandler(pool, mailer, logger, RequestPasswordResetHandler))
	router.Post("/reset_password", EnvHandler(pool, mailer, logger, ResetPasswordHandler))

	return router
}

type BaseHandlerFunc func(w http.ResponseWriter, req *http.Request, pool *pgx.ConnPool, logger log15.Logger)

func BaseHandler(pool *pgx.ConnPool, logger log15.Logger, f BaseHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f(w, req, pool, logger)
	})
}

type EnvHandlerFunc func(w http.ResponseWriter, req *http.Request, env *environment)

func EnvHandler(pool *pgx.ConnPool, mailer mail.Mailer, logger log15.Logger, f EnvHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user := GetUserFromSession(req, pool)
		env := &environment{user: user, pool: pool, mailer: mailer, logger: logger}
		f(w, req, env)
	})
}

func AuthHandler(f EnvHandlerFunc) EnvHandlerFunc {
	return EnvHandlerFunc(func(w http.ResponseWriter, req *http.Request, env *environment) {
		if env.user == nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Bad or missing X-Authentication header")
			return
		}
		f(w, req, env)
	})
}
