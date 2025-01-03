package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Service interface {
	Login(accessToken string, userSub string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Login(accessToken string, userSub string) error {
	domain := os.Getenv("AUTH0_DOMAIN")

	userDetailsByIdUrl := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userSub)
	req, err := http.NewRequest("GET", userDetailsByIdUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	fmt.Println("Response:", string(body))
	return nil
}
