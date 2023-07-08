package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"io/ioutil"
	"net/http"
)

// Telegram sends message to telegram channel
type Telegram interface {
	// Send sends a text message
	Send(ctx context.Context, bot, text string, channel int) error
}

type telegramImpl struct {
	logger log.CLoggerFunc
}

func NewTelegram(logger log.CLoggerFunc) Telegram {
	return &telegramImpl{
		logger: logger,
	}
}

func (t *telegramImpl) l() log.CLogger {
	return t.logger().Cmp("telegram")
}

func (t *telegramImpl) Send(ctx context.Context, bot, text string, channel int) error {
	l := t.l().C(ctx).Mth("send").F(log.FF{"channel": channel}).Trc(text)

	// check config
	if bot == "" {
		return ErrTelegramBotEmpty(ctx)
	}

	// send request
	rs, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?&parse_mode=html&chat_id=%d&text=%s&disable_web_page_preview=True", bot, channel, text))
	if err != nil {
		return ErrTelegramRequestFailed(ctx, err)
	}
	defer func() { _ = rs.Body.Close() }()

	// read and unmarshal response
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return ErrTelegramRequestFailed(ctx, err)
	}

	// check response
	rsBody := map[string]interface{}{}
	_ = json.Unmarshal(body, &rsBody)
	if rs.StatusCode < 300 && rsBody != nil {
		if okRes, ok := rsBody["ok"]; ok {
			if r, ok := okRes.(bool); ok && r {
				l.Trc("ok")
				return nil
			}
		}
	}
	return ErrTelegramResponseError(ctx, rs.Status, string(body))
}
