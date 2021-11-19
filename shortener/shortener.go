package shortener

import (
	"crypto/sha256"
	"encoding/base64"
	"gorm.io/gorm"
	"time"
)

const ExpirationDuration = 1 * time.Minute

type ShortUrl struct {
	gorm.Model
	Redirect, Shortened string
	ExpiresAt           time.Time
}

type ShortUrlOption func(url *ShortUrl)

func WithExpiration(duration time.Duration) ShortUrlOption {
	return func(url *ShortUrl) {
		url.ExpiresAt = time.Now().Add(duration)
	}
}

func NewShortUrl(url string, options ...ShortUrlOption) ShortUrl {
	def := ShortUrl{
		Redirect:  url,
		Shortened: "",
		ExpiresAt: time.Now().Add(ExpirationDuration),
	}
	for _, opt := range options {
		opt(&def)
	}
	return def
}

func GetShortUrl(short string, tx *gorm.DB) (result ShortUrl, err error) {
	r := tx.First(&result, "shortened = ? AND expires_at > now()", short)
	err = r.Error
	if err != nil {
		return
	}
	return
}

func (url *ShortUrl) BeforeCreate(tx *gorm.DB) (err error) {
	hasher := sha256.New()
	hasher.Write([]byte(url.Redirect))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	url.Shortened = sha[:6]
	return nil
}
