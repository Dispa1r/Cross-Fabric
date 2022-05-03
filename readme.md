# 基于中继链的联盟链跨链监管框架
2022.5.3更新：添加密钥协商阶段以及通信过程消息加密

目前实现的监管计算：线性规划
## CrossFabric使用流程

### 环境准备
两台ubuntu虚拟机作为监管与被监管链，一台云服务器作为中继链，由于所有通信流程在公网进行，因此需要内网穿透工具，这里选择使用[natapp](https://natapp.cn/)

### 被监管链准备
使用到的脚本文件均在tools文件夹内
```
# 启动基础环境（byfn）,添加第三个组织
sudo ./1-2.startNetwork.sh 
sudo ./2.addOrg3.sh

# 安装和实例化代码
# 安装伪随机数据生成链码
sudo ./3-1.installDataGenerator.sh 
# 生成随机数据,每次生成一万条
sudo ./generateAndQuery.sh
# 安装线性规划链码
sudo ./3-1.installMathTest.sh 
```

接着需要修改一些ip的配置，修改的文件是```/appcode/fccserver/src/config/ChainInfo.yaml```，修改前先启动内网映射工具，默认映射到docker的8081端口，如需修改请同时修改```appcode```目录下的docker配置文件，主要需要配置的是公网ip与端口，以及跨链路由监听的端口，以及一些关于链的身份信息。

![](/uploads/upload_9b61196adf3a1fa2b1e47a9458ad953d.png)

配置结束后运行：

```
# 启动跨链路由
sudo ./4.startAppcli.sh 
# 查看日志
sudo docker logs -f appcli
```
如果一切正常你应该看到这行日志
![](/uploads/upload_25bffbc763c35b6fff6d6c402340fb9f.png)

### 中继链准备

与被监管链基本一致，唯一不同的是，中继链安装的链码不同

```
# 安装身份信息与跨链消息记录合约
sudo ./InstallChainInfo.sh 
sudo ./InstallCrossChainMessage.sh
```

中继链不需要配置Ip地址与端口，唯一需要配置的就是relaychainaddress，也就是自身地址,配置完成后运行appcli

### 监管链准备


与被监管链基本一致，唯一不同的是，同样也是安装的链码不同

```
sudo ./InstallSupervisor.sh 
```
配置Ip地址与端口，以及中继链地址后，运行appcli

### 注册流程
直接在the god项目中运行StartRegist函数即可

### 监管流程
直接在the god中运行StartSupervisor函数即可

## 跨链路由

### Crypto Utils
* RSA sign & Verify
* AES encrypt & decrypt
* Base58 Encoding & Decoding

### RPC Interfaces
 * SendCrossChainMsg(msg Message,key *string) error // 发起跨链请求
 * RegisterInfo(args Register, reply *int) error
 * GetPubKeyById(chainId string,pubKey *[]byte) error // 注册身份信息
 * GetAllNormalChains(chainId string,chainlist *string) error // 获取全部非监管链的区块链id，格式为 ["1","2"]
 * GetCrossChainMsg(msg Message,code *int) error //监听跨链请求，根据信息种类的不同执行不同的操作
 * StartCrossChain(id string,code *int) error // 远程开启跨链
 * StartRegist(id string,code *int) error // 远程开启注册
 * GetChainPrivKey(privkey []byte,code *int) error // 接受私钥
 * GetChainId(id string,code *int) error // 接受链id
### Contract Interfaces
* GetTimeStart() string // 获取开始进行业务计算的时间
* GetCalcResult(timeStart string)  (LpProof) // 获取计算结果
* CallDataGenerator() error // 调用产生数据
* CallGetPubkeyById(chainId string) []byte // 获取对应chainId的公钥，进行消息合法性的验证
* CallGetAddressById(chainId string) string // 获取对应chainId的通信地址
* CallGetAllNormalChain() string // 使用富查询获取全部非监管链ID
* SetCrossChainInfo(reg Register) error // 注册信息上链
* SetCrossChainMsg(msg Message) error // 跨链消息上链
* CCCheckLP(msg Message) error // 调用监管函数对计算结果进行验证

### 单元测试
* RegistChainTest() error // 测试注册链
* SendCrossChainRequestTest(SCID,TCID,CalcType string) error // 测试发起跨链请求
* SendChainIdAndPrivateKeyTest (address string,id string,privateKey []byte) error // 测试从中继链获取私钥和id
* SignAndVerifyTest() bool // 消息签名与验证测试

### 合约部分
* chainInfo // 记录链身份信息
* crossChainMsg // 记录跨链消息
* dataGenerator // 模拟产生业务数据
* mathtest // 模拟进行链上计算
* writeSupervisor // 校验计算结果并上链

合约接口略

## 跨链调用栈
### 业务逻辑：
* 目标链 -> 目标链.DataGenerate
* 目标链 -> 目标链.Lp

### 注册逻辑：
* 申请链 -> 中继链.RegistInfo （分配chainid，返回私钥）
* 中继链 -> 申请链.GetChainId （设置chainid）
* 中继链 -> 申请链.GetPrivateKey （设置私钥）

### 跨链逻辑：
* 申请链 -> 中继链.SendCrossChainMsg（转发消息给目标链）
* 中继链 -> 中继链.SetCrossChainMsg (跨链连消息上链)
* 中继链 -> 中继链.TransferMsg (根据消息确定传输地址)
* 中继链 -> 目标链.GetCrossChainMsg  （获取到跨链消息后，判断类别）
* 目标链 —> 中继链.GetPublicKeyById   (申请公钥验证消息)

目标链跨链路由获取目标链数据，构造消息
* 目标链 —> 中继链.SendCrossChainMsg  (转发消息，类别为back)
* 中继链 -> 中继链.SetCrossChainMsg
* 中继链 -> 申请链.GetCrossChainMsg  （获取到跨链消息后，判断类别）
* 申请链 -> 中继链.GetPublicKeyById   (申请公钥验证消息)
* 申请链 -> 申请链.SetCrossChainResult

## bench mark

注册链信息耗时（身份信息上链）
```
2022/03/03 09:17:06 get register info cost...

2022/03/03 07:48:33 time cost = 4.428189629s
2022/03/03 12:42:03 time cost = 4.290309764s
2022/03/03 09:09:45 time cost = 4.393647887s
```

获取公钥耗时(链上读取+编码转发)
```
2022/03/02 14:24:46 get public key by id cost...
2022/03/02 14:24:48 time cost = 2.077447014s

2022/03/02 14:25:58 get public key by id cost...
2022/03/02 14:26:00 time cost = 2.094422312s

2022/03/02 14:25:39 get public key by id cost...
2022/03/02 14:25:41 time cost = 2.120858527s

2022/03/02 14:25:52 get public key by id cost...
2022/03/02 14:25:54 time cost = 2.087429543s

```

转发跨链信息耗时（跨链消息编码上链 + 转发消息）
```
2022/03/03 12:56:04 SendCrossChainMsg cost...

2022/03/03 12:34:56 time cost = 8.679100992s
2022/03/03 09:20:42 time cost = 8.75336204s
2022/03/03 09:17:36 time cost = 8.720701945s
```

TPS待测

## Todo List
* 支持异构跨链，编写bcos链版本的跨链路由
* 跨链结点去中心化
* 增加SPV验证，通过跨链路由传入目标链区块头进行监管交易验证
* 通信过程中消息采用对称加密（完成）
* 改进签名算法，RSA的公钥和私钥传输占用了过大网络性能




