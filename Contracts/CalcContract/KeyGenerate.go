package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"strconv"
	"time"
)
var RowNum = 5
var ColumnNum = 10
func GetKey()(mat.Matrix,[]float64){
	data := make([]float64, RowNum*RowNum)
	var j int
	rand.Seed(time.Now().Unix())
	for i := 0;i< len(data);i= RowNum*j+j {
		data[i] = float64(rand.Intn(10000)+1000)
		j++
	}
	//for j:=0;j<=5;j++{
	//	data[+j*7] = float64(rand.Intn(1000)+10)
	//}
	a := mat.NewDense(RowNum, RowNum, data)
	return a,data
}

func GetAnotherChainData(id string,stub shim.ChaincodeStubInterface) float64{
	response := stub.InvokeChaincode("dataGenerator", [][]byte{[]byte("get"), []byte(id)}, "mychannel")
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
		fmt.Printf(errStr)
		return 0
	}
	Aval, err := strconv.Atoi(string(response.Payload))
	if err != nil {
		errStr := fmt.Sprintf("Error retrieving state from ledger for queried chaincode: %s", err.Error())
		fmt.Printf(errStr)
		return 0
	}
	return float64(Aval)
}

func  GetRealNumber(tmp mat.VecDense) mat.VecDense{
	for i:= 0;i< RowNum;i++{
		tmpValue := tmp.At(i,0)
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", tmpValue), 64)
		tmp.SetVec(i,value)
	}
	return tmp
}

func GetR() (mat.Vector,[]float64){
	data := make([]float64, RowNum)
	rand.Seed(time.Now().Unix()+10000000000)
	for i:=0;i< len(data);i++{
		data[i] = float64(rand.Intn(10000)+1000)
	}
	r := mat.NewVecDense(RowNum,data)
	return r,data
}

func GetC(i int64,stub shim.ChaincodeStubInterface,flag *int) (mat.Vector,[]float64){
	// get C 1 * n
	data := make([]float64, ColumnNum)
	for i:=0;i< len(data);i++{
		num := GetAnotherChainData(strconv.Itoa(*flag),stub)
		data[i] = num
		//fmt.Println("get data: ",data[i])
		*flag++
		fmt.Println("flag : ",*flag)
	}

	r := mat.NewVecDense(ColumnNum,data)
	return r,data
}

func GetB(i int64,stub shim.ChaincodeStubInterface,flag *int) (mat.Vector,[]float64){
	// get B  m * 1
	data := make([]float64, RowNum)
	for i:=0;i< len(data);i++{
		num := GetAnotherChainData(strconv.Itoa(*flag),stub)
		data[i] = num
		//fmt.Println("get data: ",data[i])
		*flag++
		fmt.Println("flag : ",*flag)
	}

	r := mat.NewVecDense(RowNum,data)
	return r,data
}

func GetA()(mat.Matrix,[]float64){
	data := make([]float64, RowNum*RowNum)
	var j int
	rand.Seed(time.Now().Unix())
	for i := 0;i< len(data);i++{
		data[i] = float64(rand.Intn(10000)+1000)
		j++
	}
	a := mat.NewDense(RowNum, RowNum, data)
	return a,data
}

func DenseToSlice(dense mat.Dense) []float64{
	var result []float64
	m,n := dense.Dims()
	for i:=0;i<m;i++{
		for j:=0;j<n;j++{
			result = append(result,dense.At(i,j))
		}
	}
	return result
}

func MatrixToSlice(dense mat.Matrix) []float64{
	var result []float64
	m,n := dense.Dims()
	for i:=0;i<m;i++{
		for j:=0;j<n;j++{
			result = append(result,dense.At(i,j))
		}
	}
	return result
}

func VectorDenseToSlice(vectorDense mat.VecDense) []float64{
	var result []float64
	for i:=0;i< RowNum;i++{
		result = append(result,vectorDense.At(i,0))
	}
	return result
}

func GenerateDense(i int64,stub shim.ChaincodeStubInterface) mat.Dense{
	// get A m * n
	data := make([]float64, RowNum*ColumnNum)
	var j int
	rand.Seed(time.Now().Unix()+i)
	for i := 0;i< len(data);i++{
		data[i] = float64(rand.Intn(100))
		j++
	}
	A := mat.NewDense(RowNum, ColumnNum, data)
	return *A
}



func GenerateDense1(i int64,stub shim.ChaincodeStubInterface,flag *int) mat.Dense{
	// get A m * n
	data := make([]float64, RowNum*ColumnNum)
	//rand.Seed(time.Now().Unix()+i)
	for i := 0;i< len(data);i++{
		num := GetAnotherChainData(strconv.Itoa(*flag),stub)
		data[i] = num
		//fmt.Println("get data: ",data[i])
		*flag++
		fmt.Println("flag : ",*flag)
	}
	A := mat.NewDense(RowNum, ColumnNum, data)
	return *A
}
