package integratefxwithgin3

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func addGroups(r *gin.Engine) {
	api := r.Group("/api")
	{
		admin := api.Group("/admin")
		{
			admin.GET("/", adminFunction)
		}
		users := api.Group("/users")
		{
			users.GET("/", usersFunction)
		}
	}
}

func adminFunction(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"adminFunction": "adminFunction content"})
}
func usersFunction(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"usersFunction": "usersFunction content"})
}

func Server(lc fx.Lifecycle) *gin.Engine {

	router := gin.Default()
	addGroups(router) // define rules for router

	srv := &http.Server{Addr: ":8080", Handler: router} // define a web server

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr) // the web server starts listening on 8080
			if err != nil {
				fmt.Println("[My Demo] Failed to start HTTP Server at", srv.Addr)
				return err
			}
			go srv.Serve(ln) // process an incoming request in a go routine
			fmt.Println("[My Demo]Succeeded to start HTTP Server at", srv.Addr)
			return nil

		},
		OnStop: func(ctx context.Context) error {
			srv.Shutdown(ctx) // stop the web server
			fmt.Println("[My Demo] HTTP Server is stopped")
			return nil
		},
	})

	return router
}

// This is a simple demo to show how integrate fx with gin
//
// http://localhost:8080/api/admin  (GET)
// Result:  "adminFunction": "adminFunction content"
//
// http://localhost:8080/api/users (GET)
// Result:  "usersFunction": "usersFunction content"
func Main() {
	app := fx.New(
		fx.Provide(
			Server,
		),
		fx.Invoke(func(*gin.Engine) {}),
	)
	app.Run()
}
