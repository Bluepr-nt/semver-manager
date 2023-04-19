package services

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v51/github"
	"golang.org/x/oauth2"
)

type SemverClient interface {
	ListTags(owner string, repo string) ([]string, error)
}

type DrySemverClient struct {
	Username   string
	Password   string
	Repository string
}

func (drySv *DrySemverClient) ListTags(owner string, repo string) ([]string, error) {
	return []string{}, nil
}

type SemverSvcI interface {
	FetchSemverTags() ([]string, error)
}

type SemverSvc struct {
	client SemverClient
}

type GithubClient struct {
	client *github.Client
}

func newGithubClient(token string) *GithubClient {
	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	return &GithubClient{
		client: github.NewClient(oauthClient),
	}
}

func (ghClient *GithubClient) ListTags(owner string, repo string) ([]string, error) {
	tagList := []string{}
	ctx := context.Background()
	opts := &github.ListOptions{PerPage: 100, Page: 1}

	for {
		tags, resp, err := ghClient.client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return []string{}, fmt.Errorf("ListTags error: %w", err)
		}

		for _, tag := range tags {
			tagList = append(tagList, tag.GetName())
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return tagList, nil
}

func NewSemverSvc(platform, token string) (svc SemverSvc) {
	if platform == "github" {
		newGithubClient(token)

		svc = SemverSvc{client: newGithubClient(token)}
	} else if platform == "dry-run" {
		svc = SemverSvc{client: &DrySemverClient{}}
	}

	return svc
}

func (svSvc *SemverSvc) FetchSemverTags(owner, repo string) (tagList []string, err error) {
	if svSvc.client == nil {
		return nil, fmt.Errorf("git platform client is not defined")
	}

	tagList, err = svSvc.client.ListTags(owner, repo)

	if err != nil {
		return nil, err
	}

	semverTags := FilterSemverTags(tagList)
	return semverTags, err
}

func (svSvc *SemverSvc) GetHighestSemver(semverList []string) (string, error) {
	if semverList == nil || len(semverList) < 1 {
		return "", fmt.Errorf("error the semantic version list is empty")
	}
	semverTags := FilterSemverTags(semverList)
	versions := make([]*semver.Version, len(semverTags))
	for i, tag := range semverTags {
		version, err := semver.NewVersion(tag)
		if err != nil {
			return "", fmt.Errorf("error while sorting semver tags: %w", err)
		}

		versions[i] = version
	}

	sort.Sort(semver.Collection(versions))
	highestVersion := versions[len(versions)-1].String()
	return highestVersion, nil
}

func (svSvc *SemverSvc) FetchHighestSemver(owner, repo string) (string, error) {
	tags, err := svSvc.FetchSemverTags(owner, repo)
	if err != nil {
		return "", err
	}

	highestTag, err := svSvc.GetHighestSemver(tags)
	return highestTag, nil
}

func isRelease(version string) bool {
	regex := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)
	return regex.MatchString(version)
}

func IsSemVer(version string) bool {
	regex := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	return regex.MatchString(version)
}

func FilterSemverTags(tags []string) []string {
	semverTags := []string{}
	for _, tag := range tags {
		if IsSemVer(tag) {
			semverTags = append(semverTags, tag)
		}
	}
	return semverTags
}
