package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

type RpcServer struct {
}

func timeCost() func() {
	start := time.Now()
	return func() {
		tc:=time.Since(start)
		log.Printf("time cost = %v\n", tc)
	}
}

func (t *RpcServer) SendCrossChainMsg(msg Message,key *string) error {
	log.Println("SendCrossChainMsg cost...")
	defer timeCost()()
	TargetAddress := CallGetAddressById(msg.TCID)
	log.Println(TargetAddress)
	log.Println("Get Target Address")
	// save the cross chain msg
	err := SetCrossChainMsg(msg)
	if err != nil {
		log.Println("fail to set Crosschain Msg")
		return err
	}
	err = TransferMsg(TargetAddress,msg)
	if err != nil {
		log.Println("fail to transfer the msg")
		return err
	}
	return nil
}

func (t *RpcServer) RegisterInfo(args Register, reply *int) error {
	log.Println("get register info cost...")
	defer timeCost()()
	log.Println(args)
	err := SetCrossChainInfo(args)
	if err != nil {
		*reply = 0
		return err
	}
	*reply = 1
	return nil
}

func (t *RpcServer) GetPubKeyById(chainId string,pubKey *[]byte) error {
	log.Println("get public key by id cost...")
	defer timeCost()()
	PubKey := CallGetPubkeyById(chainId)
	if len(PubKey) == 0{
		log.Println("fail to get pub key")
		return errors.New("fail to get pub key")
	}
	*pubKey = PubKey
	return nil
}

func (t *RpcServer) GetAllNormalChains(chainId string,chainlist *string) error {
	log.Println("get all normal chains cost...")
	defer timeCost()()
	chainList := CallGetAllNormalChain()
	*chainlist = chainList
	log.Println("chainList",chainList)
	return nil
}

func (t *RpcServer) GetCrossChainMsg(msg Message,code *int) error{
	if msg.Type == "to"{
		// 说明应该返回数据了
		log.Println("get cross chain msg to cost...")
		defer timeCost()()
		log.Println("Get Message From relay chain",msg)
		result := CheckSign(msg.SCID,msg)
		if !result {
			log.Println("Sign invalid, please check the msg")
			return errors.New("invalid msg")
		}
		log.Println("sign verify success")
		err := SendDataToRelayChain(msg)
		if err != nil {
			log.Println("fail to transfer msg")
			*code = 1
			return err
		}
		*code = 0
		return nil
	}else if msg.Type == "back" {
		// 说明返回的数据到了
		log.Println("get cross chain msg back cost...")
		defer timeCost()()
		result := CheckSign(msg.SCID,msg)
		if !result {
			log.Println("Sign invalid, please check the msg")
			return errors.New("invalid msg")
		}
		// call the contract to check the proof
		err := CCCheckLP(msg)
		if err != nil{
			log.Println("fail to write cross chain result")
			*code = 1
			return err
		}
		*code = 0
		return nil
	}
	return nil
}

func (t *RpcServer) GetChainPrivKey(privkey []byte,code *int) error{

	nowPri,err  := ReadPrivateKeyFile()
	if len(nowPri) != 0 {
		return errors.New("this chain have registed")
	}
	log.Println("get get Chain Private key cost...")
	defer timeCost()()
	if len(ChainPrivateKey) != 0 {
		log.Println("you have set the ChainPrivateKey")
		return errors.New("duplicate ChainPrivateKey")
	}
	log.Println(string(privkey))
	ChainPrivateKey = privkey
	err = SavePrivateKeyFile(ChainPrivateKey)
	if err != nil {
		log.Println(err)
	}
	log.Println("update chain privateKey",string(privkey))
	return nil
}

func (t *RpcServer) GetChainId(id string,code *int) error{
	log.Println("get get Chain id key cost...")
	defer timeCost()()
	if len(ChainId) != 0 {
		log.Println("you have set the chainID")
		return errors.New("duplicate chainID")
	}
	log.Println("success to regist info")
	ChainId = id
	UpdateConfig()
	log.Println("update chain Id",ChainId)
	return nil
}

func (t *RpcServer) StartRegist(id string,code *int) error{
	err := RegistChainTest()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (t *RpcServer) StartCrossChain(id string,code *int) error{
	err := SendCrossChainRequestTest(ChainId,id,"Lp")
	if err != nil {
		log.Println(err)
	}
	return err
}

func StartRPC() {

	arith := new(RpcServer)
	rpc.Register(arith)

	tcpAddr, err := net.ResolveTCPAddr("tcp", localPort)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}
		/*
			ServeConn在单个连接上执行DefaultServer。ServeConn会阻塞，服务该连接直到客户端挂起。调用者一般应另开线程调用本函数："go serveConn(conn)"。ServeConn在该连接使用JSON编解码格式。
		*/
		go jsonrpc.ServeConn(conn)
	}
}

func checkError(err error) {
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}