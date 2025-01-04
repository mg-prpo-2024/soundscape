package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type Service interface {
	Login(accessToken string, userSub string) error
	Subscribe(planId, userId, successUrl, cancelUrl string) (*stripe.CheckoutSession, error)
	ProvisionSubscription(session *stripe.CheckoutSession) error
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

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

func (s *service) Subscribe(planId, userId, successUrl, cancelUrl string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planId),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successUrl),
		CancelURL:  stripe.String(cancelUrl),
		Metadata: map[string]string{
			"userId": userId,
		},
	}

	subscription, err := session.New(params)
	return subscription, err
}

func (s *service) ProvisionSubscription(session *stripe.CheckoutSession) error {
	return nil
}
