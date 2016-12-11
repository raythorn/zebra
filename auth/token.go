package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	Base64Pattern = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+_"
)

//Token type
//
//	claims["iss"]	Issuer - string, the issuer of this token
//	claims["uid"]	User id - string, user id the token issued for
//	claims["exp"]	Expire - time.Time, expire time of this token
//	claims["nbf"]	Not before - time.Time, time that the token becomes active
type Token struct {
	claims map[string]interface{}
}

func NewToken() *Token {
	return &Token{make(map[string]interface{})}
}

func (t *Token) Set(key string, value interface{}) {
	t.claims[key] = value
}

func (t *Token) Get(key string) interface{} {
	if value, ok := t.claims[key]; ok {
		return value
	}
	return nil
}

func (t *Token) Sign() (string, error) {

	var exp time.Time
	var nbf time.Time

	if value, ok := t.claims["iss"]; ok {
		if _, ok = value.(string); !ok {
			return "", errors.New("Token: Issuer Not Set.")
		}
	}

	if value, ok := t.claims["uid"]; ok {
		if _, ok = value.(string); !ok {
			return "", errors.New("Token: User ID Not Set.")
		}
	}

	if value, ok := t.claims["exp"]; ok {
		if exp, ok = value.(time.Time); !ok {
			return "", errors.New("Token: Expire Time Not Set.")
		}
	}

	if value, ok := t.claims["nbf"]; ok {
		if nbf, ok = value.(time.Time); !ok {
			return "", errors.New("Token: Not Before Time Not Set.")
		}
	}

	if exp.Before(time.Now()) || nbf.Before(time.Now()) || exp.Before(nbf) {
		return "", errors.New("Token: Expire time or Issue at time invalid.")
	}

	bytes, err := json.Marshal(t.claims)
	if nil != err {
		return "", err
	}

	claim := t.encode(bytes)
	sign := t.sign(claim)
	if sign == "" {
		return "", errors.New("Token: Sign error")
	}

	token := claim + "." + sign

	return token, nil
}

func (t *Token) Verify(token string) error {

	tokens := strings.Split(token, ".")
	if size := len(tokens); size != 2 {
		return errors.New("Invaid token format")
	}

	claim := tokens[0]
	sign := tokens[1]

	signTemp := t.sign(claim)
	if signTemp == "" {
		return errors.New("Sign for verify failed")
	}

	if strings.Compare(sign, signTemp) != 0 {
		return errors.New("Invalid signature")
	}

	claimstr, err := t.decode(claim)
	if err != nil {
		return errors.New("Decode claim failed")
	}

	if err := json.Unmarshal([]byte(claimstr), &t.claims); err != nil {
		return errors.New("Unmarshal claims failed")
	}

	if value, ok := t.claims["nbf"]; ok {
		if nbf, ok := value.(time.Time); ok {
			if nbf.After(time.Now()) {
				return errors.New("Token not active")
			}
		} else {
			return errors.New("Bad nbf Value")
		}
	} else {
		return errors.New("Nbf not exist")
	}

	if value, ok := t.claims["exp"]; ok {
		if exp, ok := value.(time.Time); ok {
			if exp.Before(time.Now()) {
				return errors.New("Token expired")
			}
		} else {
			return errors.New("Bad exp value")
		}
	} else {
		return errors.New("Exp not exist")
	}

	return nil
}

func (t *Token) sign(claim string) string {

	h := sha256.New224()
	_, err := h.Write([]byte(claim))
	if nil != err {
		return ""
	}

	sum := h.Sum(nil)

	return t.encode(sum[4:sha256.Size224])
}

func (t *Token) encode(bytes []byte) string {

	b64 := base64.NewEncoding(Base64Pattern)
	return strings.TrimRight(b64.EncodeToString(bytes), "=")
}

func (t *Token) decode(str string) (string, error) {
	b64 := base64.NewEncoding(Base64Pattern)

	if size := len(str) % 4; size > 0 {
		str += strings.Repeat("=", 4-size)
	}

	bytes, err := b64.DecodeString(str)
	if nil != err {
		return "", err
	}

	return string(bytes), nil
}
