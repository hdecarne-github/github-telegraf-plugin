// github_test.go
//
// Copyright (C) 2022 Holger de Carne
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
//
package github

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	gh := NewGitHub()
	require.NotNil(t, gh)
}

func TestSampleConfig(t *testing.T) {
	gh := NewGitHub()
	sampleConfig := gh.SampleConfig()
	require.NotNil(t, sampleConfig)
}

func TestDescription(t *testing.T) {
	gh := NewGitHub()
	description := gh.Description()
	require.NotNil(t, description)
}

func TestGather1(t *testing.T) {
	testServerHandler := &testServerHandler{Debug: true}
	testServer := httptest.NewServer(testServerHandler)
	defer testServer.Close()
	gh := NewGitHub()
	gh.Repos = []string{"repo_owner/repo_name"}
	gh.APIBaseURL = testServer.URL
	gh.AccessToken = "secret_token"
	gh.Log = createDummyLogger()
	gh.Debug = testServerHandler.Debug

	var a testutil.Accumulator

	require.NoError(t, a.GatherError(gh.Gather))
	require.True(t, a.HasMeasurement("github_info"))
}

func createDummyLogger() *dummyLogger {
	log.SetOutput(os.Stderr)
	return &dummyLogger{}
}

type dummyLogger struct{}

func (l *dummyLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *dummyLogger) Error(args ...interface{}) {
	log.Print(args...)
}

func (l *dummyLogger) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *dummyLogger) Debug(args ...interface{}) {
	log.Print(args...)
}

func (l *dummyLogger) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *dummyLogger) Warn(args ...interface{}) {
	log.Print(args...)
}

func (l *dummyLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *dummyLogger) Info(args ...interface{}) {
	log.Print(args...)
}

type testServerHandler struct {
	Debug bool
}

func (tsh *testServerHandler) ServeHTTP(out http.ResponseWriter, request *http.Request) {
	requestURL := request.URL.String()
	if tsh.Debug {
		log.Printf("test: request URL: %s", requestURL)
	}
	if requestURL == "/api/v3/repos/repo_owner/repo_name" {
		tsh.serveRepositoryInfo(out, request)
	} else if requestURL == "/api/v3/repos/repo_owner/repo_name/releases" {
		tsh.serveRepositoryReleases(out, request)
	} else if requestURL == "/api/v3/repos/repo_owner/repo_name/traffic/views?per=day" {
		tsh.serveRepositoryTrafficViews(out, request)
	}
}

const testResourceLight = `
{
	"stargazers_count": 1,
	"forks_count": 2,
	"subscribers_count": 3
}
`

func (tsh *testServerHandler) serveRepositoryInfo(out http.ResponseWriter, request *http.Request) {
	tsh.writeJSON(out, testResourceLight)
}

const testRepositoryReleases = `
[
  {
    "assets": [
      {
        "download_count": 1
      },
      {
        "download_count": 1
      },
      {
        "download_count": 1
      },
      {
        "download_count": 2
      },
      {
        "download_count": 1
      },
      {
        "download_count": 2
      }
    ]
  },
  {
    "assets": [
      {
        "download_count": 2
      },
      {
        "download_count": 4
      },
      {
        "download_count": 2
      },
      {
        "download_count": 3
      },
      {
        "download_count": 3
      },
      {
        "download_count": 4
      }
    ]
  },
  {
    "assets": [

    ]
  }
]
`

func (tsh *testServerHandler) serveRepositoryReleases(out http.ResponseWriter, request *http.Request) {
	tsh.writeJSON(out, testRepositoryReleases)
}

const testRepositoryTrafficViews = `
{
	"count": 14850,
	"uniques": 3782,
	"views": [
	  {
		"timestamp": "2022-10-10T00:00:00Z",
		"count": 440,
		"uniques": 143
	  },
	  {
		"timestamp": "2022-10-11T00:00:00Z",
		"count": 1308,
		"uniques": 414
	  },
	  {
		"timestamp": "2022-10-12T00:00:00Z",
		"count": 1486,
		"uniques": 452
	  },
	  {
		"timestamp": "2022-10-13T00:00:00Z",
		"count": 1170,
		"uniques": 401
	  },
	  {
		"timestamp": "2022-10-14T00:00:00Z",
		"count": 868,
		"uniques": 266
	  },
	  {
		"timestamp": "2022-10-15T00:00:00Z",
		"count": 495,
		"uniques": 157
	  },
	  {
		"timestamp": "2022-10-16T00:00:00Z",
		"count": 524,
		"uniques": 175
	  },
	  {
		"timestamp": "2022-10-17T00:00:00Z",
		"count": 1263,
		"uniques": 412
	  },
	  {
		"timestamp": "2022-10-18T00:00:00Z",
		"count": 1402,
		"uniques": 417
	  },
	  {
		"timestamp": "2022-10-19T00:00:00Z",
		"count": 1394,
		"uniques": 424
	  },
	  {
		"timestamp": "2022-10-20T00:00:00Z",
		"count": 1492,
		"uniques": 448
	  },
	  {
		"timestamp": "2022-10-21T00:00:00Z",
		"count": 1153,
		"uniques": 332
	  },
	  {
		"timestamp": "2022-10-22T00:00:00Z",
		"count": 566,
		"uniques": 168
	  },
	  {
		"timestamp": "2022-10-23T00:00:00Z",
		"count": 675,
		"uniques": 184
	  },
	  {
		"timestamp": "2022-10-24T00:00:00Z",
		"count": 614,
		"uniques": 237
	  }
	]
  }
`

func (tsh *testServerHandler) serveRepositoryTrafficViews(out http.ResponseWriter, request *http.Request) {
	tsh.writeJSON(out, testRepositoryTrafficViews)
}

func (tsh *testServerHandler) writeJSON(out http.ResponseWriter, json string) {
	out.Header().Add("Content-Type", "application/json")
	_, _ = out.Write([]byte(json))
}
