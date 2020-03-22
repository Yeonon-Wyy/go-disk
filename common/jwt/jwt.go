package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-disk/common/log4disk"
)

const (

	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3Rr47NWSHIzP5Rm5AZRA
JPSGoX0wG34IqhEfh/92ztEM4ltIAJQq0QM2CvaviqXxLtjhqqOt2c5PoT797XN8
qsi+HAMaKPTLjRjll/FNwsxkoxhup/kwfxrxZcWzCChwQxTN84a5ZzjxIDTUCPsc
OxhMfftyJ2V/sVPqaIqx/rc2bu16vebhhiGrxOCR1V7htx9o9HgmloSc2Ebudoqz
fpcfDgQZ2mJ5EE0tMvwynFFHvHP2kcZwhGwq3GZHc/aIz+Y90oigG4FHzn/I1vvm
diT9kVu7JWTs0eEcx3IT2yFIGuXU57HJyF7TKe/d+q2yqVoRclLuA3ZCgWWt0xw3
BwIDAQAB
-----END PUBLIC KEY-----
`
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA3Rr47NWSHIzP5Rm5AZRAJPSGoX0wG34IqhEfh/92ztEM4ltI
AJQq0QM2CvaviqXxLtjhqqOt2c5PoT797XN8qsi+HAMaKPTLjRjll/FNwsxkoxhu
p/kwfxrxZcWzCChwQxTN84a5ZzjxIDTUCPscOxhMfftyJ2V/sVPqaIqx/rc2bu16
vebhhiGrxOCR1V7htx9o9HgmloSc2EbudoqzfpcfDgQZ2mJ5EE0tMvwynFFHvHP2
kcZwhGwq3GZHc/aIz+Y90oigG4FHzn/I1vvmdiT9kVu7JWTs0eEcx3IT2yFIGuXU
57HJyF7TKe/d+q2yqVoRclLuA3ZCgWWt0xw3BwIDAQABAoIBAGpodKvmDK9YxSSI
wJSV+FjQpYpKaUCR4zGVlAsrUs4tpXm6XGiK5iA432VfWxPq0KuvDMvGggB0XbZI
ToRcM/8tJPDuPUTAqsV42eXJ55Z8L2Kee4KzVjeVi99iycp/S6e893DfwZJ/wOuz
AOhhkTCPfSCURlfXbSC2NfWh6g+ez1RcUIvgOw4z4DHnCpPvL70JpO0PFJBtVROp
McUsngRGYG5xv6dDFVSt1GAQLgiWk+GINm/ellpQvsKwx2jRMh75s2PcgCgx7mr/
6Juxlcw3oG3RApfUPQAcdkEZZbaW/QOcFrPQyv52ciJIHrVCSmZEp/HN4hko/62k
WYOS2KECgYEA/MgWxCqoOFSQB82vZNdJ6zOK46e/S/va8nsB6h59KsNLEoEEkimu
fRhLjyNTQMN8qDQa+ot9erdkupNLM8JiSgPKVBxWQeRFQlGlUdX7CxFiO8nmDWLy
y1xrELzUGnK5RHfJgztoP0Z0kskzuBawih7ZXMKKuvB+jonu7wLrWVUCgYEA3+uj
eS3YTVNXzD4KaBcZci2IO3sd8SwpxFRjQ6rS10LoRemsRLQGYURPVYwynndcTDje
SwPwZ8jAMqsYEqybVOfemLogSK9HhqH59fO057+wPuTlMtVyxfb9zbml35B/Oyqm
cjSWLBFfrcCXQf7c/JOZ4cufpj3EG0lo20E0XusCgYAHypMJENeGhPS7iNdzID+j
BD+vrKf0y5qABtKUSMRK4SIbO+bMKoS6Tlll3Azg3iWleZWrS0le0vBD3+5ddgxZ
g6xk00rFVSfdV27lCtdmC+8fMKXqm7YoFn0mUuumtQqI1bhcVyRrbtyA+bqiXfCr
ETBZ75UfFfKQqie3Ljva0QKBgQCvF6T2ZqDSpi7rKEAe8KEXJP9382eQZEYsnQgZ
q4O+izTxJi1sc5DhkfavIDecrhzgBT/dTLE5lkKj3CGwyIOVutHWfwQrkdPONO4u
Imj9Jmj8ZSPLwhhDMEV6Dobj8Cts6obImtIql1NHnGcqVc4bOpeDdiPabEXiSF7T
w4LeDQKBgQCGubv3GIjqJJG/RKjXQ6Lrji1PXXGJ+kUdLruP0AckRU31cSOYrIwf
8B36+zvomh5eUvlcNQsaCGHqWxomE94BvHWiLjn4q08oriI/7y1c2y9y6GsvUOHe
kW8eUWbLK+vHnehiCR+aeYqnz3udWICdmz3qpTK50+p+jt6UXOcifw==
-----END RSA PRIVATE KEY-----`
)

//internal function


func GenToken(payload map[string]interface{}) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(payload))
	securityKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}
	tokenString, err := token.SignedString(securityKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func GetPayload(tokenString string) (map[string]interface{}, bool) {
	token, ok := GetToken(tokenString)
	if !ok {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	}
	return nil, false
}

func GetToken(tokenString string) (*jwt.Token, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	})

	if err != nil {
		log4disk.E("parse token error : %v", err)
		return nil, false
	}

	return token, token.Valid
}




