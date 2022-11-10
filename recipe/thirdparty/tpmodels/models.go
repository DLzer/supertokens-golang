/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package tpmodels

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeRedirectURIQueryParams = map[string]interface{}
type TypeOAuthTokens = map[string]interface{}

type TypeRawUserInfoFromProvider struct {
	FromIdTokenPayload map[string]interface{}
	FromUserInfoAPI    map[string]interface{}
}

type TypeSupertokensUserInfo struct {
	ThirdPartyUserId string
	EmailInfo        *EmailStruct
}

type TypeUserInfo struct {
	ThirdPartyUserId        string
	Email                   *EmailStruct
	RawUserInfoFromProvider TypeRawUserInfoFromProvider
}

type EmailStruct struct {
	ID         string `json:"id"`
	IsVerified bool   `json:"isVerified"`
}

type TypeAuthorisationRedirect struct {
	URLWithQueryParams string
	PKCECodeVerifier   *string
}

type TypeRedirectURIInfo struct {
	RedirectURIOnProviderDashboard string                     `json:"redirectURIOnProviderDashboard"`
	RedirectURIQueryParams         TypeRedirectURIQueryParams `json:"redirectURIQueryParams"`
	PKCECodeVerifier               *string                    `json:"pkceCodeVerifier"`
}

type TypeFrom string

const (
	FromIdTokenPayload TypeFrom = "idTokenPayload"
	FromUserInfoAPI    TypeFrom = "userInfoAPI"
)

type TypeUserInfoMap struct {
	From          TypeFrom
	UserId        string
	Email         string
	EmailVerified string
}

type TypeProvider struct {
	ID string

	GetAuthorisationRedirectURL    func(clientType *string, tenantId *string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (TypeAuthorisationRedirect, error)
	ExchangeAuthCodeForOAuthTokens func(clientType *string, tenantId *string, redirectInfo TypeRedirectURIInfo, userContext supertokens.UserContext) (TypeOAuthTokens, error) // For apple, add userInfo from callbackInfo to oAuthTOkens
	GetUserInfo                    func(clientType *string, tenantId *string, oAuthTokens TypeOAuthTokens, userContext supertokens.UserContext) (TypeUserInfo, error)
}

/*
TypeProviderInterface allows us to define the config directly in the thirdparty.Init
instead of calling a function that returns the instance of TypeProvider. This is
needed because we may have to update the client config array from the `findProvider`,
which wouldn't be possible with the baked TypeProvider. findProvider also now builds
based on the updated client configs.

GetID needs to be a function because golang interface can contain only methods
*/
type TypeProviderInterface interface {
	GetID() string
	Build() TypeProvider
}

type User struct {
	ID         string `json:"id"`
	TimeJoined uint64 `json:"timeJoined"`
	Email      string `json:"email"`
	ThirdParty struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	} `json:"thirdParty"`
}

type TypeInputSignInAndUp struct {
	Providers []TypeProviderInterface
}

type TypeNormalisedInputSignInAndUp struct {
	Providers            []TypeProviderInterface
	GetUserPoolForTenant func(tenantId string, userContext supertokens.UserContext) (string, error)
}

type TypeInput struct {
	SignInAndUpFeature TypeInputSignInAndUp
	Override           *OverrideStruct
}

type TypeNormalisedInput struct {
	SignInAndUpFeature TypeNormalisedInputSignInAndUp
	Override           OverrideStruct
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}
