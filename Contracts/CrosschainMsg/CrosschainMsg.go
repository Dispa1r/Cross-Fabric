package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type CrossChainMsgContract struct {
}

type Message struct {
	UUID string
	SCID string
	TCID string
	CalcType string
	TimeStamp int64
	Sign string
	Proof LpProof
	Type string
}

type LpProof struct {
	C []float64 `json:"C"`
	X []float64 `json:"X"`
	B []float64 `json:"B"`
	Y []float64 `json:"Y"`
	A []float64 `json:"A"`
}

func (t *CrossChainMsgContract) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	// this function need seven parameters, uuid, SCID, TCID, CalcType, TimeStamp, Sign
	msgNum,err := stub.GetState("MsgNum")
	if err != nil {
		shim.Error("haven't init the contract")
	}
	msgNum_int,_ := strconv.Atoi(string(msgNum))

	i64, err := strconv.ParseInt(args[4], 10, 64)
	jsbytes := Base58Decoding(args[7])
	proof := LpProof{}
	err = json.Unmarshal(jsbytes,&proof)
	if err != nil {
		fmt.Println("fail to unmarshal the data")
		//return shim.Error("fail to unmarshal data")
	}
	msg := Message{
		UUID:      args[0],
		SCID:      args[1],
		TCID:      args[2],
		CalcType:  args[3],
		TimeStamp: i64,
		Sign:      args[5],
		Type:      args[6],
		Proof:     proof,
	}
	jsBytes, err := json.Marshal(msg)
	if err != nil {
		return shim.Error("marshal json error:" + err.Error())
	}
	err = stub.PutState(string(msgNum), jsBytes)
	if err != nil {
		return shim.Error("error on putstate:" + err.Error())
	}
	chainNumNew := strconv.Itoa(msgNum_int+1)
	err = stub.PutState("chainNum", []byte(chainNumNew))
	fmt.Println("msg number update finish")
	return shim.Success(nil)
}

func (t *CrossChainMsgContract) get (stub shim.ChaincodeStubInterface, args [] string) peer.Response{

	value, err := stub.GetState(args[0])
	msg := Message{}
	err = json.Unmarshal(value,&msg)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(msg.UUID)
	return shim.Success(value)
}


func (t * CrossChainMsgContract) Init(stub shim.ChaincodeStubInterface) peer.Response{
	args:= stub.GetStringArgs()
	if len(args)!=0{
		shim.Error("invalid parameter numbers")
	}
	err := stub.PutState("MsgNum",[]byte("0"))
	if err != nil {
		shim.Error(err.Error())
	}
	value,err1 := stub.GetState("MsgNum")
	fmt.Println(value)
	if err1 != nil {
		shim.Error("fail to init the number of msg")
	}
	return shim.Success(value)
}


func (t *CrossChainMsgContract) Invoke (stub shim.ChaincodeStubInterface) peer.Response{
	fn, args := stub.GetFunctionAndParameters()
	if fn =="set" {
		return t.set(stub, args)
	}else if fn == "get"{
		return t.get(stub , args)
	}
	return shim.Error("Invoke fn error")
}
func main(){
	err := shim.Start(new(CrossChainMsgContract))
	if err != nil {
		fmt.Println("start error")
	}
}


