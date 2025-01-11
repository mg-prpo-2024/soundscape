package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type MetadataRepository interface {
	GetSongs(token string, songIds []string) ([]*MetadataSong, error)
}

type metadataRepository struct {
	client http.Client
}

type MetadataSong struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func NewMetadataRepository() *metadataRepository {
	c := http.Client{Timeout: time.Duration(10) * time.Second}
	return &metadataRepository{client: c}
}

func (r *metadataRepository) GetSongs(token string, songIds []string) ([]*MetadataSong, error) {
	request, err := http.NewRequest("GET", "http://localhost:8000/songs?ids="+strings.Join(songIds, ","), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	songs := []*MetadataSong{}
	if err := json.Unmarshal(data, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}
