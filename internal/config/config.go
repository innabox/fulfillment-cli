/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/innabox/fulfillment-common/logging"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	experimentalcredentials "google.golang.org/grpc/experimental/credentials"

	"github.com/innabox/fulfillment-cli/internal/auth"
	"github.com/innabox/fulfillment-cli/internal/oauth"
	"github.com/innabox/fulfillment-cli/internal/packages"
	"github.com/innabox/fulfillment-cli/internal/version"
)

// Config is the type used to store the configuration of the client.
type Config struct {
	TokenScript       string     `json:"token_script"`
	Plaintext         bool       `json:"plaintext,omitempty"`
	Insecure          bool       `json:"insecure,omitempty"`
	Address           string     `json:"address,omitempty"`
	Private           bool       `json:"packages,omitempty"`
	AccessToken       string     `json:"access_token,omitempty"`
	RefreshToken      string     `json:"refresh_token,omitempty"`
	TokenExpiry       time.Time  `json:"token_expiry"`
	OAuthFlow         oauth.Flow `json:"oauth_flow,omitempty"`
	OauthIssuer       string     `json:"oauth_issuer,omitempty"`
	OAuthClientId     string     `json:"oauth_client_id,omitempty"`
	OAuthClientSecret string     `json:"oauth_client_secret,omitempty"`
	OAuthScopes       []string   `json:"oauth_scopes,omitempty"`
}

// Load loads the configuration from the configuration file.
func Load() (cfg *Config, err error) {
	file, err := Location()
	if err != nil {
		return
	}
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		cfg = &Config{}
		err = nil
		return
	}
	if err != nil {
		err = fmt.Errorf("failed to check if config file '%s' exists: %v", file, err)
		return
	}
	data, err := os.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("failed to read config file '%s': %v", file, err)
		return
	}
	cfg = &Config{}
	if len(data) == 0 {
		return
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		err = fmt.Errorf("failed to parse config file '%s': %v", file, err)
		return
	}
	return
}

// Save saves the given configuration to the configuration file.
func Save(cfg *Config) error {
	file, err := Location()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	dir := filepath.Dir(file)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}
	err = os.WriteFile(file, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %v", file, err)
	}
	return nil
}

// Location returns the location of the configuration file.
func Location() (result string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	result = filepath.Join(configDir, "fulfillment-cli", "config.json")
	return
}

// TokenSource creates a token source from the configuration.
func (c *Config) TokenSource(ctx context.Context) (result auth.TokenSource, err error) {
	// Get the logger:
	logger := logging.LoggerFromContext(ctx)

	// Get the token store:
	tokenStore := c.TokenStore()

	// If an OAuth flow has been configured, then use it to create a non interactive OAuth token source:
	if c.OAuthFlow != "" {
		result, err = oauth.NewTokenSource().
			SetLogger(logger).
			SetFlow(c.OAuthFlow).
			SetInteractive(false).
			SetIssuer(c.OauthIssuer).
			SetClientId(c.OAuthClientId).
			SetClientSecret(c.OAuthClientSecret).
			SetScopes(c.OAuthScopes...).
			SetInsecure(c.Insecure).
			SetStore(tokenStore).
			Build()
		return
	}

	// If a token script has been configured, then use it to create a script token source:
	if c.TokenScript != "" {
		result, err = auth.NewScriptTokenSource().
			SetLogger(logger).
			SetScript(c.TokenScript).
			SetStore(tokenStore).
			Build()
		return
	}

	// Finally, if there is an access token try to use it:
	if c.AccessToken != "" {
		result, err = auth.NewStaticTokenSource().
			SetLogger(logger).
			SetToken(&auth.Token{
				Access: c.AccessToken,
			}).
			Build()
		return
	}

	// If we are here then there is no way to get tokens, which is an error:
	err = errors.New("no token source configured")
	return
}

// Conect creates a gRPC connection from the configuration.
func (c *Config) Connect(ctx context.Context, flags *pflag.FlagSet) (result *grpc.ClientConn, err error) {
	// Get the logger:
	logger := logging.LoggerFromContext(ctx)

	// Create a token source:
	tokenSource, err := c.TokenSource(ctx)
	if err != nil {
		return
	}

	// Create the the gRPC credentials:
	tokenCredentials, err := auth.NewTokenCredentials().
		SetLogger(logger).
		SetSource(tokenSource).
		Build()
	if err != nil {
		return
	}

	// Add the credentials to the dial options:
	dialOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(tokenCredentials),
	}

	// Configure use of TLS:
	var transportCreds credentials.TransportCredentials
	if c.Plaintext {
		transportCreds = insecure.NewCredentials()
	} else {
		tlsConfig := &tls.Config{}
		if c.Insecure {
			tlsConfig.InsecureSkipVerify = true
		}

		// TODO: This should have been the non-experimental package, but we need to use this one because
		// currently the OpenShift router doesn't seem to support ALPN, and the regular credentials package
		// requires it since version 1.67. See here for details:
		//
		// https://github.com/grpc/grpc-go/issues/434
		// https://github.com/grpc/grpc-go/pull/7980
		//
		// Is there a way to configure the OpenShift router to avoid this?
		transportCreds = experimentalcredentials.NewTLSWithALPNDisabled(tlsConfig)
	}
	if transportCreds != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))
	}

	// Create the version interceptor:
	versionInterceptor, err := version.NewInterceptor().
		SetLogger(logger).
		Build()
	if err != nil {
		return
	}

	// Create the logging interceptor:
	loggingInterceptor, err := logging.NewInterceptor().
		SetLogger(logger).
		SetFlags(flags).
		Build()
	if err != nil {
		return
	}

	// Add the interceptors to the dial options:
	dialOpts = append(
		dialOpts,
		grpc.WithChainUnaryInterceptor(
			versionInterceptor.UnaryClient,
			loggingInterceptor.UnaryClient,
		),
		grpc.WithChainStreamInterceptor(
			versionInterceptor.StreamClient,
			loggingInterceptor.StreamClient,
		),
	)

	// Create the connection:
	result, err = grpc.NewClient(c.Address, dialOpts...)
	return
}

// Packages returns the list of packages that should be enabled according to the configuration. The public packages
// will always be enabled, but the private packages will be enabled only if the `private` flag is true.
func (c *Config) Packages() []string {
	if c.Private {
		return packages.All
	}
	return packages.Public
}

// TokenStore returns an implementation of the auth.TokenStore interface that loads and saves tokens from/to
// the configuration.
func (c *Config) TokenStore() auth.TokenStore {
	return &configTokenStore{
		config: c,
		lock:   &sync.RWMutex{},
	}
}

type configTokenStore struct {
	config *Config
	lock   *sync.RWMutex
}

func (s *configTokenStore) Load(ctx context.Context) (result *auth.Token, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.config.AccessToken == "" {
		return
	}
	result = &auth.Token{
		Access:  s.config.AccessToken,
		Refresh: s.config.RefreshToken,
		Expiry:  s.config.TokenExpiry,
	}
	return
}

func (s *configTokenStore) Save(ctx context.Context, token *auth.Token) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if token == nil {
		return errors.New("token cannot be nil")
	}
	accessChanged := s.config.AccessToken != token.Access
	refreshChanged := s.config.RefreshToken != token.Refresh
	expiryChanged := s.config.TokenExpiry != token.Expiry
	if !accessChanged && !refreshChanged && !expiryChanged {
		return nil
	}
	s.config.AccessToken = token.Access
	s.config.RefreshToken = token.Refresh
	s.config.TokenExpiry = token.Expiry
	return Save(s.config)
}
