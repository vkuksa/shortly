package domain

import (
	"encoding/base64"
	"time"
)

type UUID string

func NewUUIDfromString(str string) UUID {
	return UUID(base64.URLEncoding.EncodeToString([]byte(str)))
}

func (u UUID) String() string {
	return string(u)
}

type Link struct {
	// UUID is also a shortened link
	UUID UUID `json:"uuid"`
	// Original url to redirrect to
	URL string `json:"url"`
	// Amount of times that link was reddirrected
	Count int `json:"count"`

	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (l *Link) ResetExpiration() {
	if l == nil {
		panic("trying to reset expiration of nil link")
	}

	l.ExpiresAt = time.Now().Add(24 * time.Hour)
}

func NewLink(url string) *Link {
	return &Link{
		UUID:      NewUUIDfromString(url),
		URL:       url,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}
