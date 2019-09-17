// Copyright Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/dghubble/sling"
	"github.com/jawher/mow.cli"
)

const GhOrg = "joyent"
const GhRepo = "kosh"

// GithubRelease represents a 'release' for a Github project
type GithubRelease struct {
	URL        string         `json:"html_url"`
	TagName    string         `json:"tag_name"`
	SemVer     semver.Version `json:"-"` // Will be set to 0.0.0 if no releases are found
	Body       string         `json:"body"`
	Name       string         `json:"name"`
	Assets     []GithubAsset  `json:"assets"`
	PreRelease bool           `json:"prerelease"`
	Upgrade    bool           `json:"-"`
}

type GithubReleases []GithubRelease

func (g GithubReleases) Len() int {
	return len(g)
}

func (g GithubReleases) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GithubReleases) Less(i, j int) bool {
	var iSem, jSem semver.Version

	if g[i].TagName == "" {
		iSem = semver.MustParse("0.0.0")
	} else {
		iSem = semver.MustParse(
			strings.TrimLeft(g[i].TagName, "v"),
		)
	}

	if g[j].TagName == "" {
		jSem = semver.MustParse("0.0.0")
	} else {
		jSem = semver.MustParse(
			strings.TrimLeft(g[j].TagName, "v"),
		)
	}

	return iSem.GT(jSem) // reversing sort
}

