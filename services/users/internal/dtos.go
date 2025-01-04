package internal

type UserDto struct {
	Id               string  `json:"id" example:"google-oauth2|106527689641250451478" doc:"User ID"`
	Email            string  `json:"email" format:"email" example:"test@gmail.com" doc:"User email"`
	StripeCustomerId *string `json:"stripe_customer_id,omitempty" example:"cus_JZ6Z9Z9Z9Z9Z9Z" doc:"Stripe customer ID"`
}
