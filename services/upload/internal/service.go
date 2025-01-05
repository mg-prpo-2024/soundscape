package internal

type Service interface {
	CreateArtist(user ArtistDto) error
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateArtist(user ArtistDto) error {
	return s.repo.CreateArtist(user)
}
