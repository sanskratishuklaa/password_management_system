package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string)(string, error){
	salt:= make([]byte,16)
	_,err:=rand.Read(salt)
	if err!=nil{
		return "",fmt.Errorf("Failed to generate salt: %w", err)
	}

	hash:=argon2.IDKey(
		[]byte(password),
		salt, 
		1,
		64*1024,
		4,
		32,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)


	return fmt.Sprintf(
		"$argon2id$v=19$m=65536,t=1,p=4$%s$%s",
		saltEncoded,
		hashEncoded,
	), nil
}