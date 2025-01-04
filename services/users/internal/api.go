package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, options *Options) {
	service := NewService(NewRepository(db))
	stripe.Key = options.StripeSecretKey
	registerCreateUser(api, service)
	registerGetUser(api, service)
	// https://docs.stripe.com/billing/subscriptions/build-subscriptions?platform=web&ui=checkout#test
	registerSubscribe(api, service)
	registerStripeHook(api, service)
	// registerGetPayments(api, service)
	// registerGetCustomer(api, service)
}

type CreateUserInput struct {
	Body UserDto
}

type CreateUserOutput struct{}

func registerCreateUser(api huma.API, service Service) {
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
		Summary:     "Create a user",
		Description: "A webhook endpoint called when a user is created in Auth0.",
		Tags:        []string{"Users"},
		Security: []map[string][]string{
			{"Auth0WebhookSecret": []string{}},
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
		PlanId     string `json:"plan_id" example:"prod_1234" doc:"Plan ID"`
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

func registerSubscribe(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "subscribe",
		Method:      http.MethodPost,
		Path:        "/subscriptions",
		Summary:     "Create a subscription",
		Description: "Create a stripe subscription.",
		Tags:        []string{"Subscriptions"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateSubscriptionInput) (*CreateSubscriptionOutput, error) {
		fmt.Printf("input %+v\n", input.Body)
		body := input.Body
		subscription, err := service.Subscribe(body.PlanId, body.UserId, body.SuccessUrl, body.CancelUrl)
		if err != nil {
			return nil, huma.Error500InternalServerError("Error creating subscription", err)
		}
		fmt.Printf("subscription %#+v\n", subscription)
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
		log.Printf("io.ReadAll error: %v", err)
		return []error{&huma.ErrorDetail{
			Message: "Unable to read request body.",
		}}
	}
	// for local, run `make stripe-webhook-forward`, which will provide the secret
	webhookSecret := os.Getenv("SERVICE_STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(b, ctx.Header("Stripe-Signature"), webhookSecret)
	if err != nil {
		log.Printf("webhook.constructevent: %v", err)
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
		fmt.Println("Pretty Print Error:", err)
		return
	}
	fmt.Println(string(prettyJSON))
}
