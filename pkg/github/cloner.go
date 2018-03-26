package github

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Cloner, clones github repositories
type Cloner struct{}

// CloneRepository clones a github repository to given target path
func (c *Cloner) CloneRepository(ctx context.Context, out io.Writer, target string, repo domain.SourceRepository, username, token string) error {
	cloneUrl := fmt.Sprintf("https://%s:%s@github.com/%s/%s", username, token, repo.Owner, repo.Name)
	cmd := exec.Command("git", "clone", cloneUrl, target)
	cmd.Stderr = out
	cmd.Stdout = out
	return cmd.Run()
}
