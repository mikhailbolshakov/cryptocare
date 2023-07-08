package subscription

import (
	"context"
	"fmt"
	"github.com/mikhailbolshakov/cryptocare/src/domain"
	"github.com/mikhailbolshakov/cryptocare/src/kit/goroutine"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/mikhailbolshakov/cryptocare/src/kit/telegram"
	"github.com/mikhailbolshakov/cryptocare/src/service"
	"strings"
	"time"
)

type TelegramOptions struct {
	Bot string
}

const (
	emojiFlame       = "%F0%9F%94%A5"
	emojiRightArrow  = "%E2%96%B6"
	emojiExclamation = "%E2%9D%97"
	emojiRocket      = "%F0%9F%9A%80"
	newLine          = "%0a"
	numWorkers       = 10
)

type tgSendRequest struct {
	Channel int
	Rq      string
	Bot     string
}

type telegramNotifier struct {
	telegram telegram.Telegram
	opt      *TelegramOptions
	sendChan chan *tgSendRequest
}

func NewTelegramNotifier(telegram telegram.Telegram, opt *TelegramOptions) domain.TelegramNotifier {
	return &telegramNotifier{
		telegram: telegram,
		opt:      opt,
		sendChan: make(chan *tgSendRequest, numWorkers),
	}
}

func (t *telegramNotifier) l() log.CLogger {
	return service.L().Cmp("tg-notifier")
}

func (t *telegramNotifier) Init(ctx context.Context) error {
	t.worker(ctx)
	return nil
}

func (t *telegramNotifier) worker(ctx context.Context) {
	for i := 0; i < numWorkers; i++ {
		goroutine.New().
			WithLogger(t.l().C(ctx).Mth("tg-worker")).
			WithRetry(goroutine.Unrestricted).
			WithRetryDelay(time.Second*10).
			Go(ctx, func() {
				t.l().C(ctx).Mth("tg-worker").InfF("worker %d started", i)
				for rq := range t.sendChan {
					err := t.telegram.Send(ctx, rq.Bot, rq.Rq, rq.Channel)
					if err != nil {
						t.l().C(ctx).Mth("tg-worker").E(err).Err()
					}
				}
			})
	}
}

func (t *telegramNotifier) getExchangeTags(chain *domain.ProfitableChain) string {
	b := strings.Builder{}
	for i, ch := range chain.ExchangeCodes {
		b.WriteString("%23")
		b.WriteString(ch)
		if i < len(ch)-1 {
			b.WriteString(" ")
		}
	}
	return b.String()
}

func (t *telegramNotifier) getBids(chain *domain.ProfitableChain) string {
	b := strings.Builder{}
	bidsLen := len(chain.Bids)
	for i, bid := range chain.Bids {
		b.WriteString(bid.SrcAsset)
		b.WriteString(":")
		b.WriteString(bid.TrgAsset)
		b.WriteString("(")
		b.WriteString(bid.ExchangeCode)
		b.WriteString(", ")
		b.WriteString(fmt.Sprintf("%.5f", bid.Rate))
		b.WriteString(")")
		if i < bidsLen-1 {
			b.WriteString(" -> ")
		}
	}
	return b.String()
}

func (t *telegramNotifier) getProfitClass(chain *domain.ProfitableChain) int {
	if chain.ProfitShare < 1.02 {
		return 1
	} else if chain.ProfitShare < 1.05 {
		return 2
	} else if chain.ProfitShare < 1.1 {
		return 3
	} else if chain.ProfitShare < 1.2 {
		return 4
	} else {
		return 5
	}
}

func (t *telegramNotifier) getTags(chain *domain.ProfitableChain) string {
	b := strings.Builder{}
	b.WriteString("%23")
	b.WriteString(chain.Asset)
	b.WriteString(" ")
	b.WriteString(t.getExchangeTags(chain))
	b.WriteString(" %23P")
	b.WriteString(fmt.Sprintf("%d", t.getProfitClass(chain)))
	return b.String()
}

func (t *telegramNotifier) emoji(emojis ...string) string {
	b := strings.Builder{}
	for _, e := range emojis {
		b.WriteString(e)
	}
	return b.String()
}

func (t *telegramNotifier) getProfitClassEmoji(chain *domain.ProfitableChain) string {
	switch t.getProfitClass(chain) {
	case 2:
		return t.emoji(emojiFlame)
	case 3:
		return t.emoji(emojiFlame, emojiFlame, emojiFlame)
	case 4:
		return t.emoji(emojiExclamation, emojiExclamation, emojiFlame)
	case 5:
		return t.emoji(emojiRocket, emojiRocket, emojiExclamation, emojiFlame)
	}
	return ""
}

func (t *telegramNotifier) getRequest(chain *domain.ProfitableChain) string {
	b := strings.Builder{}
	b.WriteString(t.getTags(chain))
	b.WriteString(newLine)
	b.WriteString("asset: ")
	b.WriteString(fmt.Sprintf("<b>%s</b>", chain.Asset))
	b.WriteString(newLine)
	b.WriteString(t.getProfitClassEmoji(chain))
	b.WriteString("profit: ")
	b.WriteString(fmt.Sprintf("<b>%.2f%%</b>", (chain.ProfitShare-1)*100))
	b.WriteString(newLine)
	b.WriteString("chain: ")
	b.WriteString(t.getBids(chain))
	b.WriteString(newLine)
	b.WriteString("time: ")
	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(newLine)
	b.WriteString(emojiRightArrow)
	b.WriteString(fmt.Sprintf("<a href='https://panel.cryptocare.ai/trading/details/%s'>link to details</a>", chain.Id))
	b.WriteString(newLine)
	return b.String()
}

func (t *telegramNotifier) Notify(ctx context.Context, bot string, channels []int, chains []*domain.ProfitableChain) error {
	for _, chain := range chains {
		rq := t.getRequest(chain)
		for _, channel := range channels {
			t.sendChan <- &tgSendRequest{
				Channel: channel,
				Rq:      rq,
				Bot:     bot,
			}
		}
	}
	return nil
}
