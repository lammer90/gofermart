package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/lammer90/gofermart/internal/config"
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/repository/userstorage"
	"github.com/lammer90/gofermart/internal/services/authservice"
	"github.com/lammer90/gofermart/internal/web/handlers/authhandler"
	"github.com/lammer90/gofermart/internal/web/middleware/authfilter"
	"net/http"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")

	db := InitDB("pgx", config.DataSource)
	defer db.Close()

	authSrv := authservice.New(userstorage.New(db), config.PrivateKey)
	authMdl := authfilter.New(authSrv)
	authHdl := authhandler.New(authSrv)

	http.ListenAndServe(config.ServAddress, shortenerRouter(authHdl, authMdl))
}

func shortenerRouter(
	authProvider authhandler.AuthenticationRestApiProvider,
	middlewares ...func(next http.Handler) http.Handler) chi.Router {

	router := chi.NewRouter()
	for _, f := range middlewares {
		router.Use(f)
	}
	router.Post("/api/user/register", authProvider.Register)
	router.Post("/api/user/login", authProvider.Login)
	return router
}

func InitDB(driverName, dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	return db
}
