package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 假设DP问题是10*5的矩阵，部署两个合约写入数据，每个合约写入25个数据
// 写入的数据格式：ID,num
// 使用第三方合约？或者IPFS存储DP问题的解和数据来源的关系
// 数据分布式存入

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


type NumberGenerator struct {

}

var rowNum = 5
var columnNum = 10

func (t *NumberGenerator) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	err := stub.PutState(args[0],[]byte(args[1]))

	//TestRpc()
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *NumberGenerator) get (stub shim.ChaincodeStubInterface, args [] string) peer.Response{

	value, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func (t * NumberGenerator) GenerateNumber(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	num,err := strconv.Atoi(args[0])
	if err!= nil{
		return shim.Error("internal error")
	}
	size,err  := stub.GetState("size")
	sizeNum,_ := strconv.Atoi(string(size))
	rand.Seed(time.Now().Unix())
	for i:=0;i < num;i++{
		var num int
		id := strconv.Itoa(sizeNum+i)
		if (sizeNum+i) % 65 >=0 && (sizeNum+i) % 64 <= 59{
			num = rand.Intn(100)
		}else {
			num = rand.Intn(100)-90
		}

		numStr := strconv.Itoa(num)
		err := stub.PutState(id, []byte(numStr))
		fmt.Println(id,num)
		if err != nil{
			return shim.Error("fail to save the data")
		}
	}
	sizeNum += num
	err3 := stub.PutState("size",[]byte(strconv.Itoa(sizeNum)))
	if err3 != nil{
		return shim.Error("fail to update the size")
	}
	return shim.Success([]byte("success to generate number"))
}

func (t * NumberGenerator) Init(stub shim.ChaincodeStubInterface) peer.Response{

	args:= stub.GetStringArgs()
	if len(args)!=2{
		shim.Error("invalid parameter numbers")
	}
	err := stub.PutState(args[0],[]byte(args[1]))

	if err != nil {
		shim.Error(err.Error())
	}
	value,err1 := stub.GetState("size")
	if err1 != nil {
		shim.Error("fail to init the size")
	}
	return shim.Success(value)
}

func (t *NumberGenerator) Invoke (stub shim.ChaincodeStubInterface) peer.Response{

	fn, args := stub.GetFunctionAndParameters()

	if fn =="set" {
		return t.set(stub, args)
	}else if fn == "get"{
		return t.get(stub , args)
	}else if fn == "GenerateNumber"{
		return t.GenerateNumber(stub,args)
	}

	return shim.Error("Invoke fn error")
}


func main(){
	err := shim.Start(new(NumberGenerator))
	if err != nil {
		fmt.Println("start error")
	}
}


