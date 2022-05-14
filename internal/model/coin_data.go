package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/nntaoli-project/goex"
)

var pairMap = map[string]goex.CurrencyPair{
	"NEARUSDT": goex.NEAR_USDT,
}

type TradingPair struct {
	PairName string
	api      goex.API
	Pair     goex.CurrencyPair
	RunBet   RunBet `json:"runBet"`
	Config   Config `json:"config"`
}

type RunBet struct {
	NextBuyPrice  float64   `json:"next_buy_price"`  // 网格买入价格
	GridSellPrice float64   `json:"grid_sell_price"` // 网格卖出价格
	Step          int64     `json:"step"`            // 当前步数
	RecordedPrice []float64 `json:"recorded_price"`
}

type Config struct {
	ProfitRatio      float64   `json:"profit_ratio"`
	DoubleThrowRatio float64   `json:"double_throw_ratio"`
	Quantity         []float64 `json:"quantity"` // 买入量
}

type CurrentSpotData struct {
	Pair           goex.CurrencyPair
	GridBuyPrice   float64 // 网格买入价格
	GridSellPrice  float64 // 网格卖出价格
	Quantity       float64 // 买入量
	CurMarketPrice float64 // 市场现价
	Step           int64
	RightSize      string
}

func (t *TradingPair) Build(pairName string, api goex.API) {
	t.PairName = pairName
	t.Pair = pairMap[pairName]
	t.api = api
}

func (t *TradingPair) CurrentData() (*CurrentSpotData, error) {
	ticker, err := t.api.GetTicker(t.Pair)
	if err != nil {
		return nil, err
	}
	return &CurrentSpotData{
		Pair:           t.Pair,
		GridBuyPrice:   t.RunBet.NextBuyPrice,
		GridSellPrice:  t.RunBet.GridSellPrice,
		CurMarketPrice: ticker.Last,
		Quantity:       t.GetQuantity(true),
		Step:           t.RunBet.Step,
		RightSize:      strings.Split(strconv.FormatFloat(ticker.Last, 'E', -1, 64), ".")[1],
	}, nil
}

func (t *TradingPair) GetQuantity(isBuy bool) float64 {
	curStep := t.RunBet.Step
	if !isBuy {
		curStep = curStep - 1
	}
	quantityArr := t.Config.Quantity
	quantity := float64(0)
	if int(curStep) < len(quantityArr) {
		quantity = quantityArr[0]
		if curStep >= 0 {
			quantity = quantityArr[curStep]
		}
	} else {
		quantity = quantityArr[len(quantityArr)-1]
	}
	return quantity
}

// SetRatio 修改补仓止盈比率
func (t *TradingPair) SetRatio() {
	klineNum := 20
	klineList, err := t.api.GetKlineRecords(t.Pair, goex.KLINE_PERIOD_4H, klineNum)
	if err != nil {
		return
	}

	var percentTotal float64
	for _, kline := range klineList {
		percentTotal += math.Abs(kline.High-kline.Low) / kline.Close
	}
	atrValue := Round(percentTotal/float64(klineNum)*100, 1)
	t.Config.DoubleThrowRatio = atrValue
	t.Config.ProfitRatio = atrValue
}

// Round 四舍五入，ROUND_HALF_UP 模式实现
// 返回将 val 根据指定精度 precision（十进制小数点后数字的数目）进行四舍五入的结果。precision 也可以是负数或零。
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}

// SetRecordPrice 记录交易价格
func (t *TradingPair) SetRecordPrice(successPrice float64) {
	t.RunBet.RecordedPrice = append(t.RunBet.RecordedPrice, successPrice)
}

// ModifyPrice 买入后，修改 补仓价格 和 网格平仓价格以及步数
func (t *TradingPair) ModifyPrice(successPrice, curMarketPrice float64, step int64) {
	t.RunBet.NextBuyPrice = Round(successPrice*(1-t.Config.DoubleThrowRatio/100), 6)
	t.RunBet.GridSellPrice = Round(successPrice*(1+t.Config.ProfitRatio/100), 6)
	if t.RunBet.NextBuyPrice > curMarketPrice {
		t.RunBet.NextBuyPrice = Round(curMarketPrice*(1-t.Config.DoubleThrowRatio/100), 6)
	} else if t.RunBet.GridSellPrice < curMarketPrice {
		t.RunBet.GridSellPrice = Round(curMarketPrice*(1+t.Config.ProfitRatio/100), 6)
	}
	t.RunBet.Step = step
	fmt.Printf("ModifyPrice: NextBuyPrice:%v \n", t.RunBet.NextBuyPrice)
	fmt.Printf("ModifyPrice: GridSellPrice:%v \n", t.RunBet.GridSellPrice)
}

// GetLastPrice
func (t *TradingPair) GetLastPrice() float64 {
	if len(t.RunBet.RecordedPrice) == 0 {
		return 0
	}
	lastPrice := t.RunBet.RecordedPrice[0]
	step := t.RunBet.Step - 1
	if step >= 0 && int(step) <= len(t.RunBet.RecordedPrice)-1 {
		lastPrice = t.RunBet.RecordedPrice[step]
	}
	return lastPrice
}

func (t *TradingPair) RemoveRecordPrice() {
	t.RunBet.RecordedPrice = t.RunBet.RecordedPrice[:len(t.RunBet.RecordedPrice)-1]
}
