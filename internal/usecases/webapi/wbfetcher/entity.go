package wbfetcher

import "time"

type WbResponse struct {
	Cards  []Cards `json:"cards"`
	Cursor Cursor  `json:"cursor"`
}
type Photos struct {
	Big      string `json:"big"`
	C246X328 string `json:"c246x328"`
	C516X688 string `json:"c516x688"`
	Square   string `json:"square"`
	Tm       string `json:"tm"`
}
type Dimensions struct {
	Length  int  `json:"length"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`
	IsValid bool `json:"isValid"`
}
type Characteristics struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}
type Sizes struct {
	ChrtID   int      `json:"chrtID"`
	TechSize string   `json:"techSize"`
	WbSize   string   `json:"wbSize"`
	Skus     []string `json:"skus"`
}
type Tags struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type Cards struct {
	NmID            int               `json:"nmID"`
	ImtID           int               `json:"imtID"`
	NmUUID          string            `json:"nmUUID"`
	SubjectID       int               `json:"subjectID"`
	SubjectName     string            `json:"subjectName"`
	VendorCode      string            `json:"vendorCode"`
	Brand           string            `json:"brand"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Photos          []Photos          `json:"photos"`
	Video           string            `json:"video"`
	Dimensions      Dimensions        `json:"dimensions"`
	Characteristics []Characteristics `json:"characteristics"`
	Sizes           []Sizes           `json:"sizes"`
	Tags            []Tags            `json:"tags"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type Cursor struct {
	NmID      int       `json:"nmID,omitempty"`
	Total     int       `json:"total,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Filter struct {
	WithPhoto int `json:"withPhoto"`
}

type Sort struct {
	Ascending bool `json:"ascending"` // true - asc sort, false - desc sort
}

type Settings struct {
	Sort   Sort   `json:"sort"`
	Cursor Cursor `json:"cursor"`
	Filter Filter `json:"filter"`
}

type Setting struct {
	Setting Settings `json:"settings"`
}
