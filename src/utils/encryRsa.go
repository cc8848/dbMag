package utils
import (
"fmt"
"crypto/rsa"
"crypto/rand"
"crypto/x509"
"encoding/pem"
"os"
"io/ioutil"
"errors"
)

func TestEncry() {

	var data []byte
	var err error
	var privatekey,publickey []byte

	//读取公钥文件
	publickey,err=ioutil.ReadFile("public.pem")
	if err !=nil{
		os.Exit(-1)
	}
	privatekey,err=ioutil.ReadFile("private.pem")

	if err !=nil{
		os.Exit(-1)
	}

	//加密
	data,err=RsaEncrypt(publickey,[]byte("hello world22222222222"))
	if err !=nil{
		panic(err)
	}

	//解密
	origData,err:=RsaDecrypt(privatekey,data)
	if err !=nil{
		panic(err)
	}
	fmt.Println(string(origData))
	//GenRsaKey(1024)
	fmt.Println("hello world")
}

/**
加密
@parm publickey : 加密使用的公钥文件
@parm:origData : 加密的信息
 */
func RsaEncrypt(publickey []byte,origData []byte) ([]byte,error) {
	block,_:=pem.Decode(publickey)
	if block==nil{
		return nil,errors.New("public key error")
		//return nil,errors.new("public key error")
	}
	pubInterface,err:=x509.ParsePKIXPublicKey(block.Bytes)
	if err!=nil{
		return nil,err
	}

	pub:=pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader,pub,origData)
}


/*
解密
@parm publickey : 解密使用的公钥文件
@parm:origData : 解密的信息
*/

func RsaDecrypt(privatekey []byte,ciphertext []byte) ([]byte,error) {
	block,_:=pem.Decode(privatekey)
	if block==nil{
		return nil,errors.New("private key error!")
	}
	priv,err:=x509.ParsePKCS1PrivateKey(block.Bytes)
	if err!=nil{
		return nil,err
	}
	return rsa.DecryptPKCS1v15(rand.Reader,priv,ciphertext)
}





/**
生成公钥文件
@parm: bits 密钥长度
 */
func GenRsaKey(bits int)  error{
	//生成私钥
	privatekey,err:=rsa.GenerateKey(rand.Reader,bits)
	if err!=nil{
		return err
	}

	derStream:=x509.MarshalPKCS1PrivateKey(privatekey)
	block:=&pem.Block{
		Type:"私钥",
		Bytes:derStream,
	}

	file,err:=os.Create("private.pem")
	if err!=nil{
		return err
	}

	err=pem.Encode(file,block)
	if err!=nil{
		return err
	}

	//生成公钥文件
	publickey:=&privatekey.PublicKey
	derPkix,err:=x509.MarshalPKIXPublicKey(publickey)
	if err!=nil{
		return err
	}
	block2:=&pem.Block{
		Type:"公钥",
		Bytes:derPkix,
	}
	file2,err:=os.Create("public.pem")
	if err!=nil{
		return err
	}

	err=pem.Encode(file2,block2)
	if err!=nil{
		return err
	}

	return nil
}