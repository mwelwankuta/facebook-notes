package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

// GetFacebookUserAccessToken gets the access token from Facebook
func GetFacebookUserAccessToken(code string, clientID string, clientSecret string) (string, error) {
	tokenReqURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		config.FbTokenURL, clientID, url.QueryEscape(config.RedirectURI), clientSecret, code,
	)
	resp, err := http.Get(tokenReqURL)
	if err != nil {
		return "", fmt.Errorf("Failed to get access token")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokenResponse map[string]interface{}
	json.Unmarshal(body, &tokenResponse)

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to parse access token")
	}
	return accessToken, nil
}

// FetchUserProfile fetches the user profile from Facebook
func FetchUserProfile(accessToken string) (models.FacebookUser, error) {
	var userDto models.FacebookUser

	// Fetch user profile with access token
	userProfileURL := fmt.Sprintf("%s?fields=id,name,email&access_token=%s", config.FbGraphAPI, accessToken)
	userResp, err := http.Get(userProfileURL)
	if err != nil {
		return userDto, fmt.Errorf("Failed to fetch user profile")
	}
	defer userResp.Body.Close()

	userBody, _ := ioutil.ReadAll(userResp.Body)
	if err := json.Unmarshal(userBody, &userDto); err != nil {
		return userDto, err
	}

	return userDto, nil
}
