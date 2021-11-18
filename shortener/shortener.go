package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"gorm.io/gorm"
	"time"
)

var ExpirationDuration time.Duration

func init() {
	var err error
	ExpirationDuration, err = time.ParseDuration("+48h")
	if err == nil {
		panic("Cannot initialize shortener")
	}
}

type ShortUrl struct {
	gorm.Model
	redirect, shortened string
	expiresAt           time.Time
}

func NewShortUrl(url string) ShortUrl {
	return ShortUrl{
		redirect:  url,
		shortened: "",
		expiresAt: time.Now().Add(ExpirationDuration),
	}
}

func GetShortUrl(short string, tx *gorm.DB) (result ShortUrl, err error) {
	r := tx.First(&result, "shortened = ? AND expires_at > localtime", short)
	err = r.Error
	if err != nil {
		return
	}
	return
}

func (url *ShortUrl) BeforeCreate(tx *gorm.DB) (err error) {
	hasher := sha256.New()
	hasher.Write([]byte(url.redirect))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	url.shortened = sha[:6]
	return nil
}
