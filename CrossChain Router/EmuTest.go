package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
	"net/rpc/jsonrpc"
	"time"
)

func SendCrossChainRequestTest(SCID,TCID,CalcType string) error{
	u,_ := uuid.NewRandom()
	uu := u.String()
	msg := Message{
		UUID:     uu,
		SCID:      SCID,
		TCID:      TCID,
		CalcType:  CalcType,
		TimeStamp: time.Now().Unix(),
		Proof:     LpProof{},
		Type:      "to",
	}
	jsbytes,err := json.Marshal(msg)
	if err != nil {
		log.Println("fail to marshal the data")
		return err
	}
	ChainPrivateKey,err = ReadPrivateKeyFile()
	if err != nil {
		log.Println(err)
	}
	signData := RsaSignWithSha256(jsbytes,ChainPrivateKey)
	msg.Sign = Base58Encoding(signData)
	log.Println(msg)
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code string
	err = conn.Call("RpcServer.SendCrossChainMsg", msg, &code)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func RegistChainTest() error{
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	// server.natappfree.cc:37168
	reg := Register{
		Name:      "Chain1",
		Identity:  viper.GetString("ChainInfo.type"),
		Address:    viper.GetString("ChainInfo.address"),
		Port:      viper.GetString("ChainInfo.port"),
		CalcType:  viper.GetString("ChainInfo.calcResource"),
	}
	err = conn.Call("RpcServer.RegisterInfo", reg, &code)
	if err != nil {
		log.Println("call MathService.GetCrossChainMsg error:", err)
		return err
	}
	return nil
}

func SendChainIdAndPrivateKeyTest (address string,id string,privateKey []byte) error {
	conn, err := jsonrpc.Dial("tcp", address)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	log.Println(address,id,privateKey)
	err = conn.Call("RpcServer.GetChainId", id, &code)
	if err != nil {
		log.Println("call MathService.GetChainId error:", err)
		return err
	}
	err = conn.Call("RpcServer.GetChainPrivKey", privateKey, &code)
	if err != nil {
		log.Println("call MathService.GetChainPrivKey error:", err)
		return err
	}
	return nil
}

func CheckSignTest(bytes []byte,pubKey []byte) bool {
	msg := Message{}
	err := json.Unmarshal(bytes,&msg)
	if err != nil {
		log.Println("fail to unmarshal the data")
		return false
	}
	signData := Base58Decoding(msg.Sign)
	newMsg := Message{
		UUID:     msg.UUID,
		SCID:      msg.SCID,
		TCID:      msg.TCID,
		CalcType:  msg.CalcType,
		TimeStamp: msg.TimeStamp,
		Proof:     msg.Proof,
		Type:      msg.Type,
	}
	jsbytes,_ := json.Marshal(newMsg)
	log.Println([]byte(signData))
	result := RsaVerySignWithSha256(jsbytes, []byte(signData),pubKey)
	return result
}

func SignAndVerifyTest() error{
	prvKey, pubKey := GenRsaKey()
	u,_ := uuid.NewRandom()
	uu := u.String()
	msg := Message{
		UUID:     uu,
		SCID:      "1",
		TCID:      "2",
		CalcType:  "mathtest/Lp",
		TimeStamp: time.Now().Unix(),
		Proof:     LpProof{},
		Type:      "to",
	}
	jsbytes,err := json.Marshal(msg)
	//log.Println(jsbytes)
	if err != nil {
		log.Println("fail to marshal the data")
		return err
	}
	signData := RsaSignWithSha256(jsbytes,prvKey)
	base58Sign := Base58Encoding(signData)
	msg.Sign = base58Sign
	log.Println(base58Sign)
	jsbytes,err = json.Marshal(msg)
	if err != nil {
		log.Println("fail to marshal the data")
	}
	result := CheckSignTest(jsbytes,pubKey)
	log.Println(result)
	return nil
}

func GetAllChainTest() error{
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code string
	err = conn.Call("RpcServer.GetAllNormalChains", "test", &code)
	if err != nil {
		log.Println("call MathService.GetChainId error:", err)
		return err
	}
	log.Println(code)
	return nil
}