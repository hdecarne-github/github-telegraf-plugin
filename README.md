[![Downloads](https://img.shields.io/github/downloads/hdecarne-github/github-telegraf-plugin/total.svg)](https://github.com/hdecarne-github/github-telegraf-plugin/releases)
[![Build](https://github.com/hdecarne-github/github-telegraf-plugin/actions/workflows/build.yml/badge.svg)](https://github.com/hdecarne-github/github-telegraf-plugin/actions/workflows/build.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=hdecarne-github_github-telegraf-plugin&metric=coverage)](https://sonarcloud.io/summary/new_code?id=hdecarne-github_github-telegraf-plugin)

## About github-telegraf-plugin
This [Telegraf](https://github.com/influxdata/telegraf) input plugin gathers repository stats from [GitHub](https://github.com/). It uses GitHub's [REST API](https://docs.github.com/en/rest) to retrieve the stats.

### Installation
To install the plugin you have to download a suitable [release archive](https://github.com/hdecarne-github/github-telegraf-plugin/releases) and extract it or build it from source by cloning the repository and issueing a simple
```
make
```
To build the plugin, Go version 1.16 or higher is required. The resulting plugin binary will be written to **./build/bin**.
Copy the either extracted or built plugin binary to a location of your choice (e.g. /usr/local/bin/telegraf/).

### Configuration
This is an [external plugin](https://github.com/influxdata/telegraf/blob/master/docs/EXTERNAL_PLUGINS.md) which has to be integrated via Telegraf's [excecd plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd).

To use it you have to create a plugin specific config file (e.g. /etc/telegraf/github.conf) with following template content:
```toml
[[inputs.github]]
  ## The repositories (<owner>/<repo>) to query
  repos = ["influxdata/telegraf"]
  ## The API base URL to use for API access (empty URL defaults to https://api.github.com/)
  # api_base_url = ""
  ## The Personal Access Token to use for API access
  # access_token = ""
  ## The http timeout to use (in seconds)
  # timeout = 5
  ## Enable debug output
  # debug = false
```
The most important setting is the **repos** line. It defines the repositories (<owner>/<name>) to query. At least one repository has to be defined.

To enable the plugin within your Telegraf instance, add the following section to your **telegraf.conf**
```toml
[[inputs.execd]]
  command = ["/usr/local/bin/telegraf/github-telegraf-plugin", "-config", "/etc/telegraf/github.conf", "-poll_interval", "3600s"]
  signal = "none"
```
Make sure to choose a high poll interval, to not waste your rate limit. As the github stats are low-traffic stats, there is furthermore no need to poll in high frequency mode.

### License
This project is subject to the the MIT License.
See [LICENSE](./LICENSE) information for details.