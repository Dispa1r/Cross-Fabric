#!/bin/sh

CC_RUNTIME_LANGUAGE=golang
CC_SRC_PATH=github.com/chaincode/dataGenerator

CONFIG_ROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer
ORG1_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
ORG1_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
ORG2_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
ORG2_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
ORG3_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
ORG3_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt
ORDERER_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker exec \
  -e CORE_PEER_LOCALMSPID=Org2MSP \
  -e CORE_PEER_MSPCONFIGPATH=${ORG2_MSPCONFIGPATH} \
  cli \
  peer chaincode invoke \
    -o orderer.example.com:7050 \
    -C mychannel \
    -n mathtest \
    -c '{"Args":["Lp"]}' \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \










