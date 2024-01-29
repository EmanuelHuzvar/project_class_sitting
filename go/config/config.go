package config

import (
	"Visma/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.Use(handlers.CorsMiddleware())
	userGroup := router.Group("/")
	{
		userGroup.POST("/login", handlers.LoginHandler)   //login
		userGroup.POST("/logout", handlers.LogoutHandler) //logout
		userGroup.GET("/courses", handlers.GetAllCourses)
		userGroup.GET("/availablecourses", handlers.GetAvailableCoursesHandler)
		userGroup.GET("/mycourses", handlers.GetMyCourses)
		userGroup.POST("/course/:course/reserve/:seat", handlers.ReserveCourseHandler)
		userGroup.PUT("/course/:id/unreserved", handlers.UnReserveCourseHandler)

	}
	testGroup := router.Group("/test")
	{
		testGroup.GET("/log", handlers.GetLogHandler)
	}
	adminGroup := router.Group("/")
	{
		adminGroup.GET("/admincourses", handlers.GetCoursesForAdminHandler)
		adminGroup.DELETE("/course/:id/delete", handlers.DeleteCourseHandler)
		adminGroup.PUT("/addcourse", handlers.AddCourseHandler)
	}
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return
	}
}
