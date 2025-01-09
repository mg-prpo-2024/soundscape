package internal

import "time"

type CreateArtistDto struct {
	UserId string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"User ID"`
	Name   string `json:"name" example:"The Weeknd" doc:"User email"`
	Bio    string `json:"bio" example:"The Weeknd took over pop music & culture on his own terms ..." doc:"Artist's biographical information and background"`
}

type ArtistDto struct {
	Id string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID"`
	CreateArtistDto
}

type CreateAlbumDto struct {
	Title    string `json:"title" example:"After Hours" doc:"Album title"`
	ArtistId string `json:"artist_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID" format:"uuid"`
}

type AlbumDto struct {
	Id        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Album ID"`
	Title     string    `json:"title" example:"After Hours" doc:"Album title"`
	CreatedAt time.Time `json:"created_at" example:"2021-01-01T00:00:00Z" doc:"Album creation date"`
}

type CreateSongDto struct {
	Title      string `json:"title" example:"Blinding Lights" doc:"Song title"`
	AlbumId    string `json:"album_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Album ID" format:"uuid"`
	ArtistId   string `json:"artist_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Album ID" format:"uuid"`
	TrackOrder uint   `json:"track_order" example:"1" doc:"Track order"`
}

type CreateSongResponseDto struct {
	Id       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Album ID"`
	UpoadUrl string `json:"upload_url" example:"https://example.com/upload" doc:"URL to upload the song"`
}

type SongDto struct {
	Id    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Song ID"`
	Title string `json:"title" example:"Blinding Lights" doc:"Song title"`
}
