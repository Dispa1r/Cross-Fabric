package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"time"

	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.AddConfigPath(workDir + "/config/")
	viper.SetConfigName("ChainInfo")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("read config failed", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	ChainPort = viper.GetString("ChainInfo.port")
	RelayChainAddress = viper.GetString("ChainInfo.relayChainAddress")
	ChainAddress = viper.GetString("ChainInfo.address")
	ChainType = viper.GetString("ChainInfo.type")
	ChainCalcResoure = viper.GetString("ChainInfo.calcResource")
	ChainId = viper.GetString("ChainInfo.id")
	localPort = viper.GetString("ChainInfo.localPort")
}

func UpdateConfig() {
	viper.Set("ChainInfo.port", ChainPort)
	viper.Set("ChainInfo.Address", ChainAddress)
	viper.Set("ChainInfo.calcResource", ChainCalcResoure)
	viper.Set("ChainInfo.relayChainAddress", RelayChainAddress)
	viper.Set("ChainInfo.id", ChainId)
	viper.Set("ChainInfo.type", ChainType)
	viper.Set("ChainInfo.localPort", localPort)
	viper.WriteConfig()
}

func main() {
	// 3.
	InitConfig()
	err := InitCCOnStart()
	go StartRPC()
	time.Sleep(1000)
	//GetAllChainTest()
	//err := SendCrossChainRequestTest(ChainId,"2","Lp")
	//err := RegistChainTest()
	if err != nil {
		fmt.Println("fail to connect the chain")
	}
	//err := TestSignAndVerify()
	select {}

}
