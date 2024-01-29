package helpers

import (
	"Visma/models"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func GenerateToken() string {
	//Random (10-1 +1) is 0-9 9-0
	//Random (max-min+1)+min
	//35 - 122
	//cant use 39 40 41 44 45 46 58 59 91 - 96
	//65 - 90
	//49 - 57
	var token = ""
	var num = 0
	for i := 0; i < 40; i++ {
		num = rand.Intn(90-49+1) + 49

		if (num < 57 && num > 49) || (num < 90 && num > 65) {
			token = token + string(rune(num)) //rune!!!
		} else {
			i--
		}

	}
	return token
}

func GeneratePassword() string {
	var password = ""
	var num = 0
	for i := 0; i < 15; i++ {
		num = rand.Intn(90-49+1) + 49

		if (num < 57 && num > 49) || (num < 90 && num > 65) {
			password = password + string(rune(num))
		} else {
			i--
		}

	}
	return password
}
func ParseRegisterDataForEmail(data models.TeamJson) []string {

	var email []string
	email = data.Emails
	return email
}
func EncodeToBase64(code string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(code))
	return encoded
}
func DecodeFromBase64(code string) string {
	decoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return string(decoded)

}
func IpAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {

				continue
			}

			for _, addr := range addrs {

				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}

	return "localhost"
}
func SendEmail(email string, password string, name string) {
	//https://www.youtube.com/watch?v=pAPWBHxnFHM
	auth := smtp.PlainAuth("", "vismatestemail@gmail.com", "grobupcwoqoxixpy", "smtp.gmail.com")
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Login\r\n" +
		"\r\n" +
		"This is a your credentials : \n" +
		"name :  " + name + "\n" +
		"password : " + password)

	// Send the email.
	err := smtp.SendMail("smtp.gmail.com:587", auth, "vismatestemail@gmail.com", to, msg)
	if err != nil {
		panic(err)
	}

}

//	func ContainsValue(m map[string]string, value string) bool {
//		for _, v := range m {
//			if v == value {
//				return true
//			}
//		}
//		return false
//	}
func ContainsValue(m map[models.UserJson]string, value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}
func FindKeyByValue(m map[models.UserJson]string, value string) models.UserJson {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	return models.UserJson{}
}
func GetUserByToken(tokenMap map[models.UserJson]string, token string) (models.UserJson, bool) {
	for user, t := range tokenMap {
		if t == token {
			return user, true
		}
	}
	return models.UserJson{}, false
}
func ValidateAndAuthorizeAdmin(context *gin.Context, tokenMap map[models.UserJson]string) error {
	token := context.GetHeader("token")

	user, isInMap := GetUserByToken(tokenMap, token)
	if !isInMap {
		return errors.New("invalid token")
	}
	if len(user.Role) == 0 {
		return errors.New("no role assigned")
	}

	if user.Role != "admin" {
		return errors.New("user is unauthorised to do this")
	}

	return nil
}
func LogFailedLogin(username string) {
	file, err := os.OpenFile("login.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()

	logTime := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] Failed login attempt\nUsername: %s\n\n", logTime, username)

	_, err = file.WriteString(logMessage)
	if err != nil {
		fmt.Println("Failed to write to log file:", err)
		return
	}
}
func IsServerRunning() bool {
	// Make a HEAD request to the server's address
	resp, err := http.Head("http://localhost:8080")
	if err != nil {
		return false
	}

	// Check if the response status code is OK (2xx)
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
func RemoveLogFile() error {
	filePath := "login.txt"

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, nothing to remove
		return nil
	}

	// Remove the file
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove log file: %w", err)
	}

	return nil
}
