package wallet

import (
	"blockchain_go/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4
const seedLength = 64 //种子的位数

// A wallet construct by a pk sk pair
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() (*Wallet, []string) {
	//private, public := newKeyPair()
	private, public, mnemonicCode := newSeedAndKeyPair()
	wallet := Wallet{private, public}

	return &wallet, mnemonicCode
}

// GetAddress returns wallet address
// 地址由 版本号(2字节) + 公钥哈希(20字节 160位) + 校验和(4字节)
func (w Wallet) GetAddress() []byte {
	pubKeyHash := utils.HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)

	return address
}

// ValidateAddress check if address if valid
func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// checksum 对公钥取两次哈希得到4字节的校验和
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// 通过椭圆曲线生成一个新的公私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	fmt.Println("new Seed and key pair")
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// genSeedAndKeyPair 通过一个64位的种子作为源来生成公私钥对
func newSeedAndKeyPair() (ecdsa.PrivateKey, []byte, []string) {
	curve := elliptic.P256()
	randNum, err := rand.Prime(rand.Reader, seedLength)
	mnemonicCode := genMnemonicCodeBySeed(randNum)
	if err != nil {
		log.Panic(err)
	}
	// 由于io.reader长度必须大于36位，这里采用了将20位大整数种子重复一次得到40位的数
	private, err := ecdsa.GenerateKey(curve, strings.NewReader(randNum.String()+randNum.String()))
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey, mnemonicCode
}

// genKeyPairByMnemonicCode 通过助记词还原一个公私钥对
func genKeyPairByMnemonicCode(memCode []string) *Wallet {
	curve := elliptic.P256()
	seed := genSeedByMemeoryCode(memCode)
	private, err := ecdsa.GenerateKey(curve, strings.NewReader(seed.String()+seed.String()))
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	wallet := Wallet{*private, pubKey}
	return &wallet
}

// genSeedByMemeoryCode 根据助记词来生成种子
func genSeedByMemeoryCode(memCode []string) *big.Int {
	var memCodeArr = initMemCodeArray()
	var seedArray [9]byte
	for memIndex, memCodeStr := range memCode {
		for index, tmpStr := range memCodeArr {
			tmpStr = memCodeArr[index]
			if strings.Compare(tmpStr, string(memCodeStr)) == 0 {
				seedArray[memIndex] = byte(index)
				break
			}
		}
	}
	//通过验证和，检查助记词是否合法
	sha256H := sha256.New()
	sha256H.Reset()
	sha256H.Write(seedArray[0:8])
	seedHash := sha256H.Sum(nil)

	seed := new(big.Int)
	if bytes.Equal(seedHash[0:1], seedArray[8:9]) {
		seed.SetBytes(seedArray[0:8])
	} else {
		fmt.Println("Invalid mnemonic code!")
		seed = big.NewInt(0)
	}
	return seed
}

// genMnemonicCodeBySeed 根据种子产生助记词
func genMnemonicCodeBySeed(seed *big.Int) []string {
	//种子转换为字符数组
	seedArray := make([]byte, seedLength)
	seedArray = seed.Bytes()

	//对种子进行sha-256，取前8位作为校验码,得到64 + 8的数
	sha256H := sha256.New()
	sha256H.Reset()
	sha256H.Write(seedArray)
	seedHash := sha256H.Sum(nil)

	seedArrayWithCheckSum := append(seedArray, seedHash[0:1]...)

	//字符数组转换为9个助记词
	mnemonicCode := make([]string, 0)
	//fmt.Println("Please store your mnemonic code in a safe place ")
	var memCodeArr = initMemCodeArray()
	for index, value := range seedArrayWithCheckSum {
		value = seedArrayWithCheckSum[index]
		tmpStr := memCodeArr[value]
		mnemonicCode = append(mnemonicCode, tmpStr)
	}
	return mnemonicCode
}

// initMemCodeArray 初始化助记词字符串数组
func initMemCodeArray() []string {
	var memCodelist = []string{"abandon", "ability", "whisper", "about", "above", "absent", "absorb", "bag",
		"fade", "illegal", "barely", "deal", "hedgehog", "camera", "hello", "faint",
		"castle", "idea", "can", "laundry", "imitate", "basic", "heart", "gallery",
		"debate", "fall", "echo", "panther", "false", "debris", "ignore", "table",
		"ecology", "battle", "mammal", "casual", "sample", "beach", "camp", "leader",
		"talent", "decade", "original", "beauty", "kidney", "december", "parent", "marine",
		"bean", "health", "decide", "opinion", "magnet", "beauty", "hidden", "orchard",
		"jealous", "jaguar", "game", "kingdom", "manage", "party", "mansion", "window",
		"horror", "decline", "abstract", "patch", "beef", "either", "satisfy", "eight",
		"call", "machine", "initial", "magic", "patient", "gap", "lawsuit", "sausage",
		"faculty", "become", "pattern", "carbon", "veteran", "decorate", "calm", "wrestle",
		"canoe", "payment", "decrease", "target", "because", "unusual", "marble", "canal",
		"peanut", "defense", "scare", "garbage", "absurd", "car", "universe", "egg",
		"orange", "badge", "abuse", "option", "science", "famous", "scorpion", "timber",
		"cannon", "fantasy", "peasant", "scatter", "fame", "giggle", "unique", "version",
		"honey", "negative", "define", "tattoo", "before", "venture", "uniform", "wool",
		"farm", "neglect", "begin", "history", "innocent", "adapt", "hockey", "wrong",
		"capital", "inquiry", "addict", "injury", "defy", "behave", "hint", "holiday",
		"rack", "behind", "pencil", "captain", "canyon", "accident", "yellow", "young",
		"fashion", "gas", "canvas", "cat", "label", "deliver", "capable", "thought",
		"radio", "screen", "gasp", "fat", "teach", "believe", "gather", "cancel",
		"economy", "garment", "achieve", "kitchen", "across", "adult", "mechanic", "giant",
		"educate", "fatal", "raise", "knife", "galaxy", "favorite", "various", "victory",
		"nephew", "below", "kitten", "marriage", "hotel", "actor", "human", "utility",
		"hospital", "ladder", "office", "rally", "october", "season", "olympic", "learn",
		"edge", "cart", "language", "carpet", "elbow", "benefit", "zoo", "sphere",
		"secret", "actress", "demand", "material", "master", "upgrade", "spike", "zone",
		"security", "noodle", "notable", "belt", "random", "carry", "gauge", "leisure",
		"card", "nominee", "laptop", "oblige", "rapid", "obvious", "youth", "genuine",
		"gaze", "neutral", "bench", "cargo", "actual", "math", "genius", "zero",
		"network", "genre", "casino", "rare", "wheat", "case", "gentle", "sponsor",
		"fabric", "elegant", "rate", "meadow", "tenant", "electric", "measure", "tobacco"}
	return memCodelist
}
