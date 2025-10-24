/*
Copyright (c) 2025 Red Hat Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
*/

package login

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/innabox/fulfillment-common/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	experimentalcredentials "google.golang.org/grpc/experimental/credentials"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/innabox/fulfillment-cli/internal/auth"
	"github.com/innabox/fulfillment-cli/internal/config"
	"github.com/innabox/fulfillment-cli/internal/oauth"
	"github.com/innabox/fulfillment-cli/internal/templating"
	"github.com/innabox/fulfillment-cli/internal/terminal"
	metadatav1 "github.com/innabox/fulfillment-common/api/metadata/v1"
)

//go:embed templates
var templatesFS embed.FS

func Cmd() *cobra.Command {
	runner := &runnerContext{}
	result := &cobra.Command{
		Use:   "login [FLAGS] <ADDRESS>",
		Short: "Save connection and authentication details",
		RunE:  runner.run,
	}
	flags := result.Flags()
	flags.BoolVar(
		&runner.plaintext,
		"plaintext",
		false,
		"Disables use of TLS for communications with the API server.",
	)
	flags.BoolVar(
		&runner.insecure,
		"insecure",
		false,
		"Disables verification of TLS certificates and host names of the OAuth and API servers.",
	)
	flags.StringVar(
		&runner.address,
		"address",
		os.Getenv("FULFILLMENT_SERVICE_ADDRESS"),
		"Server address.",
	)
	flags.BoolVar(
		&runner.private,
		"private",
		false,
		"Enables use of the private API.",
	)
	flags.StringVar(
		&runner.token,
		"token",
		os.Getenv("FULFILLMENT_SERVICE_TOKEN"),
		"Authentication token",
	)
	flags.StringVar(
		&runner.tokenScript,
		"token-script",
		os.Getenv("FULFILLMENT_SERVICE_TOKEN_SCRIPT"),
		"Shell command that will be executed to obtain the token. For example, to automatically get the "+
			"token of the Kubernetes 'client' service account of the 'example' namespace the value "+
			"could be 'kubectl create token -n example client --duration 1h'. Note that is important "+
			"to quote this shell command correctly, as it will be passed to your shell for "+
			"execution.",
	)
	flags.StringVar(
		&runner.oauthIssuer,
		"oauth-issuer",
		"",
		"OAuth issuer URL. This is optional. By default the first issuer advertised by the server is used.",
	)
	flags.StringVar(
		&runner.oauthFlow,
		"oauth-flow",
		string(oauth.CodeFlow),
		fmt.Sprintf(
			"OAuth flow to use. Must be '%s', '%s' or '%s'.",
			oauth.CodeFlow, oauth.DeviceFlow, oauth.CredentialsFlow,
		),
	)
	flags.StringVar(
		&runner.oauthClientId,
		"oauth-client-id",
		"fulfillment-cli",
		"OAuth client identifier.",
	)
	flags.StringVar(
		&runner.oauthClientSecret,
		"oauth-client-secret",
		"",
		fmt.Sprintf(
			"OAuth client secret. This is required for the '%s' flow.",
			oauth.CredentialsFlow,
		),
	)
	flags.StringSliceVar(
		&runner.oauthScopes,
		"oauth-scopes",
		[]string{},
		"Comma separated list of OAuth scopes to request.",
	)
	flags.MarkHidden("token")
	flags.MarkHidden("token-script")
	flags.MarkHidden("private")
	return result
}

type runnerContext struct {
	logger            *slog.Logger
	console           *terminal.Console
	flags             *pflag.FlagSet
	engine            *templating.Engine
	plaintext         bool
	insecure          bool
	address           string
	private           bool
	token             string
	tokenScript       string
	tokenStore        auth.TokenStore
	oauthIssuer       string
	oauthFlow         string
	oauthClientId     string
	oauthClientSecret string
	oauthScopes       []string
}

