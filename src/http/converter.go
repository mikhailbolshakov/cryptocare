package http

import (
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
)

func (c *controllerIml) toBidsApi(bids []*domain.Bid) []*Bid {
	var r []*Bid
	for _, b := range bids {
		r = append(r, &Bid{
			Id:           b.Id,
			Type:         b.Type,
			SrcAsset:     b.SrcAsset,
			TrgAsset:     b.TrgAsset,
			Rate:         b.Rate,
			ExchangeCode: b.ExchangeCode,
			Available:    b.Available,
			MinLimit:     b.MinLimit,
			MaxLimit:     b.MaxLimit,
			Methods:      b.Methods,
			UserId:       b.UserId,
			Link:         b.Link,
		})
	}
	return r
}

func (c *controllerIml) toProfitableChainApi(ch *domain.ProfitableChain) *ProfitableChain {
	if ch == nil {
		return nil
	}
	return &ProfitableChain{
		Id:            ch.Id,
		Asset:         ch.Asset,
		ProfitShare:   ch.ProfitShare,
		Methods:       ch.Methods,
		BidAssets:     ch.BidAssets,
		Depth:         ch.Depth,
		ExchangeCodes: ch.ExchangeCodes,
		Bids:          c.toBidsApi(ch.Bids),
		CreatedAt:     ch.CreatedAt,
	}
}

func (c *controllerIml) toProfitableChainsApi(chains []*domain.ProfitableChain) *ProfitableChains {
	r := &ProfitableChains{}
	for _, ch := range chains {
		r.Chains = append(r.Chains, c.toProfitableChainApi(ch))
	}
	return r
}

func (c *controllerIml) toLoginRequest(rq *LoginRequest) *auth.LoginRequest {
	if rq == nil {
		return nil
	}
	return &auth.LoginRequest{
		Username: rq.Email,
		Password: rq.Password,
	}
}

func (c *controllerIml) toTokenApi(t *auth.SessionToken) *SessionToken {
	if t == nil {
		return nil
	}
	return &SessionToken{
		SessionId:             t.SessionId,
		AccessToken:           t.AccessToken,
		AccessTokenExpiresAt:  t.AccessTokenExpiresAt,
		RefreshToken:          t.RefreshToken,
		RefreshTokenExpiresAt: t.RefreshTokenExpiresAt,
	}
}

func (c *controllerIml) toLoginResponseApi(s *auth.Session, t *auth.SessionToken) *LoginResponse {
	return &LoginResponse{
		UserId: s.UserId,
		Token:  c.toTokenApi(t),
	}
}

func (c *controllerIml) toClientRegRequestDomain(rq *ClientRegistrationRequest) *auth.User {
	return &auth.User{
		Username:  rq.Email,
		Password:  rq.Password,
		Type:      domain.UserTypeClient,
		FirstName: rq.FirstName,
		LastName:  rq.LastName,
		Groups:    []string{domain.AuthGroupClient},
	}
}

func (c *controllerIml) toClientUserApi(usr *auth.User) *ClientUser {
	return &ClientUser{
		Id:        usr.Id,
		Email:     usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	}
}

func (c *controllerIml) toSubscriptionFilterDomain(f *SubscriptionChainFilter) *domain.SubscriptionChainFilter {
	if f == nil {
		return nil
	}
	return &domain.SubscriptionChainFilter{
		Assets:    f.Assets,
		Methods:   f.Methods,
		Exchanges: f.Exchanges,
		MaxDepth:  f.MaxDepth,
		MinProfit: f.MinProfit,
	}
}

func (c *controllerIml) toSubscriptionNotificationsRequestDomain(n []*SubscriptionNotificationRequest) []*domain.SubscriptionNotification {
	var r []*domain.SubscriptionNotification
	for _, nn := range n {
		r = append(r, &domain.SubscriptionNotification{
			Channel:  domain.SubscriptionNotificationChannelTelegram,
			IsActive: nn.IsActive,
			Telegram: &domain.SubscriptionTelegramNotificationDetails{
				Channel: nn.TelegramChannel,
			},
		})
	}
	return r
}

func (c *controllerIml) toCreateSubscriptionRequestDomain(rq *SubscriptionRequest, userId string) *domain.Subscription {
	if rq == nil {
		return nil
	}
	if rq.Filter == nil {
		rq.Filter = &SubscriptionChainFilter{}
	}
	return &domain.Subscription{
		UserId:        userId,
		Filter:        c.toSubscriptionFilterDomain(rq.Filter),
		Notifications: c.toSubscriptionNotificationsRequestDomain(rq.Notifications),
	}
}

func (c *controllerIml) toSubscriptionFilterApi(f *domain.SubscriptionChainFilter) *SubscriptionChainFilter {
	if f == nil {
		return nil
	}
	return &SubscriptionChainFilter{
		Assets:    f.Assets,
		Methods:   f.Methods,
		Exchanges: f.Exchanges,
		MaxDepth:  f.MaxDepth,
		MinProfit: f.MinProfit,
	}
}

func (c *controllerIml) toSubscriptionNotificationsApi(nn []*domain.SubscriptionNotification) []*SubscriptionNotification {
	var r []*SubscriptionNotification
	for _, n := range nn {
		notify := &SubscriptionNotification{
			Id:       n.Id,
			Channel:  n.Channel,
			IsActive: n.IsActive,
		}
		if n.Channel == domain.SubscriptionNotificationChannelTelegram {
			notify.Telegram = &SubscriptionTelegramNotificationDetails{
				Channel: n.Telegram.Channel,
			}
		}
		r = append(r, notify)
	}
	return r
}

func (c *controllerIml) toSubscriptionApi(subs *domain.Subscription) *Subscription {
	if subs == nil {
		return nil
	}
	return &Subscription{
		Id:            subs.Id,
		UserId:        subs.UserId,
		IsActive:      subs.IsActive,
		Filter:        c.toSubscriptionFilterApi(subs.Filter),
		Notifications: c.toSubscriptionNotificationsApi(subs.Notifications),
	}
}

func (c *controllerIml) toSubscriptionsApi(subs []*domain.Subscription) []*Subscription {
	var r []*Subscription
	for _, s := range subs {
		r = append(r, c.toSubscriptionApi(s))
	}
	return r
}

func (c *controllerIml) toBidRequestDomain(rq *BidRequest) *domain.Bid {
	if rq == nil {
		return nil
	}
	return &domain.Bid{
		Id:           rq.Id,
		SrcAsset:     rq.SrcAsset,
		TrgAsset:     rq.TrgAsset,
		Rate:         rq.Rate,
		ExchangeCode: rq.ExchangeCode,
		Available:    rq.Available,
		MinLimit:     rq.MinLimit,
		MaxLimit:     rq.MaxLimit,
		Methods:      rq.Methods,
		UserId:       rq.UserId,
		Link:         rq.Link,
	}
}

func (c *controllerIml) toBidApi(bid *domain.Bid) *Bid {
	if bid == nil {
		return nil
	}
	return &Bid{
		Id:           bid.Id,
		Type:         bid.Type,
		SrcAsset:     bid.SrcAsset,
		TrgAsset:     bid.TrgAsset,
		Rate:         bid.Rate,
		ExchangeCode: bid.ExchangeCode,
		Available:    bid.Available,
		MinLimit:     bid.MinLimit,
		MaxLimit:     bid.MaxLimit,
		Methods:      bid.Methods,
		UserId:       bid.UserId,
		Link:         bid.Link,
	}
}
