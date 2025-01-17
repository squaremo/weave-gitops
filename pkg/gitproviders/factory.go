package gitproviders

import (
	"fmt"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitlab"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

// GitProviderName holds a Git provider definition.
type GitProviderName string

const (
	GitProviderGitHub GitProviderName = "github"
	GitProviderGitLab GitProviderName = "gitlab"
)

// Config defines the configuration for connecting to a GitProvider.
type Config struct {
	// Provider defines the GitProvider.
	Provider GitProviderName

	// Hostname is the HTTP/S hostname of the Provider,
	// e.g. github.example.com.
	Hostname string

	// Token contains the token used to authenticate with the
	// Provider.
	Token string
}

func buildGitProvider(config Config) (gitprovider.Client, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("no git provider token present")
	}

	var client gitprovider.Client
	var err error
	switch config.Provider {
	case GitProviderGitHub:
		opts := []github.ClientOption{
			github.WithOAuth2Token(config.Token),
		}
		if config.Hostname != "" {
			opts = append(opts, github.WithDomain(config.Hostname))
		}
		if client, err = github.NewClient(opts...); err != nil {
			return nil, err
		}
	case GitProviderGitLab:
		opts := []gitlab.ClientOption{
			gitlab.WithConditionalRequests(true),
		}
		if config.Hostname != "" {
			opts = append(opts, gitlab.WithDomain(config.Hostname))
		}
		if client, err = gitlab.NewClient(config.Token, "", opts...); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported Git provider '%s'", config.Provider)
	}
	return client, err
}
