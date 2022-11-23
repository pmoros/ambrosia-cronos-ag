package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/mondracode/ambrosia-atlas-api/graph/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EnrollCourses is the resolver for the EnrollCourses field.
func (r *mutationResolver) EnrollCourses(ctx context.Context, input model.EnrollmentInput) (*model.Enrollment, error) {
	urlGrades := os.Getenv("GRADES_ARTEMIS_URL")
	urlEnrollmentsMQ := os.Getenv("ENROLLMENTS_MQ_URL")
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
		Post(gradesEndpoint)

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
func (r *mutationResolver) UploadGrades(ctx context.Context, input model.GradesInput) ([]*model.Grade, error) {
	var urlGradesService = os.Getenv("GRADES_ARTEMIS_URL")
	var urlPublishGrades = fmt.Sprintf("%s/%s", urlGradesService, "publish-grades")

	fmt.Println(urlPublishGrades)

	grades := []*model.Grade{}

	client := resty.New()
	client.R().
		SetBody(input).
		SetResult(&grades).
		EnableTrace().
		Post(urlPublishGrades)

	return grades, nil
}

// Courses is the resolver for the Courses field.
func (r *queryResolver) Courses(ctx context.Context, code *string, name *string, component *string) ([]*model.Course, error) {
	var urlCoursesService = os.Getenv("COURSES_APOLLO_URL")
	var urlEnrollmentsService = os.Getenv("ENROLLMENTS_ATHENEA_URL")
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
	var urlCoursesService = os.Getenv("COURSES_APOLLO_URL")
	var urlEnrollmentsService = os.Getenv("ENROLLMENTS_ATHENEA_URL")
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
		courseGroups := courses
		client.R().
			SetQueryParams(map[string]string{
				"code": course.CourseCode,
			}).
			SetResult(&courseGroups).
			EnableTrace().
			Get(coursesEndpoint)
		if len(courseGroups) > 0 {
			courses = courseGroups
		}
	}

	return courses, nil
}

// AcademicHistories is the resolver for the AcademicHistories field.
func (r *queryResolver) AcademicHistories(ctx context.Context, userCode string, academicHistoryCode *string) ([]*model.AcademicHistory, error) {
	var urlGradesService = os.Getenv("GRADES_ARTEMIS_URL")
	var academicHistoriesEndpoint = fmt.Sprintf("%s/%s", urlGradesService, "academic-histories")
	var queryParams = map[string]string{}

	academicHistories := []*model.AcademicHistory{}

	if academicHistoryCode != nil {
		queryParams = map[string]string{
			"userCode":            userCode,
			"academicHistoryCode": *academicHistoryCode,
		}
	} else {
		queryParams = map[string]string{
			"userCode": userCode,
		}
	}

	client := resty.New()
	client.R().
		SetQueryParams(queryParams).
		SetResult(&academicHistories).
		EnableTrace().
		Get(academicHistoriesEndpoint)

	return academicHistories, nil
}

// PendingCourses is the resolver for the PendingCourses field.
func (r *queryResolver) PendingCourses(ctx context.Context, userCode string, academicHistoryCode string) ([]*model.Course, error) {
	var urlCoursesService = os.Getenv("COURSES_APOLLO_URL")
	var allCoursesEndpoint = fmt.Sprintf("%s/%s", urlCoursesService, "subjects/get_all")

	courses := []*model.Course{}
	client := resty.New()

	client.R().
		SetResult(&courses).
		EnableTrace().
		Get(allCoursesEndpoint)

	return courses[:4], nil
}

// Appointments is the resolver for the Appointments field.
func (r *queryResolver) Appointments(ctx context.Context, userCode string) ([]*model.Appointment, error) {
	var urlGradesService = os.Getenv("GRADES_ARTEMIS_URL")
	var appointmentsEndpoint = fmt.Sprintf("%s/%s/%s", urlGradesService, "appointments", userCode)
	appointments := []*model.Appointment{}
	client := resty.New()
	client.R().
		SetResult(&appointments).
		EnableTrace().
		Get(appointmentsEndpoint)

	return appointments, nil
}

// AllCourses is the resolver for the AllCourses field.
func (r *queryResolver) AllCourses(ctx context.Context, service string) (*model.XMLResponse, error) {
	var externalSOAP = "https://campus-kid-interface.jhonatan.net/Service.svc"
	var stringResponse string
	soapBody := `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
	<soapenv:Header/>
	<soapenv:Body>
	   <tem:GetDataFromApi/>
	</soapenv:Body>
 </soapenv:Envelope>`
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/xml").
		SetBody(soapBody).
		EnableTrace().
		Post(externalSOAP)

	fmt.Print(stringResponse)
	xmlResponse := model.XMLResponse{
		Data: resp.String(),
	}

	return &xmlResponse, err
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
