package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/mondracode/ambrosia-atlas-api/graph/model"
)

// EnrollCourses is the resolver for the EnrollCourses field.
func (r *mutationResolver) EnrollCourses(ctx context.Context, input model.EnrollmentInput) ([]*model.Enrollment, error) {
	panic(fmt.Errorf("not implemented: EnrollCourses - EnrollCourses"))
}

// UploadGrades is the resolver for the UploadGrades field.
func (r *mutationResolver) UploadGrades(ctx context.Context, input []*model.GradeInput) ([]*model.Grade, error) {
	panic(fmt.Errorf("not implemented: UploadGrades - UploadGrades"))
}

// Courses is the resolver for the Courses field.
func (r *queryResolver) Courses(ctx context.Context, code *string, name *string, component *string) ([]*model.Course, error) {
	client := resty.New()
	courses := []*model.Course{}
	coursesEndpoint := fmt.Sprintf("https://ebedb84e-b0a7-4762-ba03-512fc1d81606.mock.pstmn.io/%s", "courses")
	client.R().
		SetQueryParams(map[string]string{
			"code":      *code,
			"name":      *name,
			"component": *component,
		}).
		SetResult(&courses).
		EnableTrace().
		Get(coursesEndpoint)

	for _, course := range courses {
		courseGroups := []*model.CourseGroup{}
		groupsEndpoint := fmt.Sprintf("https://7b055da9-65da-4801-a842-0426b341d991.mock.pstmn.io/%s/%s", "course-groups", course.Code)
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
