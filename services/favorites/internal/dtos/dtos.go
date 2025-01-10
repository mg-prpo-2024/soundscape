package dtos

type Playlist struct {
	Id   string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID"`
	Name string `json:"name" example:"Playist name" doc:"Chill vibes"`
}
