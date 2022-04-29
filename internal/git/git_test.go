package git_test

import (
	"testing"

	"github.com/leoluz/gitops-actions/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	t.Run("will convert git diff output to list of files", func(t *testing.T) {
		// given
		gitDiffOutput := `M	.github/workflows/push-master.yml
D	charts/external-secrets/.helmignore
T	go-actions/tweet/symlink
A	go-actions/tweet/05-retweet-test.txt`
		cloneDir := "/tmp/repo"

		// when
		files := git.ToFiles(cloneDir, gitDiffOutput)

		// then
		assert.Equal(t, 3, len(files))
		f := git.NewFile("go-actions/tweet/05-retweet-test.txt", cloneDir, git.StatusAdded)
		assert.Contains(t, files, f)
		f = git.NewFile(".github/workflows/push-master.yml", cloneDir, git.StatusModified)
		assert.Contains(t, files, f)
		f = git.NewFile("charts/external-secrets/.helmignore", cloneDir, git.StatusDeleted)
		assert.Contains(t, files, f)
	})
}
