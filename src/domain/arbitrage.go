package domain

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"time"
)

const (
	BidTypeP2P    = "p2p"
	BidTypeSpot   = "spot"
	BidTypeManual = "manual"
)

// Bid is a bid exposed on the exchange
type Bid struct {
	Id           string   `json:"id"`           // Id
	Type         string   `json:"type"`         // Type (p2p, spot)
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

// Bid is a bid exposed on the exchange
type BidLight struct {
	Id        string  `json:"id"`        // Id
	Type      string  `json:"type"`      // Type (p2p, spot)
	SrcAsset  string  `json:"src"`       // SrcAsset - source asset
	TrgAsset  string  `json:"trg"`       // TrgAsset - target asset
	Rate      float64 `json:"rate"`      // Rate - conversion rate
	Available float64 `json:"available"` // Available - available amount of asset
	MinLimit  float64 `json:"minLimit"`  // MinLimit - bid min limit
	MaxLimit  float64 `json:"maxLimit"`  // MaxLimit - bid max limit
}

// CandidateChain is a sequence of bids to be a candidate to profitable chain
type CandidateChain struct {
	BidIds    []string // BidIds sequence of bids to be applied
	Amount    float64  // Amount current amount, used to check chain limits
	TotalRate float64  // TotalRate calculated as multiplication of all rates in Rates
}

// CandidateChains bilk of chains
type CandidateChains struct {
	Chains []*CandidateChain // Chains - chains
}

// ProfitableChain is a sequence of orders to be exposed to achieve calculated profit
type ProfitableChain struct {
	Id            string    // Id - chain Id, calculated as hash from bidIds
	Asset         string    // Asset - the target asset
	ProfitShare   float64   // ProfitShare profit share
	Methods       []string  // Methods list of methods (union methods from all bids)
	BidAssets     []string  // BidAssets sequence of asset for each bids like [RUB, USD, USDT]
	Bids          []*Bid    // Bids sequence of bids
	Depth         int       // Depth chain depth
	ExchangeCodes []string  // ExchangeCodes through all bids
	CreatedAt     time.Time // CreatedAt - when this chain has been created
}

// ProfitableChains bilk of chains
type ProfitableChains struct {
	Chains []*ProfitableChain // Chains - chains
}

// GetProfitableChainsRequest request to retrieve order chains
type GetProfitableChainsRequest struct {
	kit.PagingRequest
	Assets        []string // Assets - retrieves chains by the given assets
	WithBids      bool     // WithBids - if true, retrieve chains with bids
	Methods       []string // Methods - retrieves by methods
	ExchangeCodes []string // ExchangeCodes - retrieves by exchange codes
}

type GetProfitableChainsResponse struct {
	kit.PagingResponse
	Chains []*ProfitableChain
}

// BidProvider provides bids data for analysis
type BidProvider interface {
	// Init
	Init(cfg *service.Config)
	// Run runs workers which refreshes bids data
	Run(ctx context.Context) error
	// Stop stops workers
	Stop(ctx context.Context) error
	// GetAssets returns list of available assets
	GetAssets(ctx context.Context) ([]string, error)
	// GetBidLightsBySourceAsset returns bods by the source asset
	GetBidLightsBySourceAsset(ctx context.Context, srcAsset string) ([]*BidLight, error)
	// GetBidsByIds retrieves full bids by Ids
	GetBidsByIds(ctx context.Context, ids []string) ([]*Bid, error)
	// PutBid puts a manual bid
	PutBid(ctx context.Context, bid *Bid) (*Bid, error)
}

// Notifier responsible for notification users about chains
type Notifier interface {
	// Notify notifies
	Notify(ctx context.Context, chains []*ProfitableChain) error
}

// ArbitrageService provides arbitrage functions
type ArbitrageService interface {
	// Init initializes service
	Init(cfg *service.Config)
	// RunCalculation runs calculation of profitable chains as a background process
	RunCalculationBackground(ctx context.Context) error
	// StopCalculation stops background process if running
	StopCalculation(ctx context.Context) error
	// GetProfitableChains retrieves profitable chains by criteria
	GetProfitableChains(ctx context.Context, rq *GetProfitableChainsRequest) (*GetProfitableChainsResponse, error)
	// GetProfitableChain retrieves profitable chain by id
	GetProfitableChain(ctx context.Context, chainId string) (*ProfitableChain, error)
}

// BidGenerator generates bid data (for test purposes only) // TODO: remove
type BidGenerator interface {
	Init(cfg *service.Config)
	Run(ctx context.Context)
	Stop(ctx context.Context)
}
