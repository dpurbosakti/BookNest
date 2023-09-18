package calendarhelper

import (
	"book-nest/config"
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

func GetClient(token *oauth2.Token) *http.Client {
	config := config.GetGoogleConfig()
	client := config.Client(context.Background(), token)
	return client
}
