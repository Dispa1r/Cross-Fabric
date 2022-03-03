package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)
type Supervisor struct {

}
type MathService struct {

}

type Proof struct {
	A []float64 `json:"A"`
	C []float64 `json:"C"`
	X []float64 `json:"X"`
	B []float64 `json:"B"`
	Y []float64 `json:"Y"`
}

type SuperViseData struct {
	timeStamp int64

}
var RowNum int = 5
var ColNum int = 10

func (t * Supervisor) Init(stub shim.ChaincodeStubInterface) peer.Response{
	args:= stub.GetStringArgs()
	if len(args)!=2{
		shim.Error("invalid parameter numbers")
	}
	err := stub.PutState(args[0],[]byte(args[1]))
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)

}

func (t *Supervisor) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	err := stub.PutState(args[0],[]byte(args[1]))

	//TestRpc()
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *Supervisor) get (stub shim.ChaincodeStubInterface, args [] string) peer.Response{

	value, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}
func (t *Supervisor) Invoke (stub shim.ChaincodeStubInterface) peer.Response{

	fn, args := stub.GetFunctionAndParameters()

	if fn =="set" {
		return t.set(stub, args)
	}else if fn == "get"{
		return t.get(stub , args)
	}else if fn == "WriteSupervise"{
		return t.WriteSupervise(stub,args)
	}

	return shim.Error("Invoke fn error")
}

func (t *Supervisor) WriteSupervise(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//conn, err := jsonrpc.Dial("tcp", "server.natappfree.cc:40967")
	//if err != nil {
	//	log.Fatal("fail to connect server.natappfree.cc:40967")
	//}
	//var argsTmp Args3
	//var reply int64
	reply, _ := strconv.ParseInt(args[0], 10, 64)
	argsTmp := Proof{}
	err := json.Unmarshal(Base58Decoding(args[1]),&argsTmp)
	if err != nil {
		fmt.Println("fail to unmarshal data")
		return shim.Error("fail to unmarshal data")
	}
	//fmt.Println(args)
	//err = conn.Call("MathService.GetData",reply, &argsTmp)
	//if err != nil {
	//	log.Fatal("call MathService.CheckLP error:", err)
	//}
	//fmt.Println(argsTmp,reply)
	// 获取账本数据

	var sum1 float64 = 0
	var sum2 float64 = 0
	for i := range argsTmp.C{
		tmp := argsTmp.C[i] * argsTmp.X[i]
		sum1 += tmp
		//fmt.Println(tmp)
	}
	for i := range argsTmp.B{
		tmp := argsTmp.B[i] * argsTmp.Y[i]
		sum2 += tmp
		//fmt.Println(tmp)
	}
	sub := sum1 - sum2
	if sub >= -1 && sub <= 1{
		err := stub.PutState(strconv.FormatInt(reply,10), []byte("chain calc right"))
		if err != nil{
			return shim.Error("fail to save the result")
		}
		fmt.Println("chain calc right")
	}else {
		err := stub.PutState(strconv.FormatInt(reply,10), []byte("chain calc wrong"))
		if err != nil{
			return shim.Error("fail to save the result")
		}
		fmt.Println("chain calc wrong")
	}
	return shim.Success([]byte("success to save the result"))
}



func main(){
	err := shim.Start(new(Supervisor))
	if err != nil {
		fmt.Println("start error")
	}
}