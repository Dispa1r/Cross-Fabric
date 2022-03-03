package main

import (
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "testing"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
    res := stub.MockInit("1", args)
    if res.Status != shim.OK {
        fmt.Println("Init failed", string(res.Message))
        t.FailNow()
    }
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string) {
    res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte(name)})
    if res.Status != shim.OK {
        fmt.Println("Query", name, "failed", string(res.Message))
        t.FailNow()
    }
    if res.Payload == nil {
        fmt.Println("Query", name, "failed to get value")
        t.FailNow()
    }

    fmt.Println("Query value", name, "was ", string(res.Payload))

}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
    res := stub.MockInvoke("1", args)
    if res.Status != shim.OK {
        fmt.Println("Invoke", args, "failed", string(res.Message))
        t.FailNow()
    }
}




func Test_Helloworld(t *testing.T) {
    // m,_ := GetKey()
    // n,_ := GetKey()
    // r,_ := GetR()
    // A,_ := GetA()
    // x,_ := GetR()
    // var b mat.VecDense
    // b.MulVec(A,x)
    // //fmt.Println("A: ",A)
    // fmt.Println(mat.Formatted(x))
    // var T mat.Dense
    // T.Mul(m,A)
    // T.Mul(&T,n)
    // //fmt.Println("T: ",T)
    // var c mat.VecDense
    // c.MulVec(A,r)
    // c.AddVec(&c, &b)
    // //fmt.Println("c: ",c)
    // var d mat.VecDense
    // d.MulVec(m,&c)
    ////fmt.Println("d: ",d)
    // var y mat.VecDense
    // y.SolveVec(&T,&d)
    // //fmt.Println("y: ",y)
    // conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8080")
    // var reply int
    // sliceT := DenseToSlice(T)
    // sliceY := VectorDenseToSlice(y)
    // sliceD := VectorDenseToSlice(d)
    // var args = Args1{sliceT, sliceY,sliceD}
    // //tmp,err1 := json.Marshal(args)
    // //if err1 != nil{
    // //    log.Fatal("fail to marshal the args", err)
    // //}
    //
    //
    // // 调用 Add() 方法
    // err = conn.Call("MathService.Check", args, &reply)
    // if err != nil {
    //     log.Fatal("call MathService.Add error:", err)
    // }
    // if reply == 0{
    //     log.Fatal("result check error!!!")
    // }
    //// fmt.Printf("MathService.Check %v * %v = %v True:%d", args.T, args.Y,args.D, reply)
    //
    //fmt.Println("now calculate real result X...")
    //
    // var X mat.VecDense
    // X.MulVec(n,&y)
    // X.SubVec(&X,r)
    // X = GetRealNumber(X)
    // fmt.Println(mat.Formatted(&X))
    //var test float64
    //test = 1.000000001
    //var modNum float64
    //modNum = 1
    //math.Mod(test,modNum)
    //value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", test), 64)
    //fmt.Println(value)

    //fmt.Println(b)
    //var Y mat.VecDense
    //Y.SolveVec(A, &b)
    //fmt.Println(Y)
    //
    //

    //fmt.Println(GetKey())
    //data := new(NumberGenerator)
    //stub := shim.NewMockStub("data",data)
    //checkInit(t, stub, [][]byte{[]byte("size"), []byte("0")})
    //checkInvoke(t, stub,[][]byte{[]byte("GenerateNumber"),[]byte("50")})
    //checkInvoke(t, stub,[][]byte{[]byte("GenerateNumber"),[]byte("50")})

    hello := new(Helloworld)
    stub := shim.NewMockStub("hello", hello)

    checkInit(t, stub, [][]byte{[]byte("flag"), []byte("0")})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("set"), []byte("str"), []byte("helloworld-1111")})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("MatrixCal"),})
    checkInvoke(t, stub, [][]byte{[]byte("Lp"),})
    //}
}
