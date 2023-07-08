//go:build integration
// +build integration

package aerospike

import (
	"context"
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var logger = log.Init(&log.Config{Level: log.TraceLevel})
var logf = func() log.CLogger {
	return log.L(logger)
}

type exchange struct {
	ExchangeCode string   `json:"exchCode"`
	Price        float64  `json:"price"`
	Quantity     float64  `json:"quantity"`
	MinLimit     float64  `json:"minLimit"`
	MaxLimit     float64  `json:"maxLimit"`
	Methods      []string `json:"methods"`
	PartyDetails string   `json:"partyDet"`
}

type bid struct {
	Id       string    `json:"id"`
	Src      string    `json:"src"`
	Trg      string    `json:"trg"`
	Exchange *exchange `json:"exchange"`
}

var ctx = context.Background()

func connect(t *testing.T) Aerospike {
	// open
	aes := New()
	err := aes.Open(ctx, cfg, logf)
	if err != nil {
		t.Fatal(err)
	}
	return aes
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	cfg = &Config{
		Host:      "localhost",
		Port:      3000,
		Namespace: "cryptocare",
	}
	currencies = []string{"RUB", "USD", "EUR", "BTC", "ETH", "SLN", "UAH", "CHY", "USDT", "USDC"}
	exchanges  = []string{"binance", "huobi", "bybit"}
)

func newBid() *bid {
	iSrc := rand.Int31n(int32(len(currencies)))
	iTrg := rand.Int31n(int32(len(currencies)))
	iExch := rand.Int31n(int32(len(exchanges)))
	return &bid{
		Id:  kit.NewId(),
		Src: currencies[iSrc],
		Trg: currencies[iTrg],
		Exchange: &exchange{
			ExchangeCode: exchanges[iExch],
			Price:        float64(rand.Int63n(1000000000) / 100),
			Quantity:     float64(rand.Int63n(10000) / 100),
			MinLimit:     10.5,
			MaxLimit:     1000.5,
			Methods:      []string{"M1", "M2", "M3"},
			PartyDetails: kit.NewId(),
		},
	}
}

func genMany(n int) []*bid {
	r := make([]*bid, n)
	for i := 0; i < n; i++ {
		r[i] = newBid()
	}
	return r
}

func Test_CRUD(t *testing.T) {

	aes := connect(t)
	defer aes.Close(ctx)

	b := &bid{
		Id:  kit.NewId(),
		Src: "USD",
		Trg: "ETH",
		Exchange: &exchange{
			ExchangeCode: "binance",
			Price:        1000.5,
			Quantity:     100,
			MinLimit:     10,
			MaxLimit:     100000,
			Methods:      []string{"M1", "M2", "M3"},
			PartyDetails: kit.NewId(),
		},
	}
	bj, _ := json.Marshal(b)
	putKey, err := aero.NewKey(cfg.Namespace, "", b.Id)
	if err != nil {
		t.Fatal(err)
	}
	bins := aero.BinMap{"bid": bj}
	// write the bins
	err = aes.Instance().Put(nil, putKey, bins)
	if err != nil {
		t.Fatal(err)
	}

	getKey, err := aero.NewKey(cfg.Namespace, "", b.Id)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := aes.Instance().Get(nil, getKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, rec.Bins)
	assert.NotEmpty(t, rec.Bins["bid"])

	v := rec.Bins["bid"].([]byte)
	actual := &bid{}
	_ = json.Unmarshal(v, &actual)

	assert.Equal(t, actual, b)

	_, err = aes.Instance().Delete(nil, getKey)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Put_GetBatch(t *testing.T) {
	aes := connect(t)
	defer aes.Close(ctx)

	bids := genMany(1000)
	keys := make([]*aero.Key, len(bids))

	// put
	writePolicy := aero.NewWritePolicy(0, 10)
	for i, b := range bids {
		bj, _ := json.Marshal(b)
		putKey, err := aero.NewKey(cfg.Namespace, "", b.Id)
		if err != nil {
			t.Fatal(err)
		}
		keys[i] = putKey
		bin := aero.NewBin("bid", bj)
		// write the bins
		err = aes.Instance().PutBins(writePolicy, putKey, bin)
		if err != nil {
			t.Fatal(err)
		}
	}

	batchPolicy := aero.NewBatchPolicy()
	records, err := aes.Instance().BatchGet(batchPolicy, keys, "bid")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(bids), len(records))

}

func Test_Put_ScanAll(t *testing.T) {
	aes := connect(t)
	defer aes.Close(ctx)

	bids := genMany(100)
	keys := make([]*aero.Key, len(bids))

	// put
	set := kit.NewRandString()[:10]
	writePolicy := aero.NewWritePolicy(0, 10)
	for i, b := range bids {
		bj, _ := json.Marshal(b)
		putKey, err := aero.NewKey(cfg.Namespace, set, b.Id)
		if err != nil {
			t.Fatal(err)
		}
		keys[i] = putKey
		bin := aero.NewBin("bin", bj)
		// write the bins
		err = aes.Instance().PutBins(writePolicy, putKey, bin)
		if err != nil {
			t.Fatal(err)
		}
	}

	scanPolicy := aero.NewScanPolicy()
	recordSet, err := aes.Instance().ScanAll(scanPolicy, cfg.Namespace, set, "bin")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, recordSet)

	var r []aero.BinMap

	for rec := range recordSet.Results() {
		if rec.Err != nil {
			// if there was an error, stop
			t.Fatal(rec.Err)
		}
		r = append(r, rec.Record.Bins)
	}
	assert.Equal(t, len(bids), len(r))
	for _, rec := range r {
		v := rec["bin"].([]byte)
		actual := &bid{}
		_ = json.Unmarshal(v, &actual)
		assert.NotEmpty(t, actual)
		assert.NotEmpty(t, actual.Id)
		assert.NotEmpty(t, actual.Src)
		assert.NotEmpty(t, actual.Exchange)
	}

}

