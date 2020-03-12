package test

import (
	"fmt"
	"go-disk/common/jwt"
	"log"
	"testing"
)

func TestGenToken(t *testing.T) {
	tokenString, err := jwt.GenToken(map[string]interface{}{
		"username" : "yeonon",
	})

	if err != nil {
		log.Printf("%v", err)
		return
	}

	fmt.Println(tokenString)

}
