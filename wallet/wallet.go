package wallet

import (
	"BCAddressCode/base58"
	"XianfengChain04/utils"
	"bytes"
	"github.com/bolt"
	"encoding/gob"
	"crypto/elliptic"
)

const KEYSTORE = "keystores"
const ADDANDPAIR = "addrs_keypairs"

/**
 * 定义wallet结构体，用于管理地址和对应的秘钥对信息
 */
type Wallet struct {
	Address map[string]*KeyPair
	Engine  *bolt.DB
}

//map: key    value
//    add     秘钥对(私钥、公钥）
func (wallet *Wallet) NewAddress() (string, error) {

	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}

	//3、对公钥进行sha256哈希
	pubHash := utils.Hash256(keyPair.Pub)
	//4、ripemd160计算
	ripemdPub := utils.HashRipemd160(pubHash)

	//5、添加版本号0x00
	versionPub := append([]byte{0x00}, ripemdPub...)

	//6、两次hash(双hash）
	firstHash := utils.Hash256(versionPub)
	secondHash := utils.Hash256(firstHash)

	//7、截取前4个字节作为地址校验位
	check := secondHash[:4]

	//8、拼接到versionPub后面
	originAddress := append(versionPub, check...)

	//9、base58编码
	address, err := base58.Encode(originAddress), nil
	if err != nil {
		return "", err
	}

	//把新生成的地址和对应的秘钥对存入到wallet的map结构中管理起来
	wallet.Address[address] = keyPair //仅仅是内存

	//把更新了地址信息和对应秘钥对的map结构中的数据持久化存到db文件中
	wallet.SaveAddrAndKeyPairs2DB()

	return address, nil
}

/**
 * 该函数用于检查地址是否合法，如果符合地址规范，返回true
 * 如果不符合地址规范，返回false
 */
func (wallet *Wallet) CheckAddress(addr string) bool {
	//1、使用base58对传入的地址进行解码
	reAddrBytes := base58.Decode(addr) // versionPubHash + check

	if len(reAddrBytes) < 4 {
		return false
	}

	//2、取出校验位
	reCheck := reAddrBytes[len(reAddrBytes)-4:]

	//3、截取得到versionPubHash
	reVersionPubHash := reAddrBytes[:len(reAddrBytes)-4]

	//4、对versionPub进行双hash
	reFirstHash := utils.Hash256(reVersionPubHash)
	reSecondHash := utils.Hash256(reFirstHash)

	//5、对双hash以后的内容进行前四个字节的截取
	check := reSecondHash[:4]

	return bytes.Compare(reCheck, check) == 0
}

/**
 * 该方法用于将内存中的map数据中的地址和秘钥对信息保存到持久化文件中
 */
func (wallet *Wallet) SaveAddrAndKeyPairs2DB() {
	var err error
	wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}
		}
		//桶keystores已经存在，可以向桶中存放map的数据了
		//map[key]keypair
		gob.Register(elliptic.P256())
		buff := new(bytes.Buffer)
		encoder := gob.NewEncoder(buff)
		err = encoder.Encode(wallet.Address)
		if err != nil {
			return err
		}
		bucket.Put([]byte(ADDANDPAIR), buff.Bytes())
		return nil
	})
}

/**
 * 从文件中读取已经存在的地址和对应的秘钥对信息
 */
func LoadAddrAndKeyPairsFromDB(engine *bolt.DB) (*Wallet, error) {
	address := make(map[string]*KeyPair)
	var err error
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			return nil
		}
		//如果有keystore存在，从keystore桶中读取
		addsAndKeyPairsBytes := bucket.Get([]byte(ADDANDPAIR))
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(addsAndKeyPairsBytes))
		err = decoder.Decode(&address)
		return err
	})
	if err != nil {
		return nil, err
	}
	walet := &Wallet{
		Address: address,
		Engine:  engine,
	}
	return walet, err
}
