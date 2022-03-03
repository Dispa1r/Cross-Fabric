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


    hello := new(IdentityInfo)
    stub := shim.NewMockStub("hello", hello)

    checkInit(t, stub,[][]byte{})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("set"), []byte("str"), []byte("helloworld-1111")})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("MatrixCal"),})
    checkInvoke(t, stub, [][]byte{[]byte("set"),[]byte("Test1"),[]byte("Supervisor"),[]byte("127.0.0.1"),[]byte("3306"),[]byte("Lp")})
    checkInvoke(t, stub, [][]byte{[]byte("getPubKeyById"),[]byte("0")})
    //}
}
