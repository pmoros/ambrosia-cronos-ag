package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/mondracode/ambrosia-atlas-api/graph/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EnrollCourses is the resolver for the EnrollCourses field.
func (r *mutationResolver) EnrollCourses(ctx context.Context, input model.EnrollmentInput) (*model.Enrollment, error) {
	urlGrades := "https://b04e88f7-b644-4e1a-8d02-09e58818146e.mock.pstmn.io"
	urlEnrollmentsMQ := "amqp://guest:guest@34.125.61.62:5672/"
	// urlEnrollmentsMQ := "amqp://guest:guest@localhost:5672/"

	client := resty.New()

	enrollment := model.Enrollment{
		StudentCode:         input.StudentCode,
		AcademicHistoryCode: input.AcademicHistoryCode,
		CourseGroups:        input.CourseGroups,
	}

	// validates if the student can enroll the course
	gradesEndpoint := fmt.Sprintf("%s/%s", urlGrades, "can-enroll")
	resp, err := client.R().
		SetBody(input).
		EnableTrace().
		Get(gradesEndpoint)

	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error: %s", resp.Status())
	}

	// Enroll courses using enrollments service and enrollments queue
	conn, err := amqp.Dial(urlEnrollmentsMQ)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"enrollments", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	textToSend := "{"
	textToSend += input.StudentCode
	textToSend += "#"
	textToSend += input.AcademicHistoryCode
	textToSend += "#"

	for _, group := range input.CourseGroups {
		textToSend += group
		if group != input.CourseGroups[len(input.CourseGroups)-1] {
			textToSend += ","
		}
	}
	textToSend += "}"

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	err = enc.Encode(textToSend)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	if err != nil {
		log.Fatal("encode error:", err)
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(network.Bytes()),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", input)

	return &enrollment, nil
}

// UploadGrades is the resolver for the UploadGrades field.
func (r *mutationResolver) UploadGrades(ctx context.Context, input []*model.GradeInput) ([]*model.Grade, error) {
	panic(fmt.Errorf("not implemented: UploadGrades - UploadGrades"))
}

// Courses is the resolver for the Courses field.
func (r *queryResolver) Courses(ctx context.Context, code *string, name *string, component *string) ([]*model.Course, error) {
	var urlCoursesService = "https://apollo-api-jo5b4asiwq-uc.a.run.app"
	var urlEnrollmentsService = "https://athenea-api-4axjffbidq-uc.a.run.app"
	// var urlEnrollmentsService = "http://127.0.0.1:8080"
	client := resty.New()

	// Get courses from courses service
	courses := []*model.Course{}
	coursesEndpoint := fmt.Sprintf("%s/%s", urlCoursesService, "subjects")
	client.R().
		SetQueryParams(map[string]string{
			"code":      *code,
			"name":      *name,
			"component": *component,
		}).
		SetResult(&courses).
		EnableTrace().
		Get(coursesEndpoint)

	// Get course groups from enrollments service
	for _, course := range courses {
		courseGroups := []*model.CourseGroup{}
		groupsEndpoint := fmt.Sprintf("%s/%s/%s", urlEnrollmentsService, "course-groups", course.Code)
		client.R().
			SetResult(&courseGroups).
			EnableTrace().
			Get(groupsEndpoint)

		course.Groups = courseGroups
	}

	return courses, nil
}

// UserCourses is the resolver for the UserCourses field.
func (r *queryResolver) UserCourses(ctx context.Context, userCode *string) ([]*model.UserCourse, error) {
	// var urlCoursesService = "https://ebedb84e-b0a7-4762-ba03-512fc1d81606.mock.pstmn.io"
	var urlCoursesService = "https://apollo-api-jo5b4asiwq-uc.a.run.app"
	var urlEnrollmentsService = "https://athenea-api-4axjffbidq-uc.a.run.app"
	client := resty.New()

	courses := []*model.UserCourse{}
	// Get assigned courses from enrollments service
	enrollmentsEndpoint := fmt.Sprintf("%s/%s/%s", urlEnrollmentsService, "users/course-groups", *userCode)

	client.R().
		SetResult(&courses).
		EnableTrace().
		Get(enrollmentsEndpoint)

	// Set groups name
	for _, course := range courses {
		coursesEndpoint := fmt.Sprintf("%s/%s", urlCoursesService, "subjects")
		client.R().
			SetQueryParams(map[string]string{
				"code": course.CourseCode,
			}).
			SetResult(&courses).
			EnableTrace().
			Get(coursesEndpoint)
	}

	return courses, nil
}

// AcademicHistories is the resolver for the AcademicHistories field.
func (r *queryResolver) AcademicHistories(ctx context.Context, userCode string, academicHistoryCode string) ([]*model.AcademicHistory, error) {
	var urlGradesService = "https://b04e88f7-b644-4e1a-8d02-09e58818146e.mock.pstmn.io"
	var academicHistoriesEndpoint = fmt.Sprintf("%s/%s", urlGradesService, "academic-histories")
	academicHistories := []*model.AcademicHistory{}

	client := resty.New()
	client.R().
		SetQueryParams(map[string]string{
			"userCode":            userCode,
			"academicHistoryCode": academicHistoryCode,
		}).
		SetResult(&academicHistories).
		EnableTrace().
		Get(academicHistoriesEndpoint)

	return academicHistories, nil
}

// PendingCourses is the resolver for the PendingCourses field.
func (r *queryResolver) PendingCourses(ctx context.Context, userCode string, academicHistoryCode string) ([]*model.Course, error) {
	panic(fmt.Errorf("not implemented: PendingCourses - PendingCourses"))
}

// Appointments is the resolver for the Appointments field.
func (r *queryResolver) Appointments(ctx context.Context, userCode string) ([]*model.Appointment, error) {
	panic(fmt.Errorf("not implemented: Appointments - Appointments"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Schedules(ctx context.Context, username *string) ([]*model.Schedule, error) {
	panic(fmt.Errorf("not implemented: Schedules - Schedules"))
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
