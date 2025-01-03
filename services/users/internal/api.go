package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stripe/stripe-go/v81"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, options *Options) {
	service := NewService(NewRepository(db))
	stripe.Key = options.StripeSecretKey
	registerSignIn(api, service)
	registerSubscribe(api, service)
}

type SignInInput struct {
	Body struct {
		Name string `json:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}
}

type SignInOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func registerSignIn(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "sign-in",
		Method:      http.MethodPost,
		Path:        "/sign-in",
		Summary:     "Sign in",
		Description: "Sign in to the system.",
		Tags:        []string{"Sign In"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *SignInInput) (*SignInOutput, error) {
		resp := &SignInOutput{}
		// token := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		// fmt.Printf("token %+v", token)
		// err := Login("", token.RegisteredClaims.Subject)
		// if err != nil {
		// 	return nil, err
		// }
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Body.Name)
		// service.Login(token.Raw, token.RegisteredClaims.Subject)
		return resp, nil
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
		Path:        "/subscribe",
		Summary:     "Create a subscription",
		Description: "Create a stripe subscription.",
		Tags:        []string{"Subscription"},
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

// TODO: custom serializer for webhook event

// type WebhookInput struct {}

// func registerStripeHook(api huma.API, service Service) {
// 	huma.Register(api, huma.Operation{
// 		OperationID: "stripe-hook",
// 		Method:      http.MethodPost,
// 	}, func(ctx context.Context, input *interface{}) (*interface{}, error) {
//   b, err := io.ReadAll(r.Body)
//   if err != nil {
//     http.Error(w, err.Error(), http.StatusBadRequest)
//     log.Printf("ioutil.ReadAll: %v", err)
//     return
//   }

//   event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), "{{STRIPE_WEBHOOK_SECRET}}")
//   if err != nil {
//     http.Error(w, err.Error(), http.StatusBadRequest)
//     log.Printf("webhook.ConstructEvent: %v", err)
//     return
//   }

//   switch event.Type {
//     case "checkout.session.completed":
//       // Payment is successful and the subscription is created.
//       // You should provision the subscription and save the customer ID to your database.
//     case "invoice.paid":
//       // Continue to provision the subscription as payments continue to be made.
//       // Store the status in your database and check when a user accesses your service.
//       // This approach helps you avoid hitting rate limits.
//     case "invoice.payment_failed":
//       // The payment failed or the customer does not have a valid payment method.
//       // The subscription becomes past_due. Notify your customer and send them to the
//       // customer portal to update their payment information.
//     default:
//       // unhandled event type
//   }
// 		return nil, nil
// 	})
// }
