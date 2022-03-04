# 基于中继链的联盟链跨链监管框架
目前实现的监管计算：线性规划

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
* 通信过程中消息采用对称加密
* 改进签名算法，RSA的公钥和私钥传输占用了过大网络性能




