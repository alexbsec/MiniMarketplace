package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexbsec/MiniMarketplace/src/controllers"
)


type App struct {
    Router *http.ServeMux
}

func (app *App) initializeRoutes() {
    app.Router.HandleFunc("/products", controllers.HandleProducts)
    app.Router.HandleFunc("/products/", controllers.HandleProducts)

    app.Router.HandleFunc("/users", controllers.HandleUsers)
    app.Router.HandleFunc("/users/", controllers.HandleUsers)

    app.Router.HandleFunc("/login", controllers.HandleLogin)

    app.Router.HandleFunc("/wallets", controllers.HandleWallets)
    app.Router.HandleFunc("/wallets/", controllers.HandleWallets)
}

// Initialize the app with the router and services
func (app *App) initialize() {
    app.Router = http.NewServeMux()
    app.initializeRoutes()
}

func (app *App) Run(port string) {
    app.initialize()
    fmt.Printf("Servidor iniciado na porta %s...\n", port)
    log.Fatal(http.ListenAndServe(port, app.Router))
}
