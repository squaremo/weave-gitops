package gitproviders

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/weave-gitops/pkg/utils"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

type ProviderAccountType string

const (
	AccountTypeUser ProviderAccountType = "user"
	AccountTypeOrg  ProviderAccountType = "organization"
)

// GitProvider Handler
//counterfeiter:generate . GitProvider
type GitProvider interface {
	CreateRepository(name string, owner string, private bool) error
	RepositoryExists(name string, owner string) (bool, error)
	DeployKeyExists(owner, repoName string) (bool, error)
	UploadDeployKey(owner, repoName string, deployKey []byte) error
	CreatePullRequestToUserRepo(userRepRef gitprovider.UserRepositoryRef, targetBranch string, newBranch string, files []gitprovider.CommitFile, commitMessage string, prTitle string, prDescription string) (gitprovider.PullRequest, error)
	CreatePullRequestToOrgRepo(orgRepRef gitprovider.OrgRepositoryRef, targetBranch string, newBranch string, files []gitprovider.CommitFile, commitMessage string, prTitle string, prDescription string) (gitprovider.PullRequest, error)
	GetAccountType(owner string) (ProviderAccountType, error)
}

// making sure it implements the interface
var _ GitProvider = defaultGitProvider{}

type defaultGitProvider struct {
	provider gitprovider.Client
}

func New(config Config) (GitProvider, error) {
	provider, err := buildGitProvider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build git provider: %w", err)
	}

	return defaultGitProvider{
		provider: provider,
	}, nil
}

func (p defaultGitProvider) RepositoryExists(name string, owner string) (bool, error) {
	ownerType, err := p.GetAccountType(owner)
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	if ownerType == AccountTypeOrg {
		orgRef := gitprovider.OrgRepositoryRef{
			OrganizationRef: gitprovider.OrganizationRef{Domain: github.DefaultDomain, Organization: owner},
			RepositoryName:  name,
		}
		if _, err := p.provider.OrgRepositories().Get(ctx, orgRef); err != nil {
			return false, err
		}

		return true, nil
	}

	userRepoRef := gitprovider.UserRepositoryRef{
		UserRef:        gitprovider.UserRef{Domain: github.DefaultDomain, UserLogin: owner},
		RepositoryName: name,
	}
	if _, err := p.provider.UserRepositories().Get(ctx, userRepoRef); err != nil {
		return false, err
	}

	return true, nil
}

func (p defaultGitProvider) CreateRepository(name string, owner string, private bool) error {
	visibility := gitprovider.RepositoryVisibilityPrivate
	if !private {
		visibility = gitprovider.RepositoryVisibilityPublic
	}
	repoInfo := NewRepositoryInfo("Weave Gitops repo", visibility)

	repoCreateOpts := &gitprovider.RepositoryCreateOptions{
		AutoInit:        gitprovider.BoolVar(true),
		LicenseTemplate: gitprovider.LicenseTemplateVar(gitprovider.LicenseTemplateApache2),
	}

	ownerType, err := p.GetAccountType(owner)
	if err != nil {
		return err
	}

	if ownerType == AccountTypeOrg {
		orgRef := NewOrgRepositoryRef(github.DefaultDomain, owner, name)
		if err = p.CreateOrgRepository(orgRef, repoInfo, repoCreateOpts); err != nil {
			return err
		}
	} else {
		userRef := NewUserRepositoryRef(github.DefaultDomain, owner, name)
		if err = p.CreateUserRepository(userRef, repoInfo, repoCreateOpts); err != nil {
			return err
		}
	}

	return nil
}

func (p defaultGitProvider) DeployKeyExists(owner, repoName string) (bool, error) {
	deployKeyName := "weave-gitops-deploy-key"

	ownerType, err := p.GetAccountType(owner)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	defer ctx.Done()
	switch ownerType {
	case AccountTypeOrg:
		orgRef := NewOrgRepositoryRef(github.DefaultDomain, owner, repoName)
		orgRepo, err := p.provider.OrgRepositories().Get(ctx, orgRef)
		if err != nil {
			return false, fmt.Errorf("error getting org repo reference for owner %s, repo %s, %s ", owner, repoName, err)
		}
		_, err = orgRepo.DeployKeys().Get(ctx, deployKeyName)
		if err != nil && !strings.Contains(err.Error(), "key is already in use") {
			if errors.Is(err, gitprovider.ErrNotFound) {
				return false, nil
			} else {
				return false, fmt.Errorf("error getting deploy key %s for repo %s. %s", deployKeyName, repoName, err)
			}
		} else {
			return true, nil
		}

	case AccountTypeUser:
		userRef := NewUserRepositoryRef(github.DefaultDomain, owner, repoName)
		userRepo, err := p.provider.UserRepositories().Get(ctx, userRef)
		if err != nil {
			return false, fmt.Errorf("error getting user repo reference for owner %s, repo %s, %s ", owner, repoName, err)
		}
		_, err = userRepo.DeployKeys().Get(ctx, deployKeyName)
		if err != nil && !strings.Contains(err.Error(), "key is already in use") {
			if errors.Is(err, gitprovider.ErrNotFound) {
				return false, nil
			} else {
				return false, fmt.Errorf("error getting deploy key %s for repo %s. %s", deployKeyName, repoName, err)
			}
		} else {
			return true, nil
		}
	default:
		return false, fmt.Errorf("account type not supported %s", ownerType)
	}
}

