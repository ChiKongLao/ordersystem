package routes

import (
	"github.com/chikong/ordersystem/bootstrap"
	"github.com/kataras/iris/middleware/pprof"
)

// Configure registers the necessary routes to the app.
func Configure(b *bootstrap.Bootstrapper) {
	// Register golang pprof
	b.Any("/debug/pprof/{debug:path}", pprof.New())

	LoadAPIRoutes(b)
	LoadWebRoutes(b)
	LoadWebSocketRoutes(b)
}
