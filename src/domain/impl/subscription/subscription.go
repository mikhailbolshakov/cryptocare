package subscription

import (
	"context"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"strings"
)

type subscriptionSvcImpl struct {
	storage          domain.SubscriptionStorage
	telegramNotifier domain.TelegramNotifier
	cfg              *service.Config
}

func NewSubscriptionService(storage domain.SubscriptionStorage, telegramNotifier domain.TelegramNotifier) domain.SubscriptionService {
	return &subscriptionSvcImpl{
		storage:          storage,
		telegramNotifier: telegramNotifier,
	}
}

func (s *subscriptionSvcImpl) l() log.CLogger {
	return service.L().Cmp("subscription-svc")
}

func (s *subscriptionSvcImpl) Init(cfg *service.Config) {
	s.cfg = cfg
}

func (s *subscriptionSvcImpl) validateAndPopulate(ctx context.Context, subscription *domain.Subscription) error {

	// validate and populate filters
	if subscription.Filter == nil {
		subscription.Filter = &domain.SubscriptionChainFilter{}
	}

	for i, exchange := range subscription.Filter.Exchanges {
		subscription.Filter.Exchanges[i] = strings.ToLower(strings.TrimSpace(exchange))
	}
	for i, m := range subscription.Filter.Methods {
		subscription.Filter.Methods[i] = strings.TrimSpace(m)
	}
	for i, a := range subscription.Filter.Assets {
		subscription.Filter.Assets[i] = strings.ToUpper(strings.TrimSpace(a))
	}

	if subscription.Filter.MinProfit != 0.0 && (subscription.Filter.MinProfit < 0.0001 || subscription.Filter.MinProfit > 99.9999) {
		return errors.ErrSubscriptionMinProfitInvalid(ctx)
	}
	if subscription.Filter.MaxDepth != 0 && subscription.Filter.MaxDepth < 2 {
		return errors.ErrSubscriptionMaxDepthInvalid(ctx)
	}

	for _, notify := range subscription.Notifications {
		if notify.Channel != domain.SubscriptionNotificationChannelTelegram {
			return errors.ErrSubscriptionNotificationChannelNotSupported(ctx, notify.Channel)
		}
		if notify.Id == "" {
			notify.Id = kit.NewRandString()
		}
		if notify.Telegram == nil || notify.Telegram.Channel == 0 {
			return errors.ErrSubscriptionNotificationTelegramInvalid(ctx)
		}
	}

	return nil
}

func (s *subscriptionSvcImpl) Create(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error) {
	s.l().C(ctx).Mth("create").Trc()

	err := s.validateAndPopulate(ctx, subscription)
	if err != nil {
		return nil, err
	}

	subscription.Id = kit.NewId()
	subscription.IsActive = true

	err = s.storage.SaveSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionSvcImpl) Update(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error) {
	s.l().C(ctx).Mth("update").F(log.FF{"subscriptionId": subscription.Id}).Trc()

	if subscription.Id == "" {
		return nil, errors.ErrSubscriptionIdEmpty(ctx)
	}

	stored, err := s.storage.GetSubscription(ctx, subscription.Id)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, errors.ErrSubscriptionNotFound(ctx)
	}
	if !stored.IsActive {
		return nil, errors.ErrSubscriptionNotActive(ctx)
	}

	err = s.validateAndPopulate(ctx, subscription)
	if err != nil {
		return nil, err
	}

	subscription.UserId = stored.UserId
	subscription.IsActive = true

	err = s.storage.SaveSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionSvcImpl) Delete(ctx context.Context, subscriptionId string) error {
	s.l().C(ctx).Mth("delete").F(log.FF{"subscriptionId": subscriptionId}).Trc()

	if subscriptionId == "" {
		return errors.ErrSubscriptionIdEmpty(ctx)
	}

	stored, err := s.storage.GetSubscription(ctx, subscriptionId)
	if err != nil {
		return err
	}
	if stored == nil {
		return errors.ErrSubscriptionNotFound(ctx)
	}

	return s.storage.DeleteSubscription(ctx, subscriptionId)
}

func (s *subscriptionSvcImpl) Get(ctx context.Context, subscriptionId string) (*domain.Subscription, error) {
	s.l().C(ctx).Mth("delete").F(log.FF{"subscriptionId": subscriptionId}).Trc()
	if subscriptionId == "" {
		return nil, errors.ErrSubscriptionIdEmpty(ctx)
	}
	return s.storage.GetSubscription(ctx, subscriptionId)
}

func (s *subscriptionSvcImpl) Deactivate(ctx context.Context, subscriptionId string) (*domain.Subscription, error) {
	s.l().C(ctx).Mth("deactivate").F(log.FF{"subscriptionId": subscriptionId}).Trc()

	if subscriptionId == "" {
		return nil, errors.ErrSubscriptionIdEmpty(ctx)
	}

	subscription, err := s.storage.GetSubscription(ctx, subscriptionId)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, errors.ErrSubscriptionNotFound(ctx)
	}
	if !subscription.IsActive {
		return nil, errors.ErrSubscriptionNotActive(ctx)
	}

	subscription.IsActive = false

	err = s.storage.SaveSubscription(ctx, subscription)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (s *subscriptionSvcImpl) Search(ctx context.Context, rq *domain.SearchSubscriptionsRequest) ([]*domain.Subscription, error) {
	s.l().C(ctx).Mth("search").Trc()
	return s.storage.SearchSubscriptions(ctx, rq)
}

func (s *subscriptionSvcImpl) Notify(ctx context.Context, chains []*domain.ProfitableChain) error {
	l := s.l().C(ctx).Mth("notify").Trc()

	// get active subscriptions
	subs, err := s.Search(ctx, &domain.SearchSubscriptionsRequest{WithInActive: false})
	if err != nil {
		return err
	}

	// go through chains
	for _, chain := range chains {
		chainMethods := kit.Strings(chain.Methods).Sanitize()
		// for each subscription
		var channels []int
		for _, subs := range subs {
			subFilterMethods := kit.Strings(subs.Filter.Methods).Sanitize()
			if (len(subs.Filter.Exchanges) == 0 || kit.Strings(chain.ExchangeCodes).Subset(subs.Filter.Exchanges)) &&
				(len(subs.Filter.Assets) == 0 || kit.Strings(subs.Filter.Assets).Contains(chain.Asset)) &&
				(len(subFilterMethods) == 0 || chainMethods.Subset(subFilterMethods)) &&
				(subs.Filter.MaxDepth == 0 || chain.Depth <= subs.Filter.MaxDepth) &&
				(subs.Filter.MinProfit == 0.0 || chain.ProfitShare >= 1+subs.Filter.MinProfit*0.01) {
				// for all notifications
				for _, notifier := range subs.Notifications {
					if notifier.IsActive && notifier.Channel == domain.SubscriptionNotificationChannelTelegram {
						channels = append(channels, notifier.Telegram.Channel)
					}
				}
			}
		}
		if len(channels) > 0 {
			l.DbgF("channels: %s", channels)
			if err := s.telegramNotifier.Notify(ctx, s.cfg.Arbitrage.Notification.Telegram.Bot, channels, []*domain.ProfitableChain{chain}); err != nil {
				s.l().C(ctx).Mth("notify").E(err).St().Err()
			}
		}
	}
	return nil
}
