package domain

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
)

type GetBidsRequest struct {
	kit.PagingRequest
	SrcAsset string
}

type GetBidsResponse struct {
	kit.PagingResponse
	Bids []*Bid
}

// BidStorage provides an access to bids storage
type BidStorage interface {
	// GetBidLightsBySourceAsset returns bods by the source asset
	GetBidsLightAll(ctx context.Context) ([]*BidLight, error)
	// GetBidsByIds retrieves full bids by Ids
	GetBidsByIds(ctx context.Context, ids []string) ([]*Bid, error)
	// PutBids puts bids
	PutBids(ctx context.Context, bids []*Bid, ttlSec uint32) error
}

// BidStorage provides an access to order storage
type ChainStorage interface {
	// SaveProfitableChains save profitable chains to store
	SaveProfitableChains(ctx context.Context, chains []*ProfitableChain) error
	// GetProfitableChains retrieves stored profitable chains
	GetProfitableChains(ctx context.Context, rq *GetProfitableChainsRequest) (*GetProfitableChainsResponse, error)
	// GetProfitableChain retrieves stored profitable chain by id
	GetProfitableChain(ctx context.Context, chainId string) (*ProfitableChain, error)
	// ProfitableChainExists checks if profitable chain exists
	ProfitableChainExists(ctx context.Context, chainId string) (bool, error)
}

// UserStorage manages user storage
type UserStorage interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, u *auth.User) error
	// UpdateUser updates an user
	UpdateUser(ctx context.Context, u *auth.User) error
	// GetByEmail retrieves an user by email
	GetByUsername(ctx context.Context, email string) (*auth.User, error)
	// GetUser retrieves a user by id
	GetUser(ctx context.Context, userId string) (*auth.User, error)
	// GetByIds retrieves users by IDs
	GetUserByIds(ctx context.Context, userIds []string) ([]*auth.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, u *auth.User) error
}

// SubscriptionStorage manages subscription storage
type SubscriptionStorage interface {
	// SaveSubscription creates or updates a subscription
	SaveSubscription(ctx context.Context, subs *Subscription) error
	// GetSubscription retrieves subscription by id
	GetSubscription(ctx context.Context, subsId string) (*Subscription, error)
	// DeleteSubscription deletes subscription by id
	DeleteSubscription(ctx context.Context, subsId string) error
	// SearchSubscriptions searches subscriptions
	SearchSubscriptions(ctx context.Context, rq *SearchSubscriptionsRequest) ([]*Subscription, error)
}
