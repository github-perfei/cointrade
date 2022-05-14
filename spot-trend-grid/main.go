package main

import (
	"cointrade/conf"
	"cointrade/internal/model"
	"fmt"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

func main() {
	if err := conf.InitViper(); err != nil {
		panic(fmt.Sprintf("Fatal error config file: <%v> ", err))
	}
	coinList := conf.GetCoinList()
	api := builder.DefaultAPIBuilder.APIKey("").
		APISecretkey("").
		ApiPassphrase("").Build(goex.COINEX) //创建现货api实例
	for _, coin := range coinList {
		if data := conf.GetTradePairData(coin, api); data != nil {
			spotTrend(data)
		}
	}
}

func spotTrend(tradePair *model.TradingPair) {
	for {
		cData, err := tradePair.CurrentData()
		if err != nil {
			continue
		}
		fmt.Printf("当前data:\n GridBuyPrice:%v \n GridSellPrice:%v\n CurMarketPrice:%v\n Step:%v\n Quantity:%v\n ", cData.GridBuyPrice, cData.GridSellPrice, cData.CurMarketPrice, cData.Step, cData.Quantity)

		if cData.GridBuyPrice >= cData.CurMarketPrice {
			// api.MarketBuy(strconv.FormatFloat(cData.Quantity, 'E', -1, 64), "", cData.Pair)
			fmt.Printf("币种:%v 当前市价:%v 买入:%v\n", tradePair.Pair.String(), cData.CurMarketPrice, cData.Quantity)
			// 挂单成功
			successPrice := cData.CurMarketPrice
			tradePair.SetRatio()
			tradePair.SetRecordPrice(successPrice)
			tradePair.ModifyPrice(successPrice, cData.CurMarketPrice, cData.Step+1)
			conf.SetTradePairData(tradePair)
			time.Sleep(time.Minute)
		} else if cData.GridSellPrice < cData.CurMarketPrice {
			if cData.Step == 0 {
				// step == 0 防止踏空，跟随价格上涨
				tradePair.ModifyPrice(cData.GridSellPrice, cData.CurMarketPrice, cData.Step)
				conf.SetTradePairData(tradePair)
			} else {
				lastPrice := tradePair.GetLastPrice()
				sellAmount := tradePair.GetQuantity(false)
				porUsdt := (cData.CurMarketPrice - lastPrice) * sellAmount
				fmt.Printf("币种:%v 当前市价:%v 卖出:%v usdt:%v\n", tradePair.Pair.String(), cData.CurMarketPrice, sellAmount, porUsdt)
				// res = msg.sell_market_msg(coinType, runbet.get_quantity(coinType,False),porfit_usdt)
				// if 'orderId' in res: #True 代表下单成功
				tradePair.SetRatio() //#启动动态改变比率
				tradePair.ModifyPrice(lastPrice, cData.CurMarketPrice, cData.Step-1)
				tradePair.RemoveRecordPrice()
				conf.SetTradePairData(tradePair)
				time.Sleep(time.Minute)
			}
		} else {
			fmt.Printf("币种:%v 当前市价:%v 未能满足交易,继续运行\n", tradePair.Pair.String(), cData.CurMarketPrice)
			time.Sleep(5 * time.Second)
		}
		fmt.Println("==========")
	}
}