func (p defaultGitProvider) UploadDeployKey(owner, repoName string, deployKey []byte) error {
	deployKeyName := "weave-gitops-deploy-key"
	deployKeyInfo := gitprovider.DeployKeyInfo{
		Name: deployKeyName,
		Key:  deployKey,
	}

	ownerType, err := p.GetAccountType(owner)
	if err != nil {
		return err
	}

	ctx := context.Background()
	defer ctx.Done()
	switch ownerType {
	case AccountTypeOrg:
		orgRef := NewOrgRepositoryRef(github.DefaultDomain, owner, repoName)
		orgRepo, err := p.provider.OrgRepositories().Get(ctx, orgRef)
		if err != nil {
			return fmt.Errorf("error getting org repo reference for owner %s, repo %s, %s ", owner, repoName, err)
		}
		fmt.Println("uploading deploy key")
		_, err = orgRepo.DeployKeys().Create(ctx, deployKeyInfo)
		if err != nil {
			return fmt.Errorf("error uploading deploy key %s", err)
		}
		if err = utils.WaitUntil(os.Stdout, time.Second, time.Second*30, func() error {
			_, err = orgRepo.DeployKeys().Get(ctx, deployKeyName)
			return err
		}); err != nil {
			return fmt.Errorf("error verifying deploy key %s existance for repo %s. %s", deployKeyName, repoName, err)
		}
	case AccountTypeUser:
		userRef := NewUserRepositoryRef(github.DefaultDomain, owner, repoName)
		userRepo, err := p.provider.UserRepositories().Get(ctx, userRef)
		if err != nil {
			return fmt.Errorf("error getting user repo reference for owner %s, repo %s, %s ", owner, repoName, err)
		}
		fmt.Println("uploading deploy key")
		_, err = userRepo.DeployKeys().Create(ctx, deployKeyInfo)
		if err != nil {
			return fmt.Errorf("error uploading deploy key %s", err)
		}
		if err = utils.WaitUntil(os.Stdout, time.Second, time.Second*30, func() error {
			_, err = userRepo.DeployKeys().Get(ctx, deployKeyName)
			return err
		}); err != nil {
			return fmt.Errorf("error verifying deploy key %s existance for repo %s. %s", deployKeyName, repoName, err)
		}
	default:
		return fmt.Errorf("account type not supported %s", ownerType)
	}

	return nil
}

func (p defaultGitProvider) GetAccountType(owner string) (ProviderAccountType, error) {
	ctx := context.Background()
	defer ctx.Done()

	_, err := p.provider.Organizations().Get(ctx, gitprovider.OrganizationRef{
		Domain:       github.DefaultDomain,
		Organization: owner,
	})

	if err != nil {
		if errors.Is(err, gitprovider.ErrNotFound) {
			return AccountTypeUser, nil
		}

		return "", fmt.Errorf("could not get account type %s", err)
	}

	return AccountTypeOrg, nil
}

func (p defaultGitProvider) GetRepoInfo(accountType ProviderAccountType, owner string, repoName string) error {
	ctx := context.Background()
	defer ctx.Done()

	switch accountType {
	case AccountTypeOrg:
		if err := p.GetOrgRepo(owner, repoName); err != nil {
			return err
		}
	case AccountTypeUser:
		if err := p.GetUserRepo(owner, repoName); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected account type %s", accountType)
	}

	return nil
}

func (p defaultGitProvider) GetOrgRepo(org string, repoName string) error {
	ctx := context.Background()
	defer ctx.Done()

	orgRepoRef := NewOrgRepositoryRef(github.DefaultDomain, org, repoName)

	_, err := p.provider.OrgRepositories().Get(ctx, orgRepoRef)
	if err != nil {
		return fmt.Errorf("error getting org repository %s", err)
	}

	return nil
}

func (p defaultGitProvider) GetUserRepo(user string, repoName string) error {
	ctx := context.Background()
	defer ctx.Done()

	userRepoRef := NewUserRepositoryRef(github.DefaultDomain, user, repoName)

	_, err := p.provider.UserRepositories().Get(ctx, userRepoRef)
	if err != nil {
		return err
	}

	return nil
}

func (p defaultGitProvider) CreateOrgRepository(orgRepoRef gitprovider.OrgRepositoryRef, repoInfo gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryCreateOption) error {
	ctx := context.Background()
	defer ctx.Done()

	_, err := p.provider.OrgRepositories().Create(ctx, orgRepoRef, repoInfo, opts...)
	if err != nil {
		return fmt.Errorf("error creating repo %s", err)
	}

	return p.waitUntilRepoCreated(AccountTypeOrg, orgRepoRef.Organization, orgRepoRef.RepositoryName)
}

func (p defaultGitProvider) CreateUserRepository(userRepoRef gitprovider.UserRepositoryRef, repoInfo gitprovider.RepositoryInfo, opts ...gitprovider.RepositoryCreateOption) error {
	ctx := context.Background()
	defer ctx.Done()

	_, err := p.provider.UserRepositories().Create(ctx, userRepoRef, repoInfo, opts...)
	if err != nil {
		return fmt.Errorf("error creating repo %s", err)
	}

	return p.waitUntilRepoCreated(AccountTypeUser, userRepoRef.UserLogin, userRepoRef.RepositoryName)
}

func (p defaultGitProvider) CreatePullRequestToUserRepo(userRepRef gitprovider.UserRepositoryRef, targetBranch string, newBranch string, files []gitprovider.CommitFile, commitMessage string, prTitle string, prDescription string) (gitprovider.PullRequest, error) {
	ctx := context.Background()

	ur, err := p.provider.UserRepositories().Get(ctx, userRepRef)
	if err != nil {
		return nil, fmt.Errorf("error getting info for repo [%s] err [%s]", userRepRef.String(), err)
	}

	if targetBranch == "" {
		targetBranch = *ur.Get().DefaultBranch
	}

	commits, err := ur.Commits().ListPage(ctx, targetBranch, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("error getting commits for repo[%s] err [%s]", userRepRef.String(), err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("targetBranch[%s] does not exists", targetBranch)
	}

	latestCommit := commits[0]

	if err := ur.Branches().Create(ctx, newBranch, latestCommit.Get().Sha); err != nil {
		return nil, fmt.Errorf("error creating branch[%s] for repo[%s] err [%s]", newBranch, userRepRef.String(), err)
	}

	if _, err := ur.Commits().Create(ctx, newBranch, commitMessage, files); err != nil {
		return nil, fmt.Errorf("error creating commit for branch[%s] for repo[%s] err [%s]", newBranch, userRepRef.String(), err)
	}

	pr, err := ur.PullRequests().Create(ctx, prTitle, newBranch, targetBranch, prDescription)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request[%s] for branch[%s] for repo[%s] err [%s]", prTitle, newBranch, userRepRef.String(), err)
	}

	return pr, nil
}

func (p defaultGitProvider) CreatePullRequestToOrgRepo(orgRepRef gitprovider.OrgRepositoryRef, targetBranch string, newBranch string, files []gitprovider.CommitFile, commitMessage string, prTitle string, prDescription string) (gitprovider.PullRequest, error) {
	ctx := context.Background()

	ur, err := p.provider.OrgRepositories().Get(ctx, orgRepRef)
	if err != nil {
		return nil, fmt.Errorf("error getting info for repo [%s] err [%s]", orgRepRef.String(), err)
	}

	if targetBranch == "" {
		targetBranch = *ur.Get().DefaultBranch
	}

	commits, err := ur.Commits().ListPage(ctx, targetBranch, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("error getting commits for repo[%s] err [%s]", orgRepRef.String(), err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("targetBranch[%s] does not exists", targetBranch)
	}

	latestCommit := commits[0]

	if err := ur.Branches().Create(ctx, newBranch, latestCommit.Get().Sha); err != nil {
		return nil, fmt.Errorf("error creating branch[%s] for repo[%s] err [%s]", newBranch, orgRepRef.String(), err)
	}

	if _, err := ur.Commits().Create(ctx, newBranch, commitMessage, files); err != nil {
		return nil, fmt.Errorf("error creating commit for branch[%s] for repo[%s] err [%s]", newBranch, orgRepRef.String(), err)
	}

	pr, err := ur.PullRequests().Create(ctx, prTitle, newBranch, targetBranch, prDescription)
	if err != nil {
		return nil, fmt.Errorf("error creating pull request[%s] for branch[%s] for repo[%s] err [%s]", prTitle, newBranch, orgRepRef.String(), err)
	}

	return pr, nil
}

func NewRepositoryInfo(description string, visibility gitprovider.RepositoryVisibility) gitprovider.RepositoryInfo {
	return gitprovider.RepositoryInfo{
		Description: &description,
		Visibility:  &visibility,
	}
}

func NewOrgRepositoryRef(domain, org, repoName string) gitprovider.OrgRepositoryRef {
	return gitprovider.OrgRepositoryRef{
		RepositoryName: repoName,
		OrganizationRef: gitprovider.OrganizationRef{
			Domain:       domain,
			Organization: org,
		},
	}
}

func NewUserRepositoryRef(domain, user, repoName string) gitprovider.UserRepositoryRef {
	return gitprovider.UserRepositoryRef{
		RepositoryName: repoName,
		UserRef: gitprovider.UserRef{
			Domain:    domain,
			UserLogin: user,
		},
	}
}

func (p defaultGitProvider) waitUntilRepoCreated(ownerType ProviderAccountType, owner, name string) error {
	if err := utils.WaitUntil(os.Stdout, time.Second, time.Second*30, func() error {
		return p.GetRepoInfo(ownerType, owner, name)
	}); err != nil {
		return fmt.Errorf("could not verify repo existence %s", err)
	}
	return nil
}
