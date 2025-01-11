package dtos

type CreatePlaylist struct {
	Name string `json:"name" example:"Chill vibes" doc:"Playlist name"`
}

type Playlist struct {
	Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID"`
	UserId     string `json:"user_id" example:"google-oauth2|106027649681852458278" doc:"User ID"`
	Name       string `json:"name" example:"Chill vibes" doc:"Playlist name"`
	TotalSongs int    `json:"total_songs" example:"5" doc:"Number of songs in the playlist"`
}

type Song struct {
	Id    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Song ID"`
	Title string `json:"title" example:"Alone Again" doc:"Song name"`
}
type PlaylistFull struct {
	Id     string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Artist ID"`
	UserId string  `json:"user_id" example:"google-oauth2|106027649681852458278" doc:"User ID"`
	Name   string  `json:"name" example:"Chill vibes" doc:"Playlist name"`
	Songs  []*Song `json:"songs" doc:"Songs in the playlist"`
}
