package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"log"
	"net/rpc/jsonrpc"
	"strconv"

)

// i hate the fucking society
type IdentityInfo struct {

}

type IdInfo struct {
	Id int
	Name string
	Identity string
	Address string
	Port string
	PublicKey []byte
	CalcType string
}


func (t *IdentityInfo) GetAllNormalChain(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	queryString := fmt.Sprintf("{\"selector\":{\"identity\":\"Normal\"}}")
	qis, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("query error:" + err.Error())
	}
	defer qis.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for qis.HasNext() {
		queryResponse, err := qis.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

func (t *IdentityInfo) setChainInfo(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	//err := stub.PutState(args[0],[]byte(args[1]))
	// this function need 5 parameters, Name,identity,Address,port,CalcType
	chainNum,err := stub.GetState("chainNum")
	chainNum_str := string(chainNum)
	fmt.Println("chainNum old",chainNum_str)
	if err != nil {
		shim.Error("haven't init the contract")
	}
	chainNum_int,err := strconv.Atoi(string(chainNum))
	if err != nil {
		fmt.Println("fail to transfer the chainNum")
	}
	prvKey, pubKey := GenRsaKey()
	info := IdInfo{
		chainNum_int,
		 args[0],
		 args[1],
		 args[2],
		 args[3],
		 pubKey,
		 args[4],
	}
	jsBytes, err := json.Marshal(info)
	if err != nil {
		return shim.Error("marshal json error:" + err.Error())
	}
	err = stub.PutState(string(chainNum), jsBytes)
	if err != nil {
		return shim.Error("error on putstate:" + err.Error())
	}
	//err = stub.PutState("priKey"+string(chainNum),prvKey)
	chainNumNew := strconv.Itoa(chainNum_int+1)
	fmt.Println("chainNum new",chainNumNew)
	err = stub.PutState("chainNum", []byte(chainNumNew))
	fmt.Println("chainNum update finish")
	address := info.Address + ":" + info.Port
	err = SendChainIdAndPrivateKey(address,chainNum_str,prvKey)
	return shim.Success([]byte("success to regist the chain"))
}



func SendChainIdAndPrivateKey (address string,id string,privateKey []byte) error {
	conn, err := jsonrpc.Dial("tcp", address)
	if err != nil {
		fmt.Println("fail to connect to target address")
		return err
	}
	var code int
	log.Println(address,id,privateKey)
	err = conn.Call("RpcServer.GetChainId", id, &code)
	if err != nil {
		fmt.Println("call MathService.GetChainId error:", err)
		//return err
	}
	err = conn.Call("RpcServer.GetChainPrivKey", privateKey, &code)
	if err != nil {
		fmt.Println("call MathService.GetChainPrivKey error:", err)
		//return err
	}
	return nil
}


func (t *IdentityInfo) getPubKeyById (stub shim.ChaincodeStubInterface, args [] string) peer.Response{

	value, err := stub.GetState(args[0])
	info := IdInfo{}
	err = json.Unmarshal(value,&info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(info.PublicKey)
	return shim.Success(info.PublicKey)
}

func (t *IdentityInfo) GetAddressById (stub shim.ChaincodeStubInterface, args [] string)  peer.Response {
	value, err := stub.GetState(args[0])
	info := IdInfo{}
	err = json.Unmarshal(value,&info)
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Println(info.PublicKey)
	return shim.Success([]byte(info.Address + ":" + info.Port))
}

func (t * IdentityInfo) Init(stub shim.ChaincodeStubInterface) peer.Response{
	args:= stub.GetStringArgs()
	if len(args)!=0{
		shim.Error("invalid parameter numbers")
	}
	err := stub.PutState("chainNum",[]byte("0"))
	if err != nil {
		shim.Error(err.Error())
	}
	value,err1 := stub.GetState("chainNum")
	fmt.Println(value)
	if err1 != nil {
		shim.Error("fail to init the number of chain")
	}
	return shim.Success(value)
}

func (t *IdentityInfo) Invoke (stub shim.ChaincodeStubInterface) peer.Response{

	fn, args := stub.GetFunctionAndParameters()

	if fn =="set" {
		return t.setChainInfo(stub, args)
	}else if fn == "getPubKeyById"{
		return t.getPubKeyById(stub , args)
	}else if fn == "getAddressById"{
		return t.GetAddressById(stub , args)
	}else if fn == "getNormalChain"{
		return t.GetAllNormalChain(stub , args)
	}
	return shim.Error("Invoke fn error")
}
func main(){
	err := shim.Start(new(IdentityInfo))
	if err != nil {
		fmt.Println("start error")
	}
}


