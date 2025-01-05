package internal

import (
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type Service interface {
	CreateUser(user UserDto) error
	GetUser(id string) (*UserDto, error)
	CreateSubscription(priceId, userId, successUrl, cancelUrl string) (*stripe.CheckoutSession, error)
	ProvisionSubscription(session *stripe.CheckoutSession) error
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateUser(user UserDto) error {
	return s.repo.CreateUser(user)
}

func (s *service) GetUser(id string) (*UserDto, error) {
	return s.repo.GetUser(id)
}

func (s *service) CreateSubscription(priceId, userId, successUrl, cancelUrl string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
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
	userId := session.Metadata["userId"]
	return s.repo.SetCustomerId(userId, session.Customer.ID)
}
