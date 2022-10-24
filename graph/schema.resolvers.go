package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

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
	panic(fmt.Errorf("not implemented: Courses - Courses"))
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
