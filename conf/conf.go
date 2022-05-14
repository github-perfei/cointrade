package conf

import (
	"cointrade/internal/model"
	"encoding/json"
	"os/user"
	"path"

	"github.com/nntaoli-project/goex"
	"github.com/spf13/viper"
)

// 变量
const (
	EnvPrefix = "COIN_TRADE"
	FileName  = "data"
)

// InitViper 读取配置
func InitViper() error {
	var err error

	usr, err := user.Current()
	if err != nil {
		return err
	}

	viper.SetEnvPrefix(EnvPrefix)

	viper.SetConfigName(FileName)
	viper.AddConfigPath(".")
	viper.AddConfigPath(path.Join(usr.HomeDir, path.Join(".cointrade")))
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// GetCoinList
func GetCoinList() []string {
	key := "coinList"
	if !viper.IsSet(key) {
		panic("GetCoinList" + key + "not set")
	}
	return viper.GetStringSlice(key)
}

// GetTradePairData
func GetTradePairData(coinType string, api goex.API) *model.TradingPair {
	key := coinType
	if !viper.IsSet(key) {
		panic("GetTradePairData" + key + "not set")
	}
	var result model.TradingPair
	jsonStr, _ := json.Marshal(viper.GetViper().Get(key))
	err := json.Unmarshal(jsonStr, &result)
	if err != nil {
		panic(err)
	}
	result.Build(coinType, api)
	return &result
}

// SetTradePairData
func SetTradePairData(pair *model.TradingPair) error {
	key := pair.PairName
	result := make(map[string]interface{})
	if jsonStr, err := json.Marshal(pair); err != nil {
		return err
	} else if err := json.Unmarshal(jsonStr, &result); err != nil {
		return err
	}
	viper.GetViper().Set(key, result)
	viper.GetViper().WriteConfig()
	return nil
}
