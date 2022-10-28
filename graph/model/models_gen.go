// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AcademicHistory struct {
	Username            string        `json:"username"`
	AcademicHistoryCode string        `json:"academicHistoryCode"`
	Pa                  float64       `json:"pa"`
	Papa                float64       `json:"papa"`
	ProgramInfo         *ProgramInfo  `json:"programInfo"`
	ProgressInfo        *ProgressInfo `json:"progressInfo"`
	CreditsInfo         *CreditsInfo  `json:"creditsInfo"`
	Semesters           []*Semester   `json:"semesters"`
}

type Course struct {
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	Component    string         `json:"component"`
	Requirements []string       `json:"requirements"`
	Groups       []*CourseGroup `json:"groups"`
}

type CourseGroup struct {
	Code      string      `json:"code"`
	Capacity  float64     `json:"capacity"`
	Taken     float64     `json:"taken"`
	Professor *User       `json:"professor"`
	Schedules []*Schedule `json:"schedules"`
}

type CreditsInfo struct {
	Total          float64 `json:"total"`
	Disciplinar    float64 `json:"disciplinar"`
	Fundamentacion float64 `json:"fundamentacion"`
	LibreEleccion  float64 `json:"libreEleccion"`
	Nivelacion     float64 `json:"nivelacion"`
}

type Enrollment struct {
	StudentCode         string   `json:"studentCode"`
	AcademicHistoryCode string   `json:"academicHistoryCode"`
	CourseGroups        []string `json:"courseGroups"`
}

type EnrollmentInput struct {
	StudentCode         string   `json:"studentCode"`
	AcademicHistoryCode string   `json:"academicHistoryCode"`
	CourseGroups        []string `json:"courseGroups"`
}

type FinishedCourse struct {
	Code    string  `json:"code"`
	Credits float64 `json:"credits"`
	Grade   float64 `json:"grade"`
	Name    string  `json:"name"`
	Passed  float64 `json:"passed"`
}

type Grade struct {
	StudentUsername     string  `json:"studentUsername"`
	AcademicHistoryCode string  `json:"academicHistoryCode"`
	Grade               float64 `json:"grade"`
}

type GradeInput struct {
	StudentUsername     string  `json:"studentUsername"`
	AcademicHistoryCode string  `json:"academicHistoryCode"`
	Grade               float64 `json:"grade"`
}

type ProgramInfo struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProgressInfo struct {
	Total          string `json:"total"`
	Disciplinar    string `json:"disciplinar"`
	Fundamentacion string `json:"fundamentacion"`
	LibreEleccion  string `json:"libreEleccion"`
	Nivelacion     string `json:"nivelacion"`
}

type Schedule struct {
	CourseName    string `json:"courseName"`
	GroupCode     string `json:"groupCode"`
	ProfessorName string `json:"professorName"`
	Day           string `json:"day"`
	Building      string `json:"building"`
	Classroom     string `json:"classroom"`
	TimeOfStart   string `json:"timeOfStart"`
	TimeOfEnd     string `json:"timeOfEnd"`
}

type Semester struct {
	SemesterLabel string            `json:"semesterLabel"`
	Courses       []*FinishedCourse `json:"courses"`
}

type User struct {
	Code     string `json:"code"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
