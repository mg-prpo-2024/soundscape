package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type MetadataRepository interface {
	GetSongs(token string, songIds []string) ([]*MetadataSong, error)
}

type metadataRepository struct {
	client  http.Client
	baseUrl string
}

type MetadataSong struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func NewMetadataRepository(baseUrl string) *metadataRepository {
	c := http.Client{Timeout: time.Duration(10) * time.Second}
	return &metadataRepository{client: c, baseUrl: baseUrl}
}

func (r *metadataRepository) GetSongs(token string, songIds []string) ([]*MetadataSong, error) {
	url := fmt.Sprintf("%s/songs?ids=%s", r.baseUrl, strings.Join(songIds, ","))
	request, err := http.NewRequest("GET", url, nil)
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
