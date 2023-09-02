package auth

import (
	"book-nest/config"
	ma "book-nest/internal/models/auth"
	hh "book-nest/utils/handlerhelper"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	GoogleConf *oauth2.Config
}

func NewAuthHandler() ma.AuthHandler {
	return &AuthHandler{
		GoogleConf: config.GetGoogleConfig(),
	}
}

func (hdl *AuthHandler) GoogleLogin(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "google_login",
		"scope": "auth handler",
	})
	url := hdl.GoogleConf.AuthCodeURL(config.Cfg.GoogleConf.State)

	// redirect to google login page
	c.Redirect(http.StatusTemporaryRedirect, url)
	logger.Info("redirected into: ", url)
}

func (hdl *AuthHandler) GoogleCallback(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "google_callback",
		"scope": "auth handler",
	})

	// state
	state := c.Request.URL.Query().Get("state")
	if state != config.Cfg.GoogleConf.State {
		logger.WithError(errors.New("state does not match"))
		c.JSON(http.StatusUnprocessableEntity, "state does not match")
	}

	// code
	code := c.Request.URL.Query().Get("code")

	// exchange code for token
	token, err := hdl.GoogleConf.Exchange(c, code)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code into token")
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// use google api to get user info
	resp, err := http.Get(config.Cfg.GoogleConf.TokenAccessUrl + token.AccessToken)

	if err != nil {
		logger.WithError(err).Error("failed to get user info")
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("non-OK status code received")
		c.JSON(resp.StatusCode, "Non-OK status code received")
		return
	}

	// Parse the response body as JSON
	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		logger.WithError(err).Error("failed to parse user info as JSON")
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "success",
		Data:    userData,
	})
}
