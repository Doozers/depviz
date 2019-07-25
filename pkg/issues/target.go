package issues

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Target string

func ParseTargets(inputs []string) (Targets, error) {
	targetMap := map[Target]bool{}
	for _, input := range inputs {
		target, err := ParseTarget(input)
		if err != nil {
			return nil, err
		}
		targetMap[target] = true
	}
	targets := Targets{}
	for target := range targetMap {
		targets = append(targets, target)
	}
	sort.Slice(targets, func(i, j int) bool {
		return string(targets[i]) < string(targets[j])
	})

	return targets, nil
}

func ParseTarget(input string) (Target, error) {
	// check if input is a local path
	if _, err := os.Stat(input); err == nil {
		return "", fmt.Errorf("filesystem target are not yet supported")
	}

	// parse issue
	str := input
	issue := ""
	parts := strings.Split(str, "/issues/")
	switch len(parts) {
	case 1:
	case 2:
		str = parts[0]
		issue = parts[1]
	default:
		return "", fmt.Errorf("invalid target: %q", input)
	}
	parts = strings.Split(str, "#")
	switch len(parts) {
	case 1:
	case 2:
		str = parts[0]
		issue = parts[1]
	default:
		return "", fmt.Errorf("invalid target: %q", input)
	}

	// parse scheme
	parts = strings.Split(str, "/")
	if len(parts) < 3 {
		str = fmt.Sprintf("https://github.com/%s", str)
	}

	if !strings.Contains(str, "://") {
		str = fmt.Sprintf("https://%s", str)
	}

	// append issue
	if issue != "" {
		_, err := strconv.Atoi(issue)
		if err != nil {
			return "", fmt.Errorf("invalid target (issue): %q", input)
		}
		str = str + "/issues/" + issue
	}

	return Target(str), nil
}

func (t Target) Issue() string {
	parts := strings.Split(string(t), "/issues/")
	switch len(parts) {
	case 1:
		return ""
	case 2:
		return parts[1]
	default:
		panic("invalid target")
	}
}

func (t Target) ProjectURL() string {
	return strings.Split(string(t), "/issues/")[0]
}

func (t Target) Namespace() string {
	u, err := url.Parse(t.ProjectURL())
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Path, "/")[1:]
	return strings.Join(parts[:len(parts)-1], "/")
}

func (t Target) Project() string {
	parts := strings.Split(t.ProjectURL(), "/")
	return parts[len(parts)-1]
}

func (t Target) Path() string {
	return fmt.Sprintf("%s/%s", t.Namespace(), t.Project())
}

func (t Target) Canonical() string { return string(t) }

func (t Target) Driver() ProviderDriver {
	if strings.Contains(string(t), "github.com") { // github.com
		return GithubDriver
	}
	// FIXME: support GithubEnterpriseDriver (github.company.com)
	if strings.Contains(string(t), "gitlab") { // gitlab.com and gitlab.company.com
		return GitlabDriver
	}
	if strings.Contains(string(t), "jira") || strings.Contains(string(t), "atlassian") {
		return JiraDriver
	}
	return GitlabDriver // default to gitlab
}

func (t Target) ProviderURL() string {
	switch t.Driver() {
	case GithubDriver:
		return "https://github.com"
	case GitlabDriver, JiraDriver:
		u, err := url.Parse(string(t))
		if err != nil {
			return ""
		}
		return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	default:
		panic("should not happen")
	}
}

type Targets []Target

func (t Targets) UniqueProjects() Targets {
	urlMap := map[string]bool{}
	for _, target := range t {
		urlMap[target.ProjectURL()] = true
	}

	urls := []string{}
	for url := range urlMap {
		urls = append(urls, url)
	}
	sort.Strings(urls)

	filtered := Targets{}
	for _, url := range urls {
		filtered = append(filtered, Target(url))
	}
	return filtered
}
