package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	resty "github.com/go-resty/resty/v2"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/mondracode/ambrosia-atlas-api/graph/model"
)

// EnrollCourses is the resolver for the EnrollCourses field.
func (r *mutationResolver) EnrollCourses(ctx context.Context, input model.EnrollmentInput) ([]*model.Enrollment, error) {
	urlGrades := "https://b04e88f7-b644-4e1a-8d02-09e58818146e.mock.pstmn.io"

	client := resty.New()

	enrollments := []*model.Enrollment{}

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

	return enrollments, nil
}

// UploadGrades is the resolver for the UploadGrades field.
func (r *mutationResolver) UploadGrades(ctx context.Context, input []*model.GradeInput) ([]*model.Grade, error) {
	panic(fmt.Errorf("not implemented: UploadGrades - UploadGrades"))
}

// Courses is the resolver for the Courses field.
func (r *queryResolver) Courses(ctx context.Context, code *string, name *string, component *string) ([]*model.Course, error) {
	var urlCoursesService = "https://ebedb84e-b0a7-4762-ba03-512fc1d81606.mock.pstmn.io"
	var urlEnrollmentsService = "https://athenea-api-4axjffbidq-uc.a.run.app"
	client := resty.New()

	// Get courses from courses service
	courses := []*model.Course{}
	coursesEndpoint := fmt.Sprintf("%s/%s", urlCoursesService, "courses")
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

// Schedules is the resolver for the Schedules field.
func (r *queryResolver) Schedules(ctx context.Context, username *string) ([]*model.Schedule, error) {
	panic(fmt.Errorf("not implemented: Schedules - Schedules"))
}

// AcademicHistories is the resolver for the AcademicHistories field.
func (r *queryResolver) AcademicHistories(ctx context.Context, username string, academicHistoryCode string) ([]*model.AcademicHistory, error) {
	panic(fmt.Errorf("not implemented: AcademicHistories - AcademicHistories"))
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
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
