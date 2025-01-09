package dtos

import "time"

type Artist struct {
	Id     string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID"`
	UserId string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	Name   string `json:"name" example:"The Weeknd" doc:"User email"`
	Bio    string `json:"bio" example:"The Weeknd took over pop music & culture on his own terms ..." doc:"Artist's biographical information and background"`
}

type Album struct {
	Id        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Album ID"`
	Title     string    `json:"title" example:"After Hours" doc:"Album title"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z" doc:"Album creation date"`
}

type Song struct {
	Id    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Song ID"`
	Title string `json:"title" example:"Blinding Lights" doc:"Song title"`
}
