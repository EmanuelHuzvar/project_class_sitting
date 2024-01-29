package models

import "time"

var (
	CsharpURL = "http://147.232.158.55:1337/submit"
	JavaURL   = "http://10.11.65.65:8080/receive_json"
	PythonURL = "http://147.232.159.53:8090/submit"
)

type TeamJson struct {
	TeamName   string   `json:"teamname"`
	Members    []string `json:"members"`
	Emails     []string `json:"emails"`
	LanguageID int      `json:"languageID"`
	Ai         bool     `json:"ai"`
}

type UserJson struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TaskToSubmit struct {
	Language string `json:"code"`
	Task     string `json:"language"`
	Id       int    `json:"id"`
}

type Course struct {
	ID           string            `json:"ID"`
	Date         time.Time         `firestore:"date"`
	Description  string            `firestore:"description"`
	Lector       string            `firestore:"lector"`
	Participants []ParticipantInfo `firestore:"participants"`
	Title        string            `firestore:"title"`
}

type ParticipantInfo struct {
	Seat     int    `firestore:"seat"`
	Username string `firestore:"username"`
}
