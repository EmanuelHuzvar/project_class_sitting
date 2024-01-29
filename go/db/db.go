package db

import (
	"Visma/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
)

var (
	ctx       = context.Background()
	projectID = "codefights-5b44a"
	keyPath   = "db/authentification.json"
)

func CheckCredentials(username string, password string) bool {
	var goodCredentials bool
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		// Handle error
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	query := client.Collection("UserKASV").Where("username", "==", username).Where("password", "==", password)
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error
		}
		// Do something with the document data
		data := doc.Data()

		if data["username"] == username && data["password"] == password {
			goodCredentials = true

		} else {
			goodCredentials = false

		}

	}

	return goodCredentials

}
func GetIdInDB(path string) string {
	var highestID int

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return ""
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	query := client.Collection(path)
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return ""
		}

		// Parse the document ID as an integer
		id, err := strconv.Atoi(doc.Ref.ID)

		if err != nil {
			return ""
		}

		if id > highestID {
			highestID = id
		}
	}

	highestID++ // Increment the highest ID by one
	str := strconv.Itoa(highestID)
	return str
}

func GetUserRoleByUsername(username string, password string) string {
	// Initialize the Firestore client
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return ""
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// Replace "users" with the name of your Firestore collection that stores user data
	query := client.Collection("UserKASV").Where("username", "==", username).Where("password", "==", password)

	docs, err := query.Documents(context.Background()).GetAll()
	if err != nil {
		log.Printf("Failed to query user documents: %v", err)
		return ""
	}

	if len(docs) == 0 {
		log.Println("User not found")
		return ""
	}

	doc := docs[0]

	// Replace "role" with the field name that stores the user's role in your Firestore document
	role, err := doc.DataAt("role")
	if err != nil {
		log.Printf("Failed to get user role: %v", err)
		return ""
	}

	if role == nil {
		log.Println("User role not found")
		return ""
	}

	// Convert the role to a string
	roleStr, ok := role.(string)
	if !ok {
		log.Println("User role is not a string")
		return ""
	}

	return roleStr
}
func DocumentWithIdExists(ctx context.Context, client *firestore.Client, collectionName, documentID string) (bool, error) {
	docRef := client.Collection(collectionName).Doc(documentID)
	docSnap, err := docRef.Get(ctx)

	if err != nil {

		if statusOfDocument, ok := status.FromError(err); ok && statusOfDocument.Code() == codes.NotFound {

			return false, nil
		}

		return false, err
	}

	if docSnap.Exists() {

		return true, nil
	}

	return false, nil
}

func GetAllCourses() ([]models.Course, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Printf("Failed to close Firestore client: %v", err)
		}
	}(client)

	// Query all courses
	iter := client.Collection("Courses").Documents(ctx)
	defer iter.Stop()

	var courses []models.Course

	// Iterate over the documents and retrieve the course data
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating over courses: %v", err)
			continue
		}

		var course models.Course
		if err := doc.DataTo(&course); err != nil {
			log.Printf("Error parsing course data: %v", err)
			continue
		}

		// Set the ID field of the course struct
		course.ID = doc.Ref.ID

		if course.Date.After(time.Now()) {
			courses = append(courses, course)
		}
	}

	return courses, nil
}

func GetMyCourses(user string) ([]models.Course, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	iter := client.Collection("Courses").Documents(ctx)
	defer iter.Stop()

	var courses []models.Course

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating over courses: %v", err)
			continue
		}

		var course models.Course
		if err := doc.DataTo(&course); err != nil {
			log.Printf("Error parsing course data: %v", err)
			continue
		}
		course.ID = doc.Ref.ID

		// Check if "participants" array contains the desired username
		for _, participant := range course.Participants {
			if participant.Username == user {
				if course.Date.After(time.Now()) {
					courses = append(courses, course)
				}
				break
			}
		}
	}

	return courses, nil
}
func AddParticipantToCourse(courseID string, user string, seat int) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// Get the course document to update
	courseRef := client.Collection("Courses").Doc(courseID)

	// Retrieve the course document data
	courseDoc, err := courseRef.Get(ctx)
	if err != nil {
		log.Printf("Failed to retrieve course document: %v", err)
		return err
	}

	// Parse the course data into a Course struct
	var course models.Course
	if err := courseDoc.DataTo(&course); err != nil {
		log.Printf("Error parsing course data: %v", err)
		return err
	}

	// Check if the seat is already taken by another participant
	for i := range course.Participants {
		if course.Participants[i].Seat == seat {
			log.Printf("Seat %d is already taken", seat)
			return errors.New("seat is already taken")
		}
	}

	// Check if the user is already registered in any course
	for i := range course.Participants {
		if course.Participants[i].Username == user {
			log.Printf("User is already registered in a course")
			return errors.New("user is already registered in a course")
		}
	}

	// Add the new participant to the course
	newParticipant := models.ParticipantInfo{
		Seat:     seat,
		Username: user,
	}
	course.Participants = append(course.Participants, newParticipant)

	// Update the course document
	_, err = courseRef.Set(ctx, course)
	if err != nil {
		log.Printf("Failed to update course document: %v", err)
		return err
	}

	log.Printf("Participant added to the course successfully")
	return nil
}

