package routes

import (
	"github.com/chikong/ordersystem/bootstrap"
)

func LoadWebRoutes(b *bootstrap.Bootstrapper) {

	//v1 := b.Party("/v1")
	//
	//// Register our controllers.
	//v1.Controller("/hello", new(controllers.HelloController))
	//
	//v1.Controller("/message", new(controllers.MessageController),)
	//
	//v1.Controller("/user", new(controllers.UserController),services.NewUserService())
	//
	//// Create our movie repository with some (memory) data from the datasource.
	//repo := repositories.NewMovieRepository(datasource.Movies)
	//// Create our movie service, we will bind it to the movie controller.
	//movieService := services.NewMovieService(repo)
	//
	//v1.Controller("/movies", new(controllers.MovieController),
	//	// Bind the "movieService" to the MovieController's Service (interface) field.
	//	movieService,
	//	// Add the basic authentication(admin:password) middleware
	//	// for the /movies based requests.
	//	middleware.BasicAuth)
}