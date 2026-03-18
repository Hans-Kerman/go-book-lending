package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

var coverPicDir = "/static/covers"
var defaultPicPath = coverPicDir + "/default.jpg"

func InitStaticDir() error {
	if err := os.MkdirAll("."+coverPicDir, 0755); err != nil {
		return err
	}
	_, err := os.Stat("." + defaultPicPath)
	if err != nil {
		return err
	}
	return nil
}

func ParsePicName(picBytes []byte) string {
	if len(picBytes) == 0 {
		return defaultPicPath
	}

	hash := sha256.Sum256(picBytes)
	hexStr := hex.EncodeToString(hash[:])

	// 注意：这里需要加上 "." 前缀，以相对路径的形式存入项目根目录
	err := os.WriteFile(filepath.Join("."+coverPicDir, hexStr+".jpg"), picBytes, 0644)
	if err != nil {
		return defaultPicPath
	}

	return filepath.Join(coverPicDir, hexStr+".jpg")
}
