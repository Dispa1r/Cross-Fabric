package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/rpc/jsonrpc"
	"time"
)

func TransferMsg(address string, msg Message, cipher EncMsgStruct) error {
	log.Println(msg, address)
	conn, err := jsonrpc.Dial("tcp", address)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	if _, ok := Keys[msg.UUID]; !ok {
		return errors.New("invalid uuid")
	}
	var code int
	err = conn.Call("RpcServer.GetCrossChainMsg", cipher, &code)
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
		Proof:     result,
		Type:      "back",
	}
	jsbBytes, err := json.Marshal(msgNew)
	if err != nil {
		log.Println("fail to marsha the data")
	}
	ChainPrivateKey, err = ReadPrivateKeyFile()
	if err != nil {
		log.Println(err)
	}
	sign := RsaSignWithSha256(jsbBytes, ChainPrivateKey)
	msgNew.Sign = Base58Encoding(sign)
	final := EncMsg(msg, LocalKey)
	encMsg := EncMsgStruct{UUID: msg.UUID, Cipher: final}
	// 消息让中继链转发
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code string
	err = conn.Call("RpcServer.SendCrossChainMsg", encMsg, &code)
	log.Println("new msg", msg)
	if err != nil {
		log.Println("call MathService.GetCrossChainMsg error:", err)
		return err
	}
	return nil
}
