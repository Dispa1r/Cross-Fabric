package main

import (
	"encoding/json"
	"log"
	"net/rpc/jsonrpc"
	"time"
)

func TransferMsg (address string,msg Message) error {
	log.Println(msg,address)
	conn, err := jsonrpc.Dial("tcp", address)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	err = conn.Call("RpcServer.GetCrossChainMsg", msg, &code)
	if err != nil {
		log.Println("call MathService.GetCrossChainMsg error:", err)
		return err
	}
	return nil
}

func SendDataToRelayChain(msg Message) error {
	timeStart := GetTimeStart()
	result := GetCalcResult(timeStart)
	msgNew := Message{
		UUID:      msg.UUID,
		SCID:      msg.TCID,
		TCID:      msg.SCID,
		CalcType:  msg.CalcType,
		TimeStamp: time.Now().Unix(),
		Proof:    result,
		Type:      "back",
	}
	jsbBytes,err := json.Marshal(msgNew)
	if err != nil {
		log.Println("fail to marsha the data")
	}
	ChainPrivateKey,err = ReadPrivateKeyFile()
	if err != nil {
		log.Println(err)
	}
	sign := RsaSignWithSha256(jsbBytes,ChainPrivateKey)
	msgNew.Sign = Base58Encoding(sign)
	// 消息让中继链转发
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	err = conn.Call("RpcServer.SendCrossChainMsg", msgNew, &code)
	log.Println("new msg",msg)
	if err != nil {
		log.Println("call MathService.GetCrossChainMsg error:", err)
		return err
	}
	return nil
}




