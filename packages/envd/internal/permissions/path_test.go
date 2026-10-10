package permissions

import (
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func currentUser(t *testing.T) *user.User {
	t.Helper()

	u, err := user.Current()
	require.NoError(t, err)

	return u
}

func TestExpandAndResolve(t *testing.T) {
	t.Parallel()

	u := currentUser(t)

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "empty path anchors at the sandbox root",
			path:     "",
			expected: "/",
		},
		{
			name:     "relative path anchors at the sandbox root",
			path:     "foo/bar.txt",
			expected: "/foo/bar.txt",
		},
		{
			name:     "dot-relative path anchors at the sandbox root",
			path:     "./foo",
			expected: "/foo",
		},
		{
			name:     "parent segments are cleaned inside the anchored path",
			path:     "foo/../bar",
			expected: "/bar",
		},
		{
			name:     "absolute path is returned as-is",
			path:     "/foo/bar.txt",
			expected: "/foo/bar.txt",
		},
		{
			name:     "tilde expands to the user home dir",
			path:     "~",
			expected: u.HomeDir,
		},
		{
			name:     "tilde-prefixed path expands to the user home dir",
			path:     "~/foo/bar.txt",
			expected: filepath.Join(u.HomeDir, "foo/bar.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ExpandAndResolve(tt.path, u, nil)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestExpandAndResolveInvalidTilde(t *testing.T) {
	t.Parallel()

	u := currentUser(t)

	_, err := ExpandAndResolve("~other/file.txt", u, nil)
	require.Error(t, err)
}
