package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"gonum.org/v1/gonum/optimize/convex/lp"

	"log"
	"strconv"
	"time"
)

type Helloworld struct {

}

type Args struct {
	Arg1, Arg2 int
}

type Args2 struct {
	optMax float64
	optMin float64
}

type Args3 struct {
	A []float64
	C []float64
	X []float64
	B []float64
	Y []float64
}

type SupervisorData struct {

}

//type Args1 struct {
//	T mat.Dense
//	Y mat.VecDense
//	D mat.VecDense
//}

func CalcLp(stub shim.ChaincodeStubInterface){
	// A * x = b x >= 0
	// minimize c * x
	nowTime := time.Now().Unix()
	err3 := stub.PutState("timeStart", []byte(strconv.FormatInt(nowTime,10)))
	if err3 != nil {
		fmt.Println("fail to update the timeStamp")
	}
	nowlen,_ := stub.GetState("flag")
	flag,err := strconv.Atoi(string(nowlen))
	flagP := &flag
	fmt.Println("flag before: ",flag)
	//conn, err := jsonrpc.Dial("tcp", "101.133.135.126:8080")
	if err != nil {
		log.Fatal("fail to connect 101.133.135.126:8080")
	}
	//var reply int
	for i := 0;i<100;i++{
		A := GenerateDense1(int64(i),stub,flagP) // A m * n
		_,c := GetC(int64(i),stub,flagP)         // c 1 * n
		_,b := GetB(int64(i),stub,flagP)         // b m * 1
		fmt.Println(A,c,b)
		//
		newc,newA := TransferToStandardMin(c,A)
		//m,n := newA.Dims()
		//fmt.Println("A dim",m,n)
		//fmt.Println("min",newc,newA,b)
		opt, x, err := lp.Simplex(newc, &newA, b, 0, nil)

		B := MatrixToDense(A.T())
		newb,newa := TransferToStandardMax(b,B)
		//m,n = newa.Dims()
		//fmt.Println("newA dim",m,n)
		//fmt.Println("max",newb,newa,b)
		//fmt.Println(newb,newa,c)
		opt1, y, err1 := SimplexMax1(newb, &newa, c, 0, nil)
		if err != nil {
			fmt.Println(opt)
			
			continue
		}
		if err1 != nil {
			fmt.Println(opt1)
			continue
		}
		//if err2 != nil {
		//	//fmt.Println(err)
		//	continue
		//}
		//str1 := strconv.FormatFloat(opt,'f',20,32)
		//str2 := strconv.FormatFloat(opt1,'f',20,32)
        A_slice := DenseToSlice(newa)
		var args = Args3{A_slice,newc, x,newb ,y}
		//fmt.Println(args)
		//err = conn.Call("MathService.CheckLP", args, &reply)
		//if err != nil {
		//	log.Fatal("call MathService.CheckLP error:", err)
		//}
		//if reply == 0{
		//	log.Fatal("result check error!!!")
		//}
		jsBytes,err3 := json.Marshal(args)
		if err3 != nil{
			fmt.Println("fail to marshal the data")
		}
		err4 := stub.PutState(strconv.FormatInt(nowTime+int64(i),10),jsBytes)
		if err4 != nil {
			fmt.Println("fail to store the solve result")
		}
		fmt.Printf("opt: %v\n", opt)
		fmt.Printf("x: %v\n", x)
		fmt.Printf("Maxopt: %v\n", opt1)
		fmt.Printf("y: %v\n", y)
		//fmt.Printf("Maxopt: %v\n", opt2)
		//fmt.Printf("y: %v\n", z)
	}
	fmt.Println("flag after: ",flag)
	err2 := stub.PutState("flag", []byte(strconv.Itoa(flag)))
	if err2 != nil {
		fmt.Println("fail to update the flag")
	}

}



func (t * Helloworld) GetMinimize(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	CalcLp(stub)
	return shim.Success(nil)
}


func (t * Helloworld) Init(stub shim.ChaincodeStubInterface) peer.Response{

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

func (t *Helloworld) Invoke (stub shim.ChaincodeStubInterface) peer.Response{

	fn, args := stub.GetFunctionAndParameters()

	if fn =="set" {
		return t.set(stub, args)
	}else if fn == "get"{
		return t.get(stub , args)
	}else if fn == "Lp"{
		return t.GetMinimize(stub,args)
	}else if fn == "getChain"{
		return t.chaincodeGet(stub,args)
	}

	return shim.Error("Invoke fn error")
}

func (t *Helloworld) set(stub shim.ChaincodeStubInterface , args []string) peer.Response{
	err := stub.PutState(args[0],[]byte(args[1]))

	//TestRpc()
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *Helloworld) get (stub shim.ChaincodeStubInterface, args [] string) peer.Response{

	value, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(value)
}

func (t *Helloworld) chaincodeGet(stub shim.ChaincodeStubInterface, args [] string) peer.Response{
	//value, err := stub.GetState(args[0])
	result := GetAnotherChainData(args[0],stub)
	fmt.Println(result)
	return shim.Success([]byte("success to get data from another chaincode"))
}

func main(){
	err := shim.Start(new(Helloworld))
	if err != nil {
		fmt.Println("start error")
	}
}