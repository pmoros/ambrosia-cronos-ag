package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	resty "github.com/go-resty/resty/v2"
	"github.com/mondracode/ambrosia-atlas-api/graph/generated"
	"github.com/mondracode/ambrosia-atlas-api/graph/model"
)

// EnrollSubject is the resolver for the EnrollSubject field.
func (r *mutationResolver) EnrollSubject(ctx context.Context, input model.Enrollment) ([]*model.CourseGroup, error) {
	// panic(fmt.Errorf("not implemented: EnrollSubject - EnrollSubject"))
	client := resty.New()
	courseGroups := []*model.CourseGroup{}

	client.R().
		SetBody(input).
		EnableTrace().
		// Post("https://athenea-api-nf2vn5roaq-uc.a.run.app:enrollments/students")
		Post("http://127.0.0.1:8080/enrollments/students")

	return courseGroups, nil
}

// CourseGroups is the resolver for the courseGroups field.
func (r *queryResolver) CourseGroups(ctx context.Context, courseName string, courseGroupCode string, courseCode string) ([]*model.CourseGroup, error) {
	client := resty.New()
	courseGroups := []*model.CourseGroup{}

	client.R().
		SetQueryParams(map[string]string{
			"courseName":      courseName,
			"courseGroupCode": courseGroupCode,
			"courseCode":      courseCode,
		}).
		SetResult(&courseGroups).
		EnableTrace().
		// Get("https://athenea-api-nf2vn5roaq-uc.a.run.app:/course-groups")
		Get("http://127.0.0.1:8080/course-groups")

	// panic(fmt.Errorf("not implemented: CourseGroups - courseGroups"))
	return courseGroups, nil
}

// Schedules is the resolver for the schedules field.
func (r *queryResolver) Schedules(ctx context.Context, professorUsername *string, studentUsername *string) ([]*model.Schedule, error) {
	// panic(fmt.Errorf("not implemented: Schedules - schedules"))
	client := resty.New()
	schedules := []*model.Schedule{}
	if professorUsername == nil {
		client.R().
			SetQueryParams(map[string]string{
				"studentUsername": *studentUsername,
			}).
			SetResult(&schedules).
			EnableTrace().
			// Get("https://athenea-api-nf2vn5roaq-uc.a.run.app:/schedules/students")
			Get("http://127.0.0.1:8080/schedules/students")
	} else {
		client.R().
			SetQueryParams(map[string]string{
				"professorUsername": *professorUsername,
			}).
			SetResult(&schedules).
			EnableTrace().
			// Get("https://athenea-api-nf2vn5roaq-uc.a.run.app:/schedules/professors")
			Get("http://127.0.0.1:8080/schedules/professors")
	}

	return schedules, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
