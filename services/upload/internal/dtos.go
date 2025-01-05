package internal

type ArtistDto struct {
	Id    string `json:"id" example:"google-oauth2|106527689641250451478" doc:"User ID"`
	Email string `json:"email" format:"email" example:"test@gmail.com" doc:"User email"`
}
