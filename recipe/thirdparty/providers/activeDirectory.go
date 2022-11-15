package providers

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/crypto/pkcs12"
)

const activeDirectoryID = "active-directory"

func ActiveDirectory(input tpmodels.ProviderInput) tpmodels.TypeProvider {
	if input.ThirdPartyID == "" {
		input.ThirdPartyID = activeDirectoryID
	}

	if input.Config.UserInfoMap.FromUserInfoAPI.UserId == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromUserInfoAPI.Email == "" {
		input.Config.UserInfoMap.FromUserInfoAPI.Email = "email"
	}

	if input.Config.UserInfoMap.FromIdTokenPayload.UserId == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.UserId = "sub"
	}
	if input.Config.UserInfoMap.FromIdTokenPayload.Email == "" {
		input.Config.UserInfoMap.FromIdTokenPayload.Email = "email"
	}

	oOverride := input.Override

	input.Override = func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := provider.GetConfigForClientType
		provider.GetConfigForClientType = func(clientType *string, input tpmodels.ProviderConfig, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, input, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.OIDCDiscoveryEndpoint == "" {
				config.OIDCDiscoveryEndpoint = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/", config.AdditionalConfig["directoryId"])
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" && config.AdditionalConfig["certificate"] != nil {
				if config.TokenParams == nil {
					config.TokenParams = map[string]interface{}{}
				}
				config.TokenParams["client_assertion_type"] = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
				ca, err := getADClientAssertion(config)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.TokenParams["client_assertion"] = ca
			}

			return config, err
		}

		if oOverride != nil {
			provider = oOverride(provider)
		}
		return provider
	}

	return NewProvider(input)
}

func getADClientAssertion(config tpmodels.ProviderConfigForClientType) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 3600,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Audience:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.AdditionalConfig["directoryId"]),
		Subject:   getActualClientIdFromDevelopmentClientId(config.ClientID),
		Issuer:    getActualClientIdFromDevelopmentClientId(config.ClientID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	thumbBytes, err := hex.DecodeString(config.AdditionalConfig["certificateThumbprint"].(string))
	if err != nil {
		return "", err
	}
	token.Header["x5t"] = base64.StdEncoding.EncodeToString(thumbBytes)
	token.Header["alg"] = "RS256"

	pfxbytes, err := base64.StdEncoding.DecodeString(config.AdditionalConfig["certificate"].(string))
	if err != nil {
		return "", err
	}
	pk, _, err := pkcs12.Decode(pfxbytes, "")
	if err != nil {
		return "", err
	}
	if pk == nil {
		return "", errors.New("private key not found")
	}

	return token.SignedString(pk)
}
