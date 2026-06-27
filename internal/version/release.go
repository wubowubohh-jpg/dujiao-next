package version

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CheckResult describes the local update check response.
type CheckResult struct {
	CurrentVersion string     `json:"current_version"`
	LatestVersion  string     `json:"latest_version"`
	HasUpdate      bool       `json:"has_update"`
	ReleaseURL     string     `json:"release_url,omitempty"`
	ReleaseName    string     `json:"release_name,omitempty"`
	ReleaseNotes   string     `json:"release_notes,omitempty"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`
	Source         string     `json:"source"`
}

// ErrRateLimited is kept for API compatibility with older callers.
var ErrRateLimited = errors.New("update check rate limited")

// CheckLatestRelease returns the local version only. External release checks are disabled.
func CheckLatestRelease(ctx context.Context) (*CheckResult, error) {
	_ = ctx
	current := strings.TrimSpace(Version)
	return &CheckResult{
		CurrentVersion: current,
		LatestVersion:  current,
		HasUpdate:      false,
		Source:         "",
	}, nil
}

// IsNewerVersion 判断 latest 是否比 current 更新。返回 (true, nil) 表示需要更新；
// 当任一版本号无法解析时，回退到字符串不相等比较，并返回非空 error 提示调用方。
func IsNewerVersion(latest, current string) (bool, error) {
	l, lErr := parseSemver(latest)
	c, cErr := parseSemver(current)
	if lErr != nil || cErr != nil {
		// 版本号格式无法识别时，仅在两者非空且不相等时认为有更新
		if latest == "" {
			return false, errors.New("latest version is empty")
		}
		return latest != "" && current != "" && latest != current, errors.Join(lErr, cErr)
	}

	for i := range 3 {
		if l[i] > c[i] {
			return true, nil
		}
		if l[i] < c[i] {
			return false, nil
		}
	}
	return false, nil
}

// parseSemver 将 "v1.2.3" / "1.2.3" / "v1.2.3-rc.1" 等格式解析为 [major, minor, patch]
// 仅取主.次.修订三段，忽略预发布和构建元数据
func parseSemver(v string) ([3]int, error) {
	var out [3]int
	s := strings.TrimSpace(v)
	if s == "" {
		return out, errors.New("empty version")
	}
	s = strings.TrimPrefix(s, "v")
	s = strings.TrimPrefix(s, "V")
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) == 0 {
		return out, fmt.Errorf("invalid version: %s", v)
	}
	for i := 0; i < 3 && i < len(parts); i++ {
		n, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			return out, fmt.Errorf("invalid version segment %q in %s", parts[i], v)
		}
		if n < 0 {
			return out, fmt.Errorf("negative version segment in %s", v)
		}
		out[i] = n
	}
	return out, nil
}
