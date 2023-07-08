package domain

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/service"
)

const (
	SubscriptionNotificationChannelTelegram = "telegram"
)

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

// Subscription subscription
type Subscription struct {
	Id            string                      // Id subscription
	UserId        string                      // UserId owner of the subscription. Might be empty
	IsActive      bool                        // IsActive if subscription active
	Filter        *SubscriptionChainFilter    // Filter subscription filter
	Notifications []*SubscriptionNotification // Notifications notifications
}

// SearchSubscriptionsRequest request to retrieve subscription
type SearchSubscriptionsRequest struct {
	WithInActive bool   // WithInActive if true, inactive subscription are also retrieved
	UserId       string // UserId filter by user
}

// SubscriptionService subscription service
type SubscriptionService interface {
	// Notifier implements notifier
	Notifier
	// Init initializes service
	Init(cfg *service.Config)
	// Create creates a new subscription
	Create(ctx context.Context, subscription *Subscription) (*Subscription, error)
	// Update updates a subscription
	Update(ctx context.Context, subscription *Subscription) (*Subscription, error)
	// Delete deletes a subscription
	Delete(ctx context.Context, subscriptionId string) error
	// Get retrieves subscription by id
	Get(ctx context.Context, subscriptionId string) (*Subscription, error)
	// Deactivate deactivates an active subscription
	Deactivate(ctx context.Context, subscriptionId string) (*Subscription, error)
	// Search searches subscriptions
	Search(ctx context.Context, rq *SearchSubscriptionsRequest) ([]*Subscription, error)
}

// TelegramNotifier implements telegram notification
type TelegramNotifier interface {
	// Init initializes
	Init(ctx context.Context) error
	// Notify builds and sends notification
	Notify(ctx context.Context, bot string, channels []int, chains []*ProfitableChain) error
}
