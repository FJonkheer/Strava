package Helper

// MAtrikel-Nr 3736476, 8721083
import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string { //Die Berechnung eines MD5-Hashes
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
