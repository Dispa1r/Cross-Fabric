package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

func RegistChainTest1() error {
	// tcp://server.natappfree.cc:38332
	conn, err := jsonrpc.Dial("tcp", "c6f8c208063e1586.natapp.cc:6324")
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	var id string = "1"
	// server.natappfree.cc:37168
	err = conn.Call("RpcServer.StartRegist", id, &code)
	if err != nil {
		log.Println("call RpcServer.StartRegist error:", err)
		return err
	}
	return nil
}

func RegistChainTest2() error {
	// tcp://server.natappfree.cc:38332
	conn, err := jsonrpc.Dial("tcp", "server.natappfree.cc:33012")
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	var id string = "1"
	// server.natappfree.cc:37168
	err = conn.Call("RpcServer.StartRegist", id, &code)
	if err != nil {
		log.Println("call RpcServer.StartRegist error:", err)
		return err
	}
	return nil
}

func StartSupervise() error {
	// tcp://server.natappfree.cc:38332
	conn, err := jsonrpc.Dial("tcp", "server.natappfree.cc:33012")
	if err != nil {
		log.Println("fail to connect to target address")
		return err
	}
	var code int
	var id string = "1"
	// server.natappfree.cc:37168
	err = conn.Call("RpcServer.StartCrossChain", id, &code)
	if err != nil {
		log.Println("call RpcServer.StartRegist error:", err)
		return err
	}
	return nil
}

func main() {
	err := StartSupervise()
	fmt.Println(err)
}
