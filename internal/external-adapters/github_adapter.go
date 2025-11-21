// Package externaladapters provides external service adapters.
package externaladapters

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/ochairo/wand/internal/domain/interfaces"
)

// GitHubAdapter implements the GitHubClient interface
type GitHubAdapter struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubAdapter creates a new GitHubAdapter
func NewGitHubAdapter(token string) interfaces.GitHubClient {
	ctx := context.Background()
	var client *github.Client

	if token != "" {
		client = github.NewClient(nil).WithAuthToken(token)
	} else {
		client = github.NewClient(nil)
	}

	return &GitHubAdapter{
		client: client,
		ctx:    ctx,
	}
}

// GetLatestRelease gets the latest release for a repository
func (g *GitHubAdapter) GetLatestRelease(owner, repo string) (*interfaces.GitHubRelease, error) {
	release, _, err := g.client.Repositories.GetLatestRelease(g.ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	return g.convertRelease(release), nil
}

// GetRelease gets a specific release by tag
func (g *GitHubAdapter) GetRelease(owner, repo, tag string) (*interfaces.GitHubRelease, error) {
	release, _, err := g.client.Repositories.GetReleaseByTag(g.ctx, owner, repo, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get release %s: %w", tag, err)
	}

	return g.convertRelease(release), nil
}

// ListReleases lists all releases for a repository
func (g *GitHubAdapter) ListReleases(owner, repo string) ([]*interfaces.GitHubRelease, error) {
	opts := &github.ListOptions{PerPage: 100}

	var allReleases []*interfaces.GitHubRelease
	for {
		releases, resp, err := g.client.Repositories.ListReleases(g.ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}

		for _, release := range releases {
			allReleases = append(allReleases, g.convertRelease(release))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

// DownloadAsset downloads a release asset to the specified path
func (g *GitHubAdapter) DownloadAsset(asset *interfaces.GitHubAsset, destPath string) error {
	resp, err := http.Get(asset.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	outFile, err := os.Create(destPath) //nolint:gosec // G304: destPath from application config
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = outFile.Close() }()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// convertRelease converts a GitHub release to our interface type
func (g *GitHubAdapter) convertRelease(release *github.RepositoryRelease) *interfaces.GitHubRelease {
	assets := make([]interfaces.GitHubAsset, 0, len(release.Assets))
	for _, asset := range release.Assets {
		assets = append(assets, interfaces.GitHubAsset{
			Name:        asset.GetName(),
			DownloadURL: asset.GetBrowserDownloadURL(),
		})
	}

	return &interfaces.GitHubRelease{
		TagName: release.GetTagName(),
		Assets:  assets,
	}
}