// GithubAsset represents a file inside of a github release
type GithubAsset struct {
	URL                string `json:"url"`
	Name               string `json:"name"`
	State              string `json:"state"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var ErrNoGithubRelease = errors.New("no appropriate github release found")

// LatestGithubRelease returns some fields from the latest Github Release
// that matches our major version
func LatestGithubRelease() (gh GithubRelease, err error) {
	releases := make(GithubReleases, 0)

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases",
		GhOrg,
		GhRepo,
	)

	_, err = sling.New().Get(url).Receive(&releases, nil)

	if err != nil {
		return gh, err
	}

	sort.Sort(releases)

	sem := CleanVersion(Version)

	for _, r := range releases {
		if r.PreRelease {
			continue
		}
		if r.TagName == "" {
			continue
		}
		r.SemVer = CleanVersion(
			strings.TrimLeft(r.TagName, "v"),
		)

		// Two things are at play here. First, we only care about releases that
		// share our major number. This prevents someone from updating from
		// v1.42 to v2.0 which might contain breaking changes.
		// Second, since we've sorted these in descending order, the first
		// release we find with our major number is the largest. We don't need
		// to dig any further.
		if r.SemVer.Major == sem.Major {
			if r.SemVer.GT(sem) {
				r.Upgrade = true
			}
			return r, nil
		}
	}

	return gh, ErrNoGithubRelease
}

func GithubReleasesSince(start semver.Version) GithubReleases {
	releases := make(GithubReleases, 0)

	diff := make(GithubReleases, 0)

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases",
		GhOrg,
		GhRepo,
	)

	_, err := sling.New().
		Get(url).Receive(&releases, nil)

	if err != nil {
		return diff
	}

	sort.Sort(releases)
	sem := CleanVersion(Version)

	for _, r := range releases {
		if r.PreRelease {
			continue
		}
		if r.TagName == "" {
			continue
		}
		r.SemVer = CleanVersion(
			strings.TrimLeft(r.TagName, "v"),
		)

		// We will not show changelogs for releases that do not share our major
		// version. Since we don't allow users to upgrade across a major
		// version, it's silly to show them those changelogs.
		if r.SemVer.Major == sem.Major {
			if r.SemVer.GT(start) {
				diff = append(diff, r)
			}
		}
	}

	sort.Sort(diff)

	return diff
}

// CleanVersion removes a "v" prefix, and anything after a dash
// For example, pass in v2.99.10-abcde-dirty and get back a semver containing
// 2.29.10
// Why? Git and Semver differ in their notions of what those extra bits mean.
// In Git, they mean "v2.99.10, plus some other stuff that happend". In semver,
// they indicate that this is a prerelease of v2.99.10. Obviously this screws
// up comparisions. This function lets us clean that stuff out so we can get a
// clean comparison
func CleanVersion(version string) semver.Version {
	bits := strings.Split(strings.TrimLeft(version, "v"), "-")
	return semver.MustParse(bits[0])
}

func init() {
	App.Command("update", "Information about the version status of the application and if it can be updated", func(cmd *cli.Cmd) {

		cmd.Command("status", "Verify that we have the most recent revision", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				gh, err := LatestGithubRelease()
				if err != nil {
					if err == ErrNoGithubRelease {
						fmt.Printf(
							"This is %s. No upgrade is available.\n",
							Version,
						)
						return
					}

					panic(err)
				}

				if gh.Upgrade {
					fmt.Printf(
						"This is %s. An upgrade to %s is available\n",
						Version,
						gh.TagName,
					)
				} else {
					fmt.Printf(
						"This is %s. No upgrade is available.\n",
						Version,
					)
				}

			}
		})

		cmd.Command("changelog", "Display the latest changelog", func(cmd *cli.Cmd) {
			cmd.Action = func() {
				releases := GithubReleasesSince(CleanVersion(Version))
				if len(releases) == 0 {
					fmt.Println("No changelog found")
					return
				}

				sort.Sort(sort.Reverse(releases))

				for _, gh := range releases {
					// I'm not going to try and fully sanitize the output
					// for a shell environment but removing the markdown
					// backticks seems like a no-brainer for safety.
					re := regexp.MustCompile("`")
					body := gh.Body
					re.ReplaceAllLiteralString(body, "'")
					fmt.Printf("# Version %s Changelog:\n\n", gh.TagName)
					fmt.Println(body)
					fmt.Println("- - -\n")
				}
			}
		})

		cmd.Command("self", "Update the running application to the latest release", func(cmd *cli.Cmd) {
			force := cmd.BoolOpt("force", false, "Update the binary even if it appears we are on the current release")
			cmd.Action = func() {
				gh, err := LatestGithubRelease()

				if err != nil {
					if err == ErrNoGithubRelease {
						panic("no upgrade available")
					}
					panic(err)
				}

				if !*force {
					if !gh.Upgrade {
						panic("No upgrade required")
					}
				}
				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Attempting to upgrade from %s to %s...\n",
						Version,
						gh.SemVer,
					)

					fmt.Fprintf(
						os.Stderr,
						"Detected OS to be '%s' and arch to be '%s'\n",
						runtime.GOOS,
						runtime.GOARCH,
					)
				}
				// What platform are we on?
				// XXX lookingFor := fmt.Sprintf("kosh-%s-%s", runtime.GOOS, runtime.GOARCH)
				lookingFor := fmt.Sprintf("conch-%s-%s", runtime.GOOS, runtime.GOARCH)
				downloadURL := ""

				// Is this a supported platform
				for _, a := range gh.Assets {
					if a.Name == lookingFor {
						downloadURL = a.BrowserDownloadURL
					}
				}
				if downloadURL == "" {
					panic(fmt.Errorf(
						"could not find an appropriate binary for %s-%s",
						runtime.GOOS,
						runtime.GOARCH,
					))
				}

				//*****  Download the binary
				conchBin, err := updaterDownloadFile(downloadURL)
				if err != nil {
					panic(err)
				}

				//***** Verify checksum

				// This assumes our build system is being sensible about file names.
				// At time of writing, it is.
				shaURL := downloadURL + ".sha256"
				shaBin, err := updaterDownloadFile(shaURL)
				if err != nil {
					panic(err)
				}

				// The checksum file looks like "thisisahexstring ./kosh-os-arch"
				bits := strings.Split(string(shaBin[:]), " ")
				remoteSum := bits[0]

				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Server-side SHA256 sum: %s\n",
						remoteSum,
					)
				}

				h := sha256.New()
				h.Write(conchBin)
				sum := hex.EncodeToString(h.Sum(nil))

				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"SHA256 sum of downloaded binary: %s\n",
						sum,
					)
				}

				if sum == remoteSum {
					if !API.JsonOnly {
						fmt.Fprintf(
							os.Stderr,
							"SHA256 checksums match\n",
						)
					}
				} else {
					panic(fmt.Errorf(
						"!!! SHA of downloaded file does not match the provided SHA sum: '%s' != '%s'",
						sum,
						remoteSum,
					))
				}

				//***** Write out the binary
				binPath, err := os.Executable()
				if err != nil {
					panic(err)
				}

				fullPath, err := filepath.EvalSymlinks(binPath)
				if err != nil {
					panic(err)
				}
				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Detected local binary path: %s\n",
						fullPath,
					)
				}
				existingStat, err := os.Lstat(fullPath)
				if err != nil {
					panic(err)
				}
				// On sensible operating systems, we can't open and write to our
				// own binary, because it's in use. We can, however, move a file
				// into that place.

				newPath := fmt.Sprintf("%s-%s", fullPath, gh.SemVer)
				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Writing to temp file '%s'\n",
						newPath,
					)
				}
				if err := ioutil.WriteFile(newPath, conchBin, existingStat.Mode()); err != nil {
					panic(err)
				}

				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Renaming '%s' to '%s'\n",
						newPath,
						fullPath,
					)
				}

				if err := os.Rename(newPath, fullPath); err != nil {
					panic(err)
				}

				if !API.JsonOnly {
					fmt.Fprintf(
						os.Stderr,
						"Successfully upgraded from %s to %s\n",
						Version,
						gh.SemVer,
					)
				}

			}
		})
	})

}

func updaterDownloadFile(downloadURL string) (data []byte, err error) {
	if !API.JsonOnly {
		fmt.Fprintf(
			os.Stderr,
			"Downloading '%s'\n",
			downloadURL,
		)
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return data, err
	}

	if resp.StatusCode != 200 {
		return data, fmt.Errorf(
			"could not download '%s' (status %d)",
			downloadURL,
			resp.StatusCode,
		)
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return data, err
}
