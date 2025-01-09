package internal

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, config *Config) {
	service := NewService(
		NewRepository(db),
		NewStorage(config),
	)
	registerCreateArtist(api, service)
	registerGetArtist(api, service)
	registerGetArtistAlbums(api, service)
	registerCreateAlbum(api, service)
	registerGetAlbum(api, service)
	registerCreateSong(api, service)
	registerGetSongs(api, service)
	registerDeleteSong(api, service)
}

type CreateArtistInput struct {
	Body CreateArtistDto
}

type CreateArtistOutput struct{}

func registerCreateArtist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-artist",
		Method:      http.MethodPost,
		Path:        "/artists",
		Summary:     "Create an artist",
		Description: "Create a new artist for the user id provided in the auth token.",
		Tags:        []string{"Artists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateArtistInput) (*CreateArtistOutput, error) {
		err := service.CreateArtist(input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateArtistOutput{}, nil
	})
}

type GetArtistInput struct {
	UserId string `path:"userId" doc:"User ID" example:"google-oauth2|106527649641850458478"`
}

type GetArtistOutput struct {
	Body ArtistDto
}

func registerGetArtist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-artist",
		Method:      http.MethodGet,
		Path:        "/artists/{userId}",
		Summary:     "Get an artist",
		Description: "Get the artist metadata for the user.",
		Tags:        []string{"Artists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetArtistInput) (*GetArtistOutput, error) {
		artist, err := service.GetArtist(input.UserId)
		if err != nil {
			return nil, err
		}
		return &GetArtistOutput{
			Body: *artist,
		}, nil
	})
}

type GetArtistAlbumsInput struct {
	ArtistId string `path:"artistId" doc:"Artist ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type GetArtistAlbumsOutput struct {
	Body []*AlbumDto
}

func registerGetArtistAlbums(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-artist-albums",
		Method:      http.MethodGet,
		Path:        "/artists/{artistId}/albums",
		Summary:     "Get artist albums",
		Description: "Get the albums for the artist.",
		Tags:        []string{"Artists", "Albums"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetArtistAlbumsInput) (*GetArtistAlbumsOutput, error) {
		albums, err := service.GetArtistAlbums(input.ArtistId)
		if err != nil {
			return nil, err
		}
		return &GetArtistAlbumsOutput{
			Body: albums,
		}, nil
	})
}

type CreateAlbumInput struct {
	Body CreateAlbumDto
}

type CreateAlbumOutput struct {
	Body AlbumDto
}

func registerCreateAlbum(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-album",
		Method:      http.MethodPost,
		Path:        "/albums",
		Summary:     "Create an album",
		Description: "Create a new album.",
		Tags:        []string{"Albums"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateAlbumInput) (*CreateAlbumOutput, error) {
		album, err := service.CreateAlbum(input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateAlbumOutput{
			Body: *album,
		}, nil
	})
}

type GetAlbumInput struct {
	Id string `path:"id" doc:"Album ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type GetAlbumOutput struct {
	Body AlbumDto
}

func registerGetAlbum(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-album",
		Method:      http.MethodGet,
		Path:        "/albums/{id}",
		Summary:     "Get an album",
		Description: "Get the album metadata.",
		Tags:        []string{"Albums"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetAlbumInput) (*GetAlbumOutput, error) {
		album, err := service.GetAlbum(input.Id)
		if err != nil {
			return nil, err
		}
		return &GetAlbumOutput{
			Body: *album,
		}, nil
	})
}

type CreateSongInput struct {
	Body CreateSongDto
}
type CreateSongOutput struct {
	Body CreateSongResponseDto // or something else
}

func registerCreateSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-song",
		Method:      http.MethodPost,
		Path:        "/songs",
		Summary:     "Create a song",
		Description: "Create a new song entry and return a pre-signed Azure URL that allows for secure upload of the song audio file.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateSongInput) (*CreateSongOutput, error) {
		song, err := service.CreateSong(input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateSongOutput{
			Body: *song,
		}, nil
	})
}

type GetSongsInput struct {
	AlbumId string `path:"albumId" doc:"Album ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type GetSongsOutput struct {
	Body []*SongDto
}

func registerGetSongs(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-album-songs",
		Method:      http.MethodGet,
		Path:        "/albums/{albumId}/songs",
		Summary:     "Get album songs",
		Description: "Get album songs.",
		Tags:        []string{"Albums", "Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetSongsInput) (*GetSongsOutput, error) {
		songs, err := service.GetAlbumSongs(input.AlbumId)
		if err != nil {
			return nil, err
		}
		return &GetSongsOutput{
			Body: songs,
		}, nil
	})
}

type DeleteSongInput struct {
	SongId string `path:"songId" doc:"Song ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}

func registerDeleteSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "delete-song",
		Method:      http.MethodDelete,
		Path:        "/songs/{songId}",
		Summary:     "Delete a song",
		Description: "Delete a song.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *DeleteSongInput) (*struct{}, error) {
		err := service.DeleteSong(input.SongId)
		logrus.Error(err)
		return nil, err
	})
}
