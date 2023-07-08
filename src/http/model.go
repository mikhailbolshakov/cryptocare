package http

import (
	"time"
)

// Bid is a bid exposed on the exchange
type Bid struct {
	Id           string   `json:"id"`           // Id
	Type         string   `json:"type"`         // Type
	SrcAsset     string   `json:"src"`          // SrcAsset - source asset
	TrgAsset     string   `json:"trg"`          // TrgAsset - target asset
	Rate         float64  `json:"rate"`         // Rate - conversion rate
	ExchangeCode string   `json:"exchangeCode"` // ExchangeCode - exchange code
	Available    float64  `json:"available"`    // Available - available volume
	MinLimit     float64  `json:"minLimit"`     // MinLimit - min limit
	MaxLimit     float64  `json:"maxLimit"`     // MaxLimit - max limit
	Methods      []string `json:"methods"`      // Methods - methods
	UserId       string   `json:"userId"`       // UserId - user who exposes the bid
	Link         string   `json:"link"`         // Link - link to the bid on the exchange
}

// ProfitableChain is a sequence of orders to be exposed to achieve calculated profit
type ProfitableChain struct {
	Id            string    `json:"id"`             // Id - chain Id, calculated as hash from bidIds
	Asset         string    `json:"asset"`          // Asset - the target asset
	ProfitShare   float64   `json:"profitShare"`    // ProfitShare profit share
	Methods       []string  `json:"methods"`        // Methods list of methods (union methods from all bids)
	BidAssets     []string  `json:"bidAssets"`      // BidAssets sequence of asset for each bids like [RUB, USD, USDT]
	Depth         int       `json:"depth"`          // Depth chain depth
	ExchangeCodes []string  `json:"exchangeCodes"`  // ExchangeCodes through all bids
	Bids          []*Bid    `json:"bids,omitempty"` // Bids sequence of bids
	CreatedAt     time.Time `json:"createdAt"`      // CreatedAt - when this chain has been created
}

type ProfitableChains struct {
	Chains []*ProfitableChain `json:"chains"` // Chains
}

type LoginRequest struct {
	Email    string `json:"email"`    // Email - login
	Password string `json:"password"` // Password - password
}

// SessionToken specifies a session token
type SessionToken struct {
	SessionId             string    // SessionId - session ID
	AccessToken           string    // AccessToken
	AccessTokenExpiresAt  time.Time // AccessTokenExpiresAt - when access token expires
	RefreshToken          string    // RefreshToken
	RefreshTokenExpiresAt time.Time // RefreshToken - when refresh token expires
}

type LoginResponse struct {
	Token  *SessionToken `json:"token"`  // Token - auth token must be passed as  "Authorization Bearer" header for all the requests (except ones which don't require authorization)
	UserId string        `json:"userId"` // UserId - ID of account
}

type ClientRegistrationRequest struct {
	Email     string `json:"email"`     // Email - user's email
	Password  string `json:"password"`  // Password - password
	FirstName string `json:"firstName"` // FirstName - user's first name
	LastName  string `json:"lastName"`  // LastName - user's last name
}

type ClientUser struct {
	Id        string `json:"id"`                  // Id - user ID
	Email     string `json:"email"`               // Email - email
	FirstName string `json:"firstName,omitempty"` // FirstName - user's first name
	LastName  string `json:"lastName,omitempty"`  // LastName - user's last name
}

type SetPasswordRequest struct {
	PrevPassword string `json:"prevPassword"` // PrevPassword - current password
	NewPassword  string `json:"newPassword"`  // NewPassword - new password
}

// SubscriptionChainFilter allows conditional subscription
type SubscriptionChainFilter struct {
	Assets    []string `json:"assets,omitempty"`    // Assets filters by assets
	Methods   []string `json:"methods,omitempty"`   // Methods filters by methods
	Exchanges []string `json:"exchanges,omitempty"` // Exchanges filters by exchange codes
	MaxDepth  int      `json:"maxDepth,omitempty"`  // MaxDepth max depth of chains
	MinProfit float64  `json:"minProfit,omitempty"` // MinProfit min profit of chains
}

// SubscriptionTelegramNotificationDetails details of telegram notification
type SubscriptionTelegramNotificationDetails struct {
	Channel int `json:"channel"` // Channel telegram channel
}

// SubscriptionNotification notification details
type SubscriptionNotification struct {
	Id       string                                   `json:"id"`                 // Id notification id
	Channel  string                                   `json:"channel"`            // Channel notification channel
	IsActive bool                                     `json:"isActive"`           // IsActive if notification active
	Telegram *SubscriptionTelegramNotificationDetails `json:"telegram,omitempty"` // Telegram telegram details
}

// SubscriptionNotification notification details
type SubscriptionNotificationRequest struct {
	TelegramChannel int  `json:"tgChannel"`
	IsActive        bool `json:"isActive"`
}

// Subscription subscription
type Subscription struct {
	Id            string                      `json:"id"`                      // Id subscription
	UserId        string                      `json:"userId"`                  // UserId owner of the subscription. Might be empty
	IsActive      bool                        `json:"isActive"`                // IsActive if subscription active
	Filter        *SubscriptionChainFilter    `json:"filter,omitempty"`        // Filter subscription filter
	Notifications []*SubscriptionNotification `json:"notifications,omitempty"` // Notifications notifications
}

type Subscriptions struct {
	Items []*Subscription `json:"items"`
}

type SubscriptionRequest struct {
	Filter        *SubscriptionChainFilter           `json:"filter,omitempty"`        // Filter subscription filter
	Notifications []*SubscriptionNotificationRequest `json:"notifications,omitempty"` // Notifications notifications
}

// Bid is a bid exposed on the exchange
type BidRequest struct {
	Id           string   `json:"id"`           // Id
	SrcAsset     string   `json:"src"`          // SrcAsset - source asset
	TrgAsset     string   `json:"trg"`          // TrgAsset - target asset
	Rate         float64  `json:"rate"`         // Rate - conversion rate
	ExchangeCode string   `json:"exchangeCode"` // ExchangeCode - exchange code
	Available    float64  `json:"available"`    // Available available volume
	MinLimit     float64  `json:"minLimit"`     // MinLimit - minimum limit
	MaxLimit     float64  `json:"maxLimit"`     // MaxLimit - max limit
	Methods      []string `json:"methods"`      // Methods - methods
	UserId       string   `json:"userId"`       // UserId - user who expose the bid
	Link         string   `json:"link"`         // Link - link to the bid
}
