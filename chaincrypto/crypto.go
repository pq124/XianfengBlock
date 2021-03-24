package chaincrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 * 使用密码学随机生成私钥：椭圆曲线数字签名算法ECDSA
 * ECDSA：elliptic curve digital signature algorithm
 * ECC：elliptic curve crypto
 */
func NewPriKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand.Reader)
}

/**
 * 根据私钥获得公钥
 */
func GetPub(curve elliptic.Curve, pri *ecdsa.PrivateKey) []byte {
	return elliptic.Marshal(curve, pri.X, pri.Y)
}


