package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/peter6866/foodie/controllers"
	"github.com/peter6866/foodie/middlewares"
	"github.com/peter6866/foodie/repositories"
	"github.com/peter6866/foodie/services"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(client *mongo.Client) *gin.Engine {
	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(config))

	userRepo := repositories.NewUserRepository(client)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	router.GET("/google_login", controllers.GoogleLogin)
	router.POST("/api/auth/loginOrRegister", authController.LoginOrRegister)

	categoryRepo := repositories.NewCategoryRepository(client)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	menuItemRepo := repositories.NewMenuItemRepository(client)
	menuItemService := services.NewMenuItemService(menuItemRepo, categoryRepo, userRepo)
	menuItemController := controllers.NewMenuItemController(menuItemService)

	router.GET("/api/categories", categoryController.GetAllCategories)

	// Authenticated routes
	authenticatedRoutes := router.Group("/api")
	authenticatedRoutes.Use(middlewares.AuthMiddleware)
	{
		authenticatedRoutes.GET("/auth/getUser", authController.GetUser)
		authenticatedRoutes.POST("/auth/setChefAndPartner", authController.SetChefAndPartner)

		menuItemRoutes := authenticatedRoutes.Group("/menuItems")
		{
			menuItemRoutes.GET("/:id", menuItemController.GetMenuItem)

			menuItemRoutes.GET("", menuItemController.GetAllMenuItems)
			menuItemRoutes.POST("", middlewares.ChefMiddleware(), menuItemController.CreateMenuItem)
			// menuItemRoutes.PUT("/:id", middlewares.ChefMiddleware(), menuItemController.UpdateMenuItem)
			menuItemRoutes.DELETE("/:id", middlewares.ChefMiddleware(), menuItemController.DeleteMenuItem)
		}

		categoryRoutes := authenticatedRoutes.Group("/categories")
		{
			categoryRoutes.POST("", middlewares.ChefMiddleware(), categoryController.CreateCategory)
		}
	}

	return router
}
