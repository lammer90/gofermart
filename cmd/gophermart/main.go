package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lammer90/gofermart/internal/config"
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/repository/orderstorage"
	"github.com/lammer90/gofermart/internal/repository/userstorage"
	"github.com/lammer90/gofermart/internal/services/authservice"
	"github.com/lammer90/gofermart/internal/services/orderservice"
	"github.com/lammer90/gofermart/internal/web/handlers/authhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/orderhandler"
	"github.com/lammer90/gofermart/internal/web/middleware/authfilter"
	"net/http"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")

	db := InitDB("pgx", config.DataSource)
	defer db.Close()

	cookieStore := buildSession()

	authSrv := authservice.New(userstorage.New(db), config.PrivateKey)
	authMdl := authfilter.New(authSrv, cookieStore)
	authHdl := authhandler.New(authSrv)

	orderSrv := orderservice.New(orderstorage.New(db))
	orderHdl := orderhandler.New(orderSrv, cookieStore)

	http.ListenAndServe(config.ServAddress, shortenerRouter(authHdl, orderHdl, authMdl))
}

func shortenerRouter(
	authProvider authhandler.AuthenticationRestApiProvider,
	orderProvider orderhandler.OrderRestApiProvider,
	middlewares ...func(next http.Handler) http.Handler) chi.Router {

	router := chi.NewRouter()
	for _, f := range middlewares {
		router.Use(f)
	}
	router.Post("/api/user/register", authProvider.Register)
	router.Post("/api/user/login", authProvider.Login)
	router.Post("/api/user/orders", orderProvider.Save)
	router.Get("/api/user/orders", orderProvider.FindAll)
	return router
}

func InitDB(driverName, dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	return db
}

func buildSession() *sessions.CookieStore {
	key := []byte("abc123")
	return sessions.NewCookieStore(key)
}
