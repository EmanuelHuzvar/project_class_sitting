package handlers

import (
	"Visma/db"
	"Visma/helpers"
	"Visma/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	tokenMap = map[models.UserJson]string{}
)

func LoginHandler(context *gin.Context) {

	var user models.UserJson
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInTokenMap := tokenMap[user]

	if len(userInTokenMap) == 0 {

		if db.CheckCredentials(user.UserName, user.Password) {

			HeaderToken := helpers.GenerateToken()
			user.Role = db.GetUserRoleByUsername(user.UserName, user.Password)
			tokenMap[user] = HeaderToken
			context.JSON(http.StatusOK, gin.H{"token": HeaderToken, "role": user.Role, "username": user.UserName})

		} else {
			helpers.LogFailedLogin(user.UserName)
			context.JSON(http.StatusUnauthorized, gin.H{"error": "bad credentials"})

		}
	} else {

		context.JSON(http.StatusUnauthorized, 401)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)
	zapLogger := logger.Sugar()
	zapLogger.Infow("Login request", "method", context.Request.Method, "path", context.Request.URL.Path, "status", context.Writer.Status(), "username", user.UserName)

}

func LogoutHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")

	if helpers.ContainsValue(tokenMap, HeaderToken) {

		delete(tokenMap, helpers.FindKeyByValue(tokenMap, HeaderToken))

		context.JSON(http.StatusOK, 200)

	} else {
		context.JSON(http.StatusUnauthorized, 401)
	}

}

func GetAllCourses(context *gin.Context) {
	courses, err := db.GetAllCourses()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "no courses right now"})
		return
	}
	context.JSON(http.StatusOK, courses)
}
func GetMyCourses(context *gin.Context) {
	user, err := helpers.GetUserByToken(tokenMap, context.GetHeader("token"))
	if !err {
		context.JSON(http.StatusUnauthorized, gin.H{"err2": "no suitable courses"})
	}
	courses, err2 := db.GetMyCourses(user.UserName)
	if err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err2": "no suitable courses"})
		return
	}
	context.JSON(http.StatusOK, courses)

}

func GetLogHandler(c *gin.Context) {
	wd, err := os.Getwd()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting working directory")
		return
	}

	// Construct the file path
	filePath := filepath.Join(wd, "login.txt")

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file")
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}

	// Send the file contents as the response
	c.Data(http.StatusOK, "application/octet-stream", data)
}
func ReserveCourseHandler(context *gin.Context) {
	token := context.GetHeader("token")
	seat := context.Param("seat")
	course := context.Param("course")
	if !helpers.ContainsValue(tokenMap, token) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := helpers.GetUserByToken(tokenMap, token)
	if !err {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	num, errorr := strconv.Atoi(seat)
	if errorr != nil {
		fmt.Println("Error:", err)
		return
	}

	err2 := db.AddParticipantToCourse(course, user.UserName, num)

	if err2 != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "status created"})
}
func UnReserveCourseHandler(context *gin.Context) {
	token := context.GetHeader("token")
	course := context.Param("id")
	if !helpers.ContainsValue(tokenMap, token) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := helpers.GetUserByToken(tokenMap, token)
	if !err {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	err2 := db.RemoveParticipantFromCourse(course, user.UserName)
	if err2 != nil {
		context.JSON(http.StatusNotModified, err2.Error())
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "you are not participant no more"})
}
func DeleteCourseHandler(context *gin.Context) {
	courseID := context.Param("id")
	if err := helpers.ValidateAndAuthorizeAdmin(context, tokenMap); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err := db.DeleteCourseByID(courseID)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "removed successfully"})

}
func AddCourseHandler(context *gin.Context) {
	token := context.GetHeader("token")
	var course models.Course

	if err := helpers.ValidateAndAuthorizeAdmin(context, tokenMap); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, isUser := helpers.GetUserByToken(tokenMap, token)
	if !isUser {
		context.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	course.Lector = user.UserName

	if err := context.ShouldBindJSON(&course); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//add this in front end T00:00:00Z

	err := db.AddCourseToDb(course)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Course added successfully"})
}
func GetAvailableCoursesHandler(context *gin.Context) {
	token := context.GetHeader("token")
	user, err1 := helpers.GetUserByToken(tokenMap, token)
	if !err1 {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
		return
	}
	courses, err := db.GetAvailableCourses(user.UserName)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "no courses right now"})
		return
	}

	context.JSON(http.StatusOK, courses)
}
func GetCoursesForAdminHandler(context *gin.Context) {
	token := context.GetHeader("token")
	user, err := helpers.GetUserByToken(tokenMap, token)
	if !err {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
		return
	}
	if err := helpers.ValidateAndAuthorizeAdmin(context, tokenMap); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	courses, err1 := db.GetCoursesByLector(user.UserName)
	if err1 != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "no courses right now"})
		return
	}
	context.JSON(http.StatusOK, courses)
}
