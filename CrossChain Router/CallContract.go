package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
)


func GetTimeStart() string {
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return "err"
	}
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccMathTest, "get",
		[]string{"timeStart"})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(bs))
	}
	return string(bs)
}

func GetCalcResult(timeStart string)  LpProof{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return LpProof{}
	}
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccMathTest, "get",
		[]string{timeStart})
	if err != nil {
		log.Println(err)
	}
	arg := LpProof{}
	err = json.Unmarshal(bs, &arg)
	if err != nil {
		log.Println("fail to Unmarshal the data")
	} else {
		log.Println(arg)
	}
	return arg
}


func CallDataGenerator() error{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return errors.New("fail to connect to the chain")
	}
	//check if bookid exist
	//new
	_, err := CCinvoke(cclient, ccDataGenerator, "GenerateNumber",
		[]string{"10000"})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func CallGetPubkeyById(chainId string) []byte{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return nil
	}
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccRegistInfo, "getPubKeyById",
		[]string{chainId})
	if err != nil {
		log.Println(err)
	}
	return bs
}

func CallGetAddressById(chainId string) string{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return ""
	}
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccRegistInfo, "getAddressById",
		[]string{chainId})
	if err != nil {
		log.Println(err)
	}
	return string(bs)
}

func CallGetAllNormalChain() string{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return ""
	}
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccRegistInfo, "getNormalChain",
		[]string{})
	if err != nil {
		log.Println(err)
	}
	return string(bs)
}


func SetCrossChainInfo(reg Register) error{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return errors.New("fail to connect to the chain")
	}
	//check if bookid exist
	//check if bookid exist
	//new
	bs, err := CCinvoke(cclient, ccRegistInfo, "set",
		[]string{reg.Name,reg.Identity,reg.Address,reg.Port,reg.CalcType})
	log.Println(bs)
	if err != nil {
		log.Println(err)
	}
	return nil
}


func SetCrossChainMsg(msg Message) error{
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return errors.New("fail to connect to the chain")
	}
	//check if bookid exist
	timestamp := strconv.FormatInt(msg.TimeStamp,10)
	jsbytes,err := json.Marshal(msg)
	if err != nil {
		log.Println("fail to marashal the data")
		return err
	}
	_, err = CCinvoke(cclient, ccCrossChainMsg, "set",
		[]string{msg.UUID,msg.SCID,msg.TCID,msg.CalcType,timestamp,msg.Sign,msg.Type,string(jsbytes)})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func CCCheckLP(msg Message) error {
	cclient = GetChannelClient()
	if cclient == nil {
		log.Println("fail to connect the chain")
		return errors.New("fail to connect to the chain")
	}
	//check if bookid exist
	timestamp := strconv.FormatInt(msg.TimeStamp,10)
	jsbytes,err := json.Marshal(msg.Proof)
	bs, err := CCinvoke(cclient, ccCrossChainResult, "WriteSupervise",
		[]string{timestamp,Base58Encoding(jsbytes)})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(bs))
	return nil
}


