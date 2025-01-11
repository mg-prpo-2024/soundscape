package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, config *Config) {
	service := NewService(NewRepository(db))
	stripe.Key = config.StripeSecretKey
	registerCreateUser(api, service, config)
	registerGetUser(api, service)
	// https://docs.stripe.com/billing/subscriptions/build-subscriptions?platform=web&ui=checkout#test
	registerCreateSubscription(api, service, config)
	registerStripeHook(api, service)
	registerGetCustomer(api, service)
	registerGetPayments(api, service)
}

type CreateUserInput struct {
	Secret string `header:"X-Auth0-Webhook-Secret" doc:"Auth0 Webhook Secret"`
	Body   CreateUserDto
}

type CreateUserOutput struct{}

func registerCreateUser(api huma.API, service Service, opts *Config) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "Create a user",
		Description: "A webhook endpoint called when a user is created in Auth0.",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"Auth0WebhookSecret": []string{}},
		},
	}, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
		if input.Secret != opts.Auth0HookSecret {
			return nil, huma.Error401Unauthorized("Invalid webhook secret")
		}
		err := service.CreateUser(input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateUserOutput{}, nil
	})
}

type GetUserInput struct {
	Id string `path:"id" doc:"User ID"`
}

type GetUserOutput struct {
	Body UserDto
}

func registerGetUser(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get a user",
		Description: "Get a user by ID.",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
		user, err := service.GetUser(input.Id)
		if err != nil {
			return nil, err
		}
		return &GetUserOutput{Body: *user}, nil
	})
}

type CreateSubscriptionInput struct {
	Body struct {
		UserId     string `json:"user_id" example:"123123123" doc:"User ID"`
		PlanId     string `json:"plan_id" example:"premium" doc:"Plan ID"`
		SuccessUrl string `json:"success_url" example:"http://example.com/success" doc:"Success URL"`
		CancelUrl  string `json:"cancel_url" example:"http://example.com/cancel" doc:"Cancel URL"`
	}
}

type CreateSubscriptionOutput_Body struct {
	Id  string `json:"id" example:"sub_xxxx" doc:"Subscription ID"`
	URL string `json:"url" example:"https://checkout.stripe.com/..." doc:"Stripe Checkout URL"`
}
type CreateSubscriptionOutput struct {
	Body CreateSubscriptionOutput_Body
}

func registerCreateSubscription(api huma.API, service Service, config *Config) {
	huma.Register(api, huma.Operation{
		OperationID: "create-subscription",
		Method:      http.MethodPost,
		Path:        "/subscriptions",
		Summary:     "Create a subscription",
		Description: "Create a stripe subscription.",
		Tags:        []string{"Subscriptions"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateSubscriptionInput) (*CreateSubscriptionOutput, error) {
		body := input.Body
		switch body.PlanId {
		case "premium":
			body.PlanId = config.PlanStripeId.Premium
		case "student":
			body.PlanId = config.PlanStripeId.Student
		default:
			return nil, huma.Error422UnprocessableEntity("Invalid plan ID")
		}
		subscription, err := service.CreateSubscription(body.PlanId, body.UserId, body.SuccessUrl, body.CancelUrl)
		if err != nil {
			return nil, huma.Error500InternalServerError("Error creating subscription", err)
		}
		return &CreateSubscriptionOutput{
			Body: CreateSubscriptionOutput_Body{Id: subscription.ID, URL: subscription.URL},
		}, nil
	})
}

type WebhookInput struct {
	Event *stripe.Event
}

func (wi *WebhookInput) Resolve(ctx huma.Context) []error {
	b, err := io.ReadAll(ctx.BodyReader())
	if err != nil {
		logrus.Errorf("io.ReadAll error: %v", err)
		return []error{&huma.ErrorDetail{
			Message: "Unable to read request body.",
		}}
	}
	// for local, run `make stripe-webhook-forward`, which will provide the secret
	webhookSecret := os.Getenv("SERVICE_STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(b, ctx.Header("Stripe-Signature"), webhookSecret)
	if err != nil {
		logrus.Errorf("webhook.constructevent: %v", err)
		return []error{&huma.ErrorDetail{
			Message: fmt.Sprintf("Unable to construct webhook event: %v", err),
		}}
	}
	wi.Event = &event
	return []error{}
}

func registerStripeHook(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "stripe-webhook",
		Method:      http.MethodPost,
		Path:        "/stripe-webhook",
		Summary:     "Stripe Webhook",
		Description: "Stripe webhook endpoint.",
		Tags:        []string{"Subscriptions"},
	}, func(ctx context.Context, input *WebhookInput) (*struct{}, error) {
		event := *input.Event
		switch event.Type {
		case "checkout.session.completed":
			// Payment is successful and the subscription is created.
			// You should provision the subscription and save the customer ID to your database.
			PrettyPrint(event)
			var checkoutSession stripe.CheckoutSession
			err := json.Unmarshal(event.Data.Raw, &checkoutSession)
			if err != nil {
				return nil, huma.Error422UnprocessableEntity("Unable to parse stripe.CheckoutSession")
			}
			// sets customer ID in users table
			// for now, everything else will be fetched from stripe
			err = service.ProvisionSubscription(&checkoutSession)
			return nil, err
		case "invoice.paid":
			// Continue to provision the subscription as payments continue to be made.
			// Store the status in your database and check when a user accesses your service.
			// This approach helps you avoid hitting rate limits.
		case "invoice.payment_failed":
			// The payment failed or the customer does not have a valid payment method.
			// The subscription becomes past_due. Notify your customer and send them to the
			// customer portal to update their payment information.
		default:
			// unhandled event type
		}
		return nil, nil
	})
}

func PrettyPrint(v interface{}) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Println("Pretty Print Error:", err)
		return
	}
	logrus.Println(string(prettyJSON))
}

type GetCustomerInput struct {
	Id string `path:"id" doc:"User ID" example:"2abe81c8-5b1a-4625-904b-ef4a8594de4c"`
}
type GetCustomerOutput struct {
	Body *stripe.Customer
}

func registerGetCustomer(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-customer",
		Method:      http.MethodGet,
		Path:        "/users/{id}/customer",
		Summary:     "Get stripe customer object for a user",
		Description: "Get stripe customer object for a user by stripe customer id.",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetCustomerInput) (*GetCustomerOutput, error) {
		user, err := service.GetCustomer(input.Id)
		if err != nil {
			return nil, err
		}
		return &GetCustomerOutput{Body: user}, nil
	})
}

type GetPaymentsInput struct {
	Id string `path:"id" doc:"User ID" example:"2abe81c8-5b1a-4625-904b-ef4a8594de4c"`
}
type GetPaymentsOutput struct {
	Body []*stripe.PaymentIntent
}

func registerGetPayments(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-payments",
		Method:      http.MethodGet,
		Path:        "/users/{id}/payments",
		Summary:     "Get stripe payments for a user",
		Description: "Get stripe payments for a user.",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetPaymentsInput) (*GetPaymentsOutput, error) {
		payments, err := service.GetPayments(input.Id)
		if err != nil {
			return nil, err
		}
		return &GetPaymentsOutput{Body: payments}, nil
	})
}
