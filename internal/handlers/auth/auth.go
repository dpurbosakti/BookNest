package auth

import (
	"book-nest/config"
	ma "book-nest/internal/models/auth"
	hh "book-nest/utils/handlerhelper"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	AuthService ma.AuthService
	GoogleConf  *oauth2.Config
	TwitterConf *oauth2.Config
	GithubConf  *oauth2.Config
}

func NewAuthHandler(authService ma.AuthService) ma.AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		GoogleConf:  config.GetGoogleConfig(),
		TwitterConf: config.GetTwitterConfig(),
		GithubConf:  config.GetGithubConfig(),
	}
}

func (hdl *AuthHandler) Login(c *gin.Context) {
	var authReq ma.LoginRequest

	err := c.ShouldBindJSON(&authReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "login",
		"scope": "auth handler",
		"data":  authReq,
	})
	if err != nil {
		logger.WithError(err).Error("failed to bind user")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	token, err := hdl.AuthService.Login(authReq)
	if err != nil {
		if err.Error() == "password incorrect" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"messaage": "ok",
		"token":    token,
	})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "state does not match"})
	}

	// code
	code := c.Request.URL.Query().Get("code")

	// exchange code for token
	token, err := hdl.GoogleConf.Exchange(c, code)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code into token")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	}
	if token == nil {
		logger.WithError(err).Error("failed to get token")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		// use google api to get user info
		resp, err := http.Get(config.Cfg.GoogleConf.TokenAccessUrl + token.AccessToken)

		if err != nil {
			logger.WithError(err).Error("failed to get user info")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Error("non-OK status code received")
			c.JSON(resp.StatusCode, gin.H{"message": "Non-OK status code received"})
			return
		}

		// Parse the response body as JSON
		var userData *ma.GoogleResponse
		err = json.NewDecoder(resp.Body).Decode(&userData)
		if err != nil {
			logger.WithError(err).Error("failed to parse user info as JSON")
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		userData.OauthAccessToken = token.AccessToken
		jwtToken, err := hdl.AuthService.LoginByGoogle(userData)
		if err != nil {
			logger.WithError(err).Error("failed to login by google")
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"token":   jwtToken,
		})
	}

}

func (hdl *AuthHandler) TwitterLogin(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "twitter_login",
		"scope": "auth handler",
	})

	// Generate the Twitter OAuth2 authorization URL
	url := hdl.TwitterConf.AuthCodeURL("")

	// Redirect the user to the Twitter login page
	c.Redirect(http.StatusTemporaryRedirect, url)
	logger.Info("redirected into: ", url)
}

func (hdl *AuthHandler) TwitterCallback(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "twitter_callback",
		"scope": "auth handler",
	})

	// // state
	// state := c.Request.URL.Query().Get("state")
	// if state != config.Cfg.GoogleConf.State {
	// 	logger.WithError(errors.New("state does not match"))
	// 	c.JSON(http.StatusUnprocessableEntity, "state does not match")
	// }

	// code
	code := c.Request.URL.Query().Get("code")
	fmt.Println("code: ", code)

	// exchange code for token
	token, err := hdl.TwitterConf.Exchange(c, code)
	fmt.Println("token: ", token)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code into token")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	}

	// Create an HTTP client with the access token
	httpClient := oauth2.NewClient(c, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.AccessToken}))

	// Make the authenticated API request
	resp, err := httpClient.Get(config.Cfg.TwitterConf.ApiEndpoint)
	if err != nil {
		logger.WithError(err).Error("failed to make API request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("non-OK status code received")
		c.JSON(resp.StatusCode, gin.H{"message": "Non-OK status code received"})
		return
	}

	// Parse the response body as JSON
	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		logger.WithError(err).Error("failed to parse user info as JSON")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "success",
		Data:    userData,
	})
}
func (hdl *AuthHandler) GithubLogin(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "github_login",
		"scope": "auth handler",
	})
	url := hdl.GithubConf.AuthCodeURL("")
	fmt.Println("url: ", url)

	// // redirect to twitter login page
	c.Redirect(http.StatusTemporaryRedirect, url)
	logger.Info("redirected into: ", url)
}

func (hdl *AuthHandler) GithubCallback(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "github_callback",
		"scope": "auth handler",
	})

	// code
	code := c.Request.URL.Query().Get("code")

	// exchange code for token
	token, err := hdl.GithubConf.Exchange(c, code)
	if err != nil {
		logger.WithError(err).Error("failed to exchange code into token")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
	}
	if token == nil {
		logger.WithError(err).Error("failed to get token")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	// Create an HTTP client with the access token
	httpClient := hdl.GithubConf.Client(c, token)
	if httpClient == nil {
		logger.WithError(err).Error("failed to initiate http client")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	// Make the authenticated API request
	resp, err := httpClient.Get(config.Cfg.GithubConf.TokenAccessUrl)
	if err != nil {
		logger.WithError(err).Error("failed to make API request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("non-OK status code received")
		c.JSON(resp.StatusCode, gin.H{"message": "Non-OK status code received"})
		return
	}

	// Parse the response body as JSON
	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		logger.WithError(err).Error("failed to parse user info as JSON")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "success",
		Data:    userData,
	})
}
