package gitcli

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/sourcegraph/sourcegraph/cmd/gitserver/internal/git"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/gitserver/gitdomain"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func (g *gitCLIBackend) ListRefs(ctx context.Context, opt git.ListRefsOpts) (git.RefIterator, error) {
	r, err := g.NewCommand(ctx, WithArguments(buildListRefsArgs(opt)...))
	if err != nil {
		return nil, err
	}

	return &refIterator{
		Closer: r,
		sc:     bufio.NewScanner(r),
	}, nil
}

func buildListRefsArgs(opt git.ListRefsOpts) []string {
	cmdArgs := []string{
		"for-each-ref",
		"--sort", "-refname",
		"--sort", "-creatordate",
		"--sort", "-HEAD",
		// tag * refs/tags/v5.3.1-rc.1 v5.3.1-rc.1 26750071c89a4a6536834a10bf9a97c5e70060a 26750071c89a4a6536834a10bf9a97c5e70060a 11708577416
		"--format", "%(objecttype)%00%(HEAD)%00%(refname)%00%(refname:short)%00%(objectname)%00%(*objectname)%00%(creatordate:unix)",
	}

	if opt.HeadsOnly {
		cmdArgs = append(cmdArgs, "refs/heads/")
	}

	if opt.TagsOnly {
		cmdArgs = append(cmdArgs, "refs/tags/")
	}

	for _, c := range opt.PointsAtCommit {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--points-at=%s", string(c)))
	}

	for _, c := range opt.Contains {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--contains=%s", string(c)))
	}

	return cmdArgs
}

type refIterator struct {
	io.Closer
	sc *bufio.Scanner
}

func (it *refIterator) Next() (*gitdomain.Ref, error) {
	for {
		if it.sc.Scan() {
			line := it.sc.Bytes()
			if len(line) == 0 {
				// Skip over empty output.
				continue
			}
			parts := bytes.Split(line, []byte("\x00"))
			if len(parts) != 7 {
				return nil, errors.Errorf("unexpected output from git for-each-ref %q", string(line))
			}
			// Only tags point to a target object, so for non-tags we set the target
			// equal to the commit ID.
			if string(parts[0]) != "tag" {
				parts[5] = parts[4]
			}
			var createdDate time.Time
			if string(parts[6]) != "" {
				ts, err := strconv.Atoi(string(parts[6]))
				if err != nil {
					return nil, errors.Errorf("unexpected output from git for-each-ref (bad date format) %q", string(line))
				}
				createdDate = time.Unix(int64(ts), 0)
			}
			return &gitdomain.Ref{
				Type:        refTypeForString(string(parts[0])),
				Name:        string(parts[2]),
				ShortName:   string(parts[3]),
				CommitID:    api.CommitID(parts[5]),
				RefOID:      api.CommitID(parts[4]),
				CreatedDate: createdDate,
				IsHead:      string(parts[1]) == "*",
			}, nil
		}
		break
	}
	if err := it.sc.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

func refTypeForString(s string) gitdomain.RefType {
	switch s {
	case "tag":
		return gitdomain.RefTypeTag
	case "commit":
		return gitdomain.RefTypeBranch
	default:
		return gitdomain.RefTypeUnknown
	}
}
