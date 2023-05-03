package semverSvc

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
	listTags(owner string, repo string) ([]string, error)
}

type DrySemverClient struct {
	Username   string
	Password   string
	Repository string
}

func (drySv *DrySemverClient) listTags(owner string, repo string) ([]string, error) {
	return []string{}, nil
}

type SemverSvcI interface {
	FetchSemverTags() ([]string, error)
}

type SemverSvc struct {
	client SemverClient
}

type Filters struct {
	Highest bool
	Release bool
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

func (ghClient *GithubClient) listTags(owner string, repo string) ([]string, error) {
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

func (svSvc *SemverSvc) FetchSemverTags(owner, repo string, filters *Filters) (tagList []string, err error) {
	if svSvc.client == nil {
		return nil, fmt.Errorf("git platform client is not defined")
	}

	tagList, err = svSvc.client.listTags(owner, repo)

	if err != nil {
		return nil, err
	}

	semverTags, err := svSvc.FilterSemverTags(tagList, filters)
	return semverTags, err
}

func (svSvc *SemverSvc) FilterHighestSemver(semverList []string) (string, error) {
	if semverList == nil || len(semverList) < 1 {
		return "", fmt.Errorf("error the semantic version list is empty")
	}
	semverTags, err := svSvc.FilterSemverTags(semverList, nil)
	if err != nil {
		return "", err
	}
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

func Filter(semverList []string, filters *Filters) {

}
func IsRelease(version string) bool {
	regex := regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)
	return regex.MatchString(version)
}

func IsSemver(version string) bool {
	regex := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	return regex.MatchString(version)
}

func AreSemver(versions []string) bool {
	for _, version := range versions {
		if !IsSemver(version) {
			return false
		}
	}
	return true
}

func (svSvc *SemverSvc) FilterSemverTags(tags []string, filters *Filters) ([]string, error) {
	var err error
	filteredTags := removeNonCompliantTags(tags)
	filteredTags, err = svSvc.SortTags(filteredTags)
	if err != nil {
		return nil, err
	}

	if filters != nil && filters.Release {
		filteredTags, err = svSvc.FilterSemverRelease(filteredTags)
		if err != nil {
			return nil, err
		}
	}

	if filters != nil && filters.Highest && len(filteredTags) > 0 {
		highestVersion, err := svSvc.FilterHighestSemver(filteredTags)
		if err != nil {
			return nil, err
		}
		filteredTags = []string{highestVersion}
	}

	return filteredTags, nil
}

func (svSvc *SemverSvc) FilterSemverRelease(tags []string) ([]string, error) {
	var releaseVersions []*semver.Version
	versions, err := svSvc.stringsToVersions(tags)
	if err != nil {
		return nil, err
	}
	for _, v := range versions {
		if v.Prerelease() == "" {
			releaseVersions = append(releaseVersions, v)
		}
	}
	releaseTag := svSvc.versionsToStrings(releaseVersions)
	return releaseTag, nil
}

func (svSvc *SemverSvc) SortTags(tags []string) ([]string, error) {
	versions, err := svSvc.stringsToVersions(tags)
	if err != nil {
		return nil, err
	}
	sort.Sort(semver.Collection(versions))
	sortedTags := svSvc.versionsToStrings(versions)
	return sortedTags, nil
}
func removeNonCompliantTags(tags []string) []string {
	var semverTags []string
	for _, tag := range tags {
		if IsSemver(tag) {
			semverTags = append(semverTags, tag)
		}
	}
	return semverTags
}

func (svSvc *SemverSvc) versionsToStrings(versions []*semver.Version) []string {
	var tags []string
	for _, version := range versions {
		tags = append(tags, version.String())
	}
	return tags
}

func (svSvc *SemverSvc) stringsToVersions(tags []string) ([]*semver.Version, error) {
	var semverVersions []*semver.Version

	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			return semverVersions, err
		}
		semverVersions = append(semverVersions, v)
	}
	return semverVersions, nil
}
