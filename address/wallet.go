package address

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	_"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

//创建并返回一个wallet
func NewWallet() *Wallet{
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte){
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		fmt.Println(err)
		panic("Error##newKeyPair")
	}
	//在基于椭圆曲线的算法中，公钥是曲线上的点。
	//因此，公钥是 X、Y 坐标的组合。在比特币中，这些坐标连接起来形成一个公钥。
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

//返回wallet的地址
func (w Wallet) GetAddress() []byte{
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload,checksum...)
	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte{
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := crypto.RIPEMD160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil{
		panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

//给public key生成一个checksum
func checksum(payload []byte) []byte{
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

//检查地址是否有效
func ValidateAddress (address string) bool{
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1:len(pubKeyHash) - addressChecksumLen]
	targetChecksum := checksum(append([]byte{version},pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}


