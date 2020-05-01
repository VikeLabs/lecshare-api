package utils

type Course []struct {
	CatalogCourseID       string      `json:"__catalogCourseId"`
	PassedCatalogQuery    bool        `json:"__passedCatalogQuery"`
	DateStart             string      `json:"dateStart"`
	Pid                   string      `json:"pid"`
	ID                    string      `json:"id"`
	Title                 string      `json:"title"`
	SubjectCode           SubjectCode `json:"subjectCode"`
	CatalogActivationDate string      `json:"catalogActivationDate"`
	Score                 int         `json:"_score"`
}

type SubjectCode struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	LinkedGroup string `json:"linkedGroup"`
}