func (c *runnerContext) run(cmd *cobra.Command, args []string) error {
	var err error

	// Get the context:
	ctx := cmd.Context()

	// Get the logger, console and flags:
	c.logger = logging.LoggerFromContext(ctx)
	c.console = terminal.ConsoleFromContext(ctx)
	c.flags = cmd.Flags()

	// Create the templating engine:
	c.engine, err = templating.NewEngine().
		SetLogger(c.logger).
		SetFS(templatesFS).
		SetDir("templates").
		Build()
	if err != nil {
		return fmt.Errorf("failed to create templating engine: %w", err)
	}

	// The address used to be specified with a command line flag, but now we also take it from the arguments:
	if c.address == "" {
		if len(args) == 1 {
			c.address = args[0]
		} else {
			return fmt.Errorf("address is mandatory")
		}
	}

	// Configure use of TLS:
	dialOpts := []grpc.DialOption{}
	var transportCreds credentials.TransportCredentials
	if c.plaintext {
		transportCreds = insecure.NewCredentials()
	} else {
		tlsConfig := &tls.Config{}
		if c.insecure {
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

	// Create an empty configuration and a token store that will load/save tokens from/to that configuration:
	cfg := &config.Config{}
	c.tokenStore = cfg.TokenStore()

	// Create the token source:
	tokenSource, err := c.createTokenSource(ctx, dialOpts)
	if err != nil {
		return fmt.Errorf("failed to create token source: %w", err)
	}

	// Try to obtain a token using the token source, as this will trigger the authentication flow and verify that
	// it works correctly.
	_, err = tokenSource.Token(ctx)
	if err != nil {
		return fmt.Errorf("failed to obtain token: %w", err)
	}

	// Populate the configuration:
	cfg.OauthIssuer = c.oauthIssuer
	cfg.OAuthFlow = oauth.Flow(c.oauthFlow)
	cfg.OAuthClientId = c.oauthClientId
	cfg.OAuthClientSecret = c.oauthClientSecret
	cfg.OAuthScopes = c.oauthScopes
	cfg.Plaintext = c.plaintext
	cfg.Insecure = c.insecure
	cfg.Address = c.address
	cfg.Private = c.private

	// Check if the configuration is working by invoking the health check method:
	conn, err := cfg.Connect(ctx, c.flags)
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer conn.Close()
	client := healthv1.NewHealthClient(conn)
	response, err := client.Check(ctx, &healthv1.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if response.Status != healthv1.HealthCheckResponse_SERVING {
		return fmt.Errorf("server is not serving, status is '%s'", response.Status)
	}

	// Everything is working, so we can save the configuration:
	err = config.Save(cfg)
	if err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	return nil
}

func (c *runnerContext) createTokenSource(ctx context.Context,
	dialOpts []grpc.DialOption) (result auth.TokenSource, err error) {
	// Use a token if specified:
	if c.token != "" {
		result, err = auth.NewStaticTokenSource().
			SetLogger(c.logger).
			SetToken(&auth.Token{
				Access: c.token,
			}).
			Build()
		return
	}

	// Use a token script if specified::
	if c.tokenScript != "" {
		result, err = auth.NewScriptTokenSource().
			SetLogger(c.logger).
			SetScript(c.tokenScript).
			SetStore(c.tokenStore).
			Build()
		return
	}

	// If we don't have a token or token script, then we need to use OAuth, and for that we need to find the token
	// issuers trusted by the server. To do so we need an annonymous gRPC connecto to fetch the metadata.
	anonymousConn, err := grpc.NewClient(c.address, dialOpts...)
	if err != nil {
		err = fmt.Errorf("failed to create anonymous gRPC connection: %w", err)
		return
	}
	defer anonymousConn.Close()
	metadataClient := metadatav1.NewMetadataClient(anonymousConn)
	metadataResponse, err := metadataClient.Get(ctx, metadatav1.MetadataGetRequest_builder{}.Build())
	if err != nil {
		err = fmt.Errorf("failed to get metadata: %w", err)
		return
	}
	authnMetadata := metadataResponse.GetAuthn()
	if authnMetadata == nil {
		err = errors.New("no authentication metadata found")
		return
	}
	err = anonymousConn.Close()
	if err != nil {
		c.logger.ErrorContext(
			ctx,
			"Failed to close anonymous connection",
			slog.Any("error", err),
		)
	}
	c.logger.DebugContext(
		ctx,
		"Obtained authentication metadata",
		slog.Any("metadata", authnMetadata),
	)

	// If the issuer has been explicitly specified, use it, otherwise use the first one advertised by the server.
	if c.oauthIssuer == "" {
		advertisedIssuers := authnMetadata.GetTrustedTokenIssuers()
		if len(advertisedIssuers) > 0 {
			c.oauthIssuer = advertisedIssuers[0]
			if len(advertisedIssuers) > 1 {
				c.logger.WarnContext(
					ctx,
					"Server advertises multiple issuers, selecting the first one",
					slog.Any("advertised", advertisedIssuers),
					slog.Any("selected", c.oauthIssuer),
				)
			}
		} else {
			err = errors.New(
				"server advertises no issuers, and no issuer has been specified in the command line",
			)
			return
		}
	} else {
		c.logger.DebugContext(
			ctx,
			"Using issuer from the command line",
			slog.String("issuer", c.oauthIssuer),
		)
	}

	// Create the OAuth token source:
	result, err = oauth.NewTokenSource().
		SetLogger(c.logger).
		SetStore(c.tokenStore).
		SetListener(&oauthFlowListener{
			runner: c,
		}).
		SetInsecure(c.insecure).
		SetInteractive(true).
		SetIssuer(c.oauthIssuer).
		SetFlow(oauth.Flow(c.oauthFlow)).
		SetClientId(c.oauthClientId).
		SetClientSecret(c.oauthClientSecret).
		SetScopes(c.oauthScopes...).
		Build()
	return
}

type oauthFlowListener struct {
	runner *runnerContext
}

func (l *oauthFlowListener) Start(ctx context.Context, event oauth.FlowStartEvent) error {
	switch event.Flow {
	case oauth.CodeFlow:
		return l.startCodeFlow(ctx, event)
	case oauth.DeviceFlow:
		return l.startDeviceFlow(ctx, event)
	default:
		return fmt.Errorf(
			"unsupported flow '%s', must be '%s' or '%s'",
			event.Flow, oauth.CodeFlow, oauth.DeviceFlow,
		)
	}
}

func (l *oauthFlowListener) startCodeFlow(ctx context.Context, event oauth.FlowStartEvent) error {
	l.runner.console.Render(ctx, l.runner.engine, "start_code_flow.txt", map[string]any{
		"AuthorizationUri": event.AuthorizationUri,
	})
	return nil
}

func (l *oauthFlowListener) startDeviceFlow(ctx context.Context, event oauth.FlowStartEvent) error {
	// If the authorizatoin server has provided a complete URL, with the code included, then use it, otherwise use
	// the URL without the code:
	verficationUri := event.VerificationUriComplete
	if verficationUri == "" {
		verficationUri = event.VerificationUri
	}

	// Calculate the expiration time to show to the user::
	now := time.Now()
	expiresIn := humanize.RelTime(now, now.Add(event.ExpiresIn), "from now", "")
	l.runner.console.Render(ctx, l.runner.engine, "start_device_flow.txt", map[string]any{
		"VerificationUri": verficationUri,
		"UserCode":        event.UserCode,
		"ExpiresIn":       expiresIn,
	})
	return nil
}

func (l *oauthFlowListener) End(ctx context.Context, event oauth.FlowEndEvent) error {
	if event.Outcome {
		l.runner.console.Render(ctx, l.runner.engine, "auth_success.txt", nil)
	} else {
		l.runner.console.Render(ctx, l.runner.engine, "auth_failure.txt", nil)
	}
	return nil
}