func Test_ExpList_Contains(t *testing.T) {
	aes := connect(t)
	defer aes.Close(ctx)

	type obj struct {
		Id     string
		Values []string
	}
	v1 := kit.NewRandString()
	v2 := kit.NewRandString()
	v3 := kit.NewRandString()
	v4 := kit.NewRandString()
	values := []*obj{
		{
			Id:     kit.NewRandString(),
			Values: []string{v1, v2},
		},
		{
			Id:     kit.NewRandString(),
			Values: []string{v2, v3},
		},
		{
			Id:     kit.NewRandString(),
			Values: []string{v3, v4, v1, v2},
		},
	}

	// put
	writePolicy := aero.NewWritePolicy(0, 60)
	writePolicy.SendKey = true
	for _, v := range values {
		putKey, err := aero.NewKey(cfg.Namespace, "test_values", v.Id)
		if err != nil {
			t.Fatal(err)
		}
		err = aes.Instance().PutBins(nil, putKey, aero.NewBin("values", v.Values))
		if err != nil {
			t.Fatal(err)
		}
	}

	queryPolicy := aero.NewQueryPolicy()
	queryPolicy.SendKey = true
	queryPolicy.FilterExpression =
		aero.ExpGreater(
			aero.ExpListGetByValueList(
				aero.ListReturnTypeCount,
				aero.ExpListValueVal(v3, v4),
				aero.ExpListBin("values"),
			),
			aero.ExpIntVal(0),
		)

	statement := aero.NewStatement(cfg.Namespace, "test_values")

	recordSet, err := aes.Instance().Query(queryPolicy, statement)
	if err != nil {
		t.Fatal(err)
	}
	var res []*obj
	for r := range recordSet.Results() {
		if r.Err != nil {
			t.Fatal(err)
		} else {
			v, err := AsStrings(ctx, r.Record.Bins, "values")
			if err != nil {
				t.Fatal(err)
			}
			res = append(res, &obj{
				//Id:     r.Record.Key.Value().String(),
				Values: v},
			)
		}
	}
	assert.NotEmpty(t, res)
}
