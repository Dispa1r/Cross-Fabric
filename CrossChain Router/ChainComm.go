/*
Copyright xujf000@gmail.com .2020. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

var (
	ccFile       = "./config.yaml"
	userName     = "User1"
	orgName      = "Org1"
	channelName  = "mychannel"
	ccDataGenerator    = "dataGenerator"
	ccMathTest = "mathtest"
	ccRegistInfo = "chainInfo"
	ccCrossChainResult = "supervisor"
	ccCrossChainMsg = "crossChainMsg"
)

var sdk *fabsdk.FabricSDK
var cclient *channel.Client

func InitCCOnStart() error {
	sdk, err := fabsdk.New(config.FromFile(ccFile))
	if err != nil {
		log.Println("WARN: init Chaincode SDK error:", err.Error())
		return err
	}
	clientContext := sdk.ChannelContext(channelName, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	if clientContext == nil {
		log.Println("WARN: init Chaincode clientContext error:", err.Error())
		return err
	} else {
		cclient, err = channel.New(clientContext)
		if err != nil {
			log.Println("WARN: init Chaincode cclient error:", err.Error())
		}
	}
	log.Println("Chaincode client initialed successfully.")
	return nil
}

func GetChannelClient() *channel.Client {
	return cclient
}

func CCinvoke(channelClient *channel.Client, ccname, fcn string, args []string) ([]byte, error) {
	var tempArgs [][]byte
	for i := 0; i < len(args); i++ {
		tempArgs = append(tempArgs, []byte(args[i]))
	}
	qrequest := channel.Request{
		ChaincodeID:     ccname,
		Fcn:             fcn,
		Args:            tempArgs,
		TransientMap:    nil,
		InvocationChain: nil,
	}
	//log.Println("cc exec request:",qrequest.ChaincodeID,"\t",qrequest.Fcn,"\t",qrequest.Args)
	response, err := channelClient.Execute(qrequest)
	if err != nil {
		return nil, err
	}
	return response.Payload, nil
}

func CCquery(channelClient *channel.Client, ccname, fcn string, args []string) ([]byte, error) {
	var tempArgs [][]byte
	if args == nil {
		tempArgs = nil
	} else {
		for i := 0; i < len(args); i++ {
			tempArgs = append(tempArgs, []byte(args[i]))
		}
	}
	qrequest := channel.Request{
		ChaincodeID:     ccname,
		Fcn:             fcn,
		Args:            tempArgs,
		TransientMap:    nil,
		InvocationChain: nil,
	}
	response, err := channelClient.Query(qrequest)
	if err != nil {
		return nil, err
	}
	return response.Payload, nil
}

//
//func (that *RouterService) CheckLP(args Args3,reply *int) error{
//	m := make(map[int64]Args3)
//	timeStamp := time.Now().Unix()
//	m[timeStamp] = args
//	//fmt.Println(args.C, args.X, args.B, args.Y)
//	var sum1 float64 = 0
//	var sum2 float64 = 0
//	for i := range args.C{
//		tmp := args.C[i] * args.X[i]
//		sum1 += tmp
//		//fmt.Println(tmp)
//	}
//	for i := range args.B{
//		tmp := args.B[i] * args.Y[i]
//		sum2 += tmp
//		//fmt.Println(tmp)
//	}
//	fmt.Println(sum1,sum2)
//
//	sub := sum1 - sum2
//	if sub >= -1 && sub <= 1{
//		* reply = 1
//		fmt.Println("chain calc right")
//	}else {
//		*reply = 0
//		fmt.Println("chain calc wrong")
//	}
//	return nil
//}
//
//func (that *RouterService) GetData(reply int64,args *Args3) error{
//	timeStart := strconv.FormatInt(reply,10)
//	arg3 := GetCalcResult(timeStart)
//	args.Y = arg3.Y
//	args.X = arg3.X
//	args.C = arg3.C
//	args.B = arg3.B
//	return nil
//}
//
//func  (that *RouterService) ReceiveSupervise(args SupervisorMessage,reply *ReturnMessage) error{
//
//}

