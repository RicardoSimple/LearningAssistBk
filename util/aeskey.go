package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func encryptAES(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, []byte(key)[:aes.BlockSize])
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, _ := base64.StdEncoding.DecodeString(text)
	cfb := cipher.NewCFBDecrypter(block, []byte(key)[:aes.BlockSize])
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 密码对比
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// IsValidPhoneNumber 手机正则校验
func IsValidPhoneNumber(phoneNumber string) bool {
	const phoneRegex = `^1[3-9]\d{9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phoneNumber)
}

// IsValidEmail 邮箱正则校验
func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
