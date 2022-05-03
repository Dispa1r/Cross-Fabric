package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/big"
	"net/rpc/jsonrpc"
	"os"
)

var base58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encoding(strByte []byte) string { //Base58编码
	//1. 转换成ascii码对应的值
	//strByte := []byte(str)
	//log.Println(strByte) // 结果[70 97 110]
	//2. 转换十进制
	strTen := big.NewInt(0).SetBytes(strByte)
	//log.Println(strTen)  // 结果4612462
	//3. 取出余数
	var modSlice []byte
	for strTen.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0) //余数
		strTen58 := big.NewInt(58)
		strTen.DivMod(strTen, strTen58, mod)             //取余运算
		modSlice = append(modSlice, base58[mod.Int64()]) //存储余数,并将对应值放入其中
	}
	//  处理0就是1的情况 0使用字节'1'代替
	for _, elem := range strByte {
		if elem != 0 {
			break
		} else if elem == 0 {
			modSlice = append(modSlice, byte('1'))
		}
	}
	//log.Println(modSlice)   //结果 [12 7 37 23] 但是要进行反转，因为求余的时候是相反的。
	//log.Println(string(modSlice))  //结果D8eQ
	ReverseModSlice := ReverseByteArr(modSlice)
	//log.Println(ReverseModSlice)  //反转[81 101 56 68]
	//log.Println(string(ReverseModSlice))  //结果Qe8D
	return string(ReverseModSlice)
}
func ReverseByteArr(bytes []byte) []byte { //将字节的数组反转
	for i := 0; i < len(bytes)/2; i++ {
		bytes[i], bytes[len(bytes)-1-i] = bytes[len(bytes)-1-i], bytes[i] //前后交换
	}
	return bytes
}

//就是编码的逆过程
func Base58Decoding(str string) []byte { //Base58解码
	strByte := []byte(str)
	//log.Println(strByte)  //[81 101 56 68]
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte(base58, byteElem) //获取base58对应数组的下标
		ret.Mul(ret, big.NewInt(58))               //相乘回去
		ret.Add(ret, big.NewInt(int64(index)))     //相加
	}
	//log.Println(ret) 	// 拿到了十进制 4612462
	//log.Println(ret.Bytes())  //[70 97 110]
	//log.Println(string(ret.Bytes()))
	return ret.Bytes()
}

func GenRsaKey() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		log.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//验证
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

func GetRemotePublicKey(chainid string) []byte {
	conn, err := jsonrpc.Dial("tcp", RelayChainAddress)
	if err != nil {
		log.Println("fail to connect to target address")
		return nil
	}
	var pubKey []byte
	err = conn.Call("RpcServer.GetPubKeyById", chainid, &pubKey)
	if err != nil {
		log.Println("call MathService.GetCrossChainMsg error:", err)
		return nil
	}
	return pubKey
}

func RsaEncrypt() {

}

func CheckSign(chainId string, msg Message) bool {
	pubkey := GetRemotePublicKey(chainId)
	tmpMsg := Message{
		UUID:      msg.UUID,
		SCID:      msg.SCID,
		TCID:      msg.TCID,
		CalcType:  msg.CalcType,
		TimeStamp: msg.TimeStamp,
		Proof:     msg.Proof,
		Type:      msg.Type,
	}
	jsbytes, err := json.Marshal(tmpMsg)
	if err != nil {
		log.Println("fail to marshal the data")
	}

	signdata := Base58Decoding(msg.Sign)
	checkResult := RsaVerySignWithSha256(jsbytes, signdata, pubkey)
	return checkResult
}

func ReadPrivateKeyFile() ([]byte, error) {
	file, err := os.Open("./pri.key")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	fmt.Println(string(content))
	return content, nil
}

func SavePrivateKeyFile(privKey []byte) error {
	err := ioutil.WriteFile("./pri.key", privKey, 0644)
	if err != nil {
		return err
	}
	return nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func getKey() [16]byte {
	//方法一
	u, _ := uuid.NewRandom()
	str := u.String()
	has := md5.Sum([]byte(str))
	//md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	//md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return has
}

func DecryptMsg(msg string) Message {
	msgAESBytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		log.Println("fail to base decode msg ", err)
	}
	if len(LocalKey) == 0 {
		log.Println("invalid local key")
		return Message{}
	}
	msgBytes, err := AesDecrypt(msgAESBytes, LocalKey)
	msgObejct := Message{}
	err1 := json.Unmarshal(msgBytes, &msgObejct)
	if err1 != nil {
		log.Println("fail to unmarshal msg")
		return Message{}
	}
	return msgObejct
}

func EncMsg(msg Message, key []byte) string {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("fail to marsha the message ", err)
		return ""
	}
	encBytes, err1 := AesEncrypt(msgBytes, key)
	if err1 != nil {
		log.Println("fail to encrypt the message ", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(encBytes)

}
