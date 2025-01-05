package internal

type Service interface {
	CreateArtist(artist ArtistDto) error
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateArtist(artist ArtistDto) error {
	return s.repo.CreateArtist(artist)
}
