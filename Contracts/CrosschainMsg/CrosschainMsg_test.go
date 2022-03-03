package main

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "strconv"
    "testing"
    "time"
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


    hello := new(CrossChainMsgContract)
    stub := shim.NewMockStub("hello", hello)

    checkInit(t, stub,[][]byte{})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("set"), []byte("str"), []byte("helloworld-1111")})
    //    //checkQuery(t, stub, "str")
    //    //checkInvoke(t, stub, [][]byte{[]byte("MatrixCal"),})
    uuid := uuid.New()
    key := uuid.String()
    timeStamp := time.Now().Unix()
    timeStamp_str := strconv.FormatInt(timeStamp,10)

    checkInvoke(t, stub, [][]byte{[]byte("set"),[]byte(key),[]byte("1"),[]byte("1"),[]byte("Lp"),[]byte(timeStamp_str),[]byte("i am a base58")})
    checkInvoke(t, stub, [][]byte{[]byte("get"),[]byte("0")})
    //}
}
