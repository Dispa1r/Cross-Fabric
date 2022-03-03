sudo ./1-2.startNetwork.sh
echo "finish start byfn, now add Org3"
sudo ./2.addOrg3.sh
echo "finish to add Org3, now add chaincode dataGenerator"
sudo ./InstallChainInfo.sh
sudo ./InstallCrossChainMsg.sh
echo "finish to install Relay chain contract, test the function"
sudo ./invokeFn.sh
echo "start the cross-chain node"
sudo ./4.startAppcli.sh

