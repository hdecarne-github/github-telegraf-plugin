// github.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	githubApi "github.com/google/go-github/v44/github"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"golang.org/x/oauth2"
)

type GitHub struct {
	Repos       []string `toml:"repos"`
	APIBaseURL  string   `toml:"api_base_url"`
	AccessToken string   `toml:"access_token"`

	Timeout int  `toml:"timeout"`
	Debug   bool `toml:"debug"`

	Log telegraf.Logger
}

func NewGitHub() *GitHub {
	return &GitHub{
		Repos:       []string{},
		AccessToken: "",
		Timeout:     10,
	}
}

func (plugin *GitHub) SampleConfig() string {
	return `
  ## The repositories (<owner>/<repo>) to query
  repos = ["influxdata/telegraf"]
  ## The API base URL to use for API access (empty URL defaults to https://api.github.com/)
  # api_base_url = ""
  ## The Personal Access Token to use for API access
  # access_token = ""
  ## The http timeout to use (in seconds)
  # timeout = 10
  ## Enable debug output
  # debug = false
 `
}

func (plugin *GitHub) Description() string {
	return "Gather GitHub stats"
}

func (plugin *GitHub) Gather(a telegraf.Accumulator) error {
	if len(plugin.Repos) == 0 {
		return errors.New("github: Empty repo list")
	}
	ctx := context.Background()
	client, err := plugin.getClient(ctx)
	if err != nil {
		return err
	}
	for _, repo := range plugin.Repos {
		a.AddError(plugin.processRepo(ctx, client, a, repo))
	}
	return nil
}

func (plugin *GitHub) processRepo(ctx context.Context, client *githubApi.Client, a telegraf.Accumulator, repo string) error {
	if plugin.Debug {
		plugin.Log.Infof("Processing repo: %s", repo)
	}
	repoOwner, repoName, err := plugin.splitRepoId(repo)
	if err != nil {
		return err
	}
	repoInfo, _, err := client.Repositories.Get(ctx, repoOwner, repoName)
	if err != nil {
		return err
	}
	repoReleases, _, err := client.Repositories.ListReleases(ctx, repoOwner, repoName, nil)
	if err != nil {
		return err
	}
	totalDownloadCount := 0
	for _, repoRelease := range repoReleases {
		for _, repoReleaseAsset := range repoRelease.Assets {
			totalDownloadCount += repoReleaseAsset.GetDownloadCount()
		}
	}

	viewTimestamp := time.Time{}
	var totalViews int
	var uniqueViews int

	if plugin.AccessToken != "" {
		repoTrafficViews, _, err := client.Repositories.ListTrafficViews(ctx, repoOwner, repoName, &githubApi.TrafficBreakdownOptions{Per: "day"})
		if err != nil {
			return err
		}
		for _, repoTrafficView := range repoTrafficViews.Views {
			if repoTrafficView.Timestamp.After(viewTimestamp) {
				viewTimestamp = repoTrafficView.Timestamp.Time
				totalViews = repoTrafficView.GetCount()
				uniqueViews = repoTrafficView.GetUniques()
			}
		}
	}
	tags := make(map[string]string)
	tags["github_repo"] = repo
	fields := make(map[string]interface{})
	fields["forks_count"] = repoInfo.ForksCount
	fields["stargazers_count"] = repoInfo.StargazersCount
	fields["subscribers_count"] = repoInfo.SubscribersCount
	fields["total_download_count"] = totalDownloadCount
	fields["total_views"] = totalViews
	fields["unique_views"] = uniqueViews
	a.AddCounter("github_info", fields, tags)
	return nil
}

func (plugin *GitHub) splitRepoId(repo string) (string, string, error) {
	repoParts := strings.Split(repo, "/")
	if len(repoParts) != 2 {
		return "", "", fmt.Errorf("github: Invalid repo identifier '%s'", repo)
	}
	return repoParts[0], repoParts[1], nil
}

func (plugin *GitHub) getClient(ctx context.Context) (*githubApi.Client, error) {
	if plugin.Debug {
		plugin.Log.Debug("Creating GitHub client...")
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: time.Duration(plugin.Timeout) * time.Second,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(plugin.Timeout) * time.Second,
	}
	if plugin.AccessToken != "" {
		if plugin.Debug {
			plugin.Log.Debug("Using oauth2 access token...")
		}
		token := &oauth2.Token{AccessToken: plugin.AccessToken}
		tokenSource := oauth2.StaticTokenSource(token)
		httpClient = oauth2.NewClient(ctx, tokenSource)
	}
	if plugin.APIBaseURL != "" {
		if plugin.Debug {
			plugin.Log.Debug("Using API base URL: '%s'...", plugin.APIBaseURL)
		}
		return githubApi.NewEnterpriseClient(plugin.APIBaseURL, "", httpClient)
	}
	return githubApi.NewClient(httpClient), nil
}

func init() {
	inputs.Add("github", func() telegraf.Input {
		return NewGitHub()
	})
}