func DeleteCourseByID(courseID string) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	_, err = client.Collection("Courses").Doc(courseID).Delete(ctx)
	if err != nil {
		log.Printf("Failed to delete course: %v", err)
		return err
	}

	log.Printf("Course deleted successfully")
	return nil
}
func RemoveParticipantFromCourse(courseID string, user string) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// Get the course document to update
	courseDoc, err := client.Collection("Courses").Doc(courseID).Get(ctx)
	if err != nil {
		log.Printf("Failed to retrieve course document: %v", err)
		return err
	}

	var course models.Course
	if err := courseDoc.DataTo(&course); err != nil {
		log.Printf("Error parsing course data: %v", err)
		return err
	}

	// Find the index of the participant to remove
	participantIndex := -1
	for i, participant := range course.Participants {
		if participant.Username == user {
			participantIndex = i
			break
		}
	}

	// Check if the participant is found
	if participantIndex == -1 {
		log.Printf("Participant not found in the course")
		return errors.New("participant not found in the course")
	}

	// Remove the participant from the course
	course.Participants = append(course.Participants[:participantIndex], course.Participants[participantIndex+1:]...)

	// Update the course document
	_, err = client.Collection("Courses").Doc(courseID).Set(ctx, course)
	if err != nil {
		log.Printf("Failed to update course document: %v", err)
		return err
	}

	log.Printf("Participant removed from the course successfully")
	return nil
}
func AddCourseToDb(course models.Course) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	// Prepare the document data to be added
	docData := map[string]interface{}{
		"date":        course.Date,
		"description": course.Description,
		"lector":      course.Lector,
		"title":       course.Title,
	}

	// Add participants data if available
	if len(course.Participants) > 0 {
		participants := make([]map[string]interface{}, len(course.Participants))
		for i, participant := range course.Participants {
			participants[i] = map[string]interface{}{
				"seat":     participant.Seat,
				"username": participant.Username,
			}
		}
		docData["participants"] = participants
	} else {
		docData["participants"] = []map[string]interface{}{}
	}

	// Add the course document to the "courses" collection
	_, err = client.Collection("Courses").Doc(GetIdInDB("Courses")).Set(ctx, docData)
	if err != nil {
		return err
	}

	// Close the Firestore client
	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}
func GetAvailableCourses(user string) ([]models.Course, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			// Handle error if necessary
		}
	}(client)

	iter := client.Collection("Courses").Documents(ctx)
	defer iter.Stop()

	var courses []models.Course

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating over courses: %v", err)
			continue
		}

		var course models.Course
		if err := doc.DataTo(&course); err != nil {
			log.Printf("Error parsing course data: %v", err)
			continue
		}
		course.ID = doc.Ref.ID

		// Check if "participants" array does not contain the desired username
		participantFound := false
		for _, participant := range course.Participants {
			if participant.Username == user {
				participantFound = true
				break
			}
		}

		if !participantFound {
			courses = append(courses, course)
		}
	}

	return courses, nil
}
func GetCoursesByLector(lectorName string) ([]models.Course, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			// Handle error if necessary
		}
	}(client)

	query := client.Collection("Courses").Where("lector", "==", lectorName)
	iter := query.Documents(ctx)
	defer iter.Stop()

	var courses []models.Course

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating over courses: %v", err)
			continue
		}

		var course models.Course
		if err := doc.DataTo(&course); err != nil {
			log.Printf("Error parsing course data: %v", err)
			continue
		}
		course.ID = doc.Ref.ID

		courses = append(courses, course)
	}

	return courses, nil
}
