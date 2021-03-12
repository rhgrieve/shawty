package urls

import (
	"math/rand"
	"net/url"
	"time"
)

// TODO: Optimize selection of alphabetical chars
func getRandomLetter() rune {
	rand.Seed(time.Now().UnixNano())
	randomChar := rand.Intn(127)

	switch {
	case 97 <= randomChar && randomChar <= 122:
		return rune(randomChar)
	case 65 <= randomChar && randomChar <= 90:
		return rune(randomChar)
	default:
		break
	}

	return getRandomLetter()
}

func generateRandomString() string {
	var runes []rune

	for i := 0; i < 8; i++ {
		randomChar := getRandomLetter()
		runes = append(runes, rune(randomChar))
	}

	return string(runes)
}

type ShortURL struct {
	URL   url.URL `json:"url"`
	Shortcode string  `json:"short"`
}

func NewShortURL(u url.URL) *ShortURL {
	return &ShortURL{
		URL:       u,
		Shortcode: generateRandomString(),
	}
}

func (u *ShortURL) IsEmpty() bool {
	return u.Shortcode == ""
}
