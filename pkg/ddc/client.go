package ddc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type FeatureCode int

type Client struct{}

var getVcpRe = regexp.MustCompile(`\(sl=0x([0-9a-f]{2})\)$`)

func (c *Client) getVcp(ctx context.Context, featureCode FeatureCode) (int64, error) {
	var out bytes.Buffer

	if err := ddcutil(ctx, nil, &out, os.Stderr, "getvcp", fmt.Sprintf("%d", featureCode)); err != nil {
		return 0, fmt.Errorf("unable to execute the ddcutil command: %w", err)
	}

	m := getVcpRe.FindStringSubmatch(strings.TrimSpace(out.String()))

	if m == nil {
		return 0, fmt.Errorf("invalid ddcutil command output")
	}

	res, err := strconv.ParseInt(m[1], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to convert the value: %w", err)
	}

	return res, nil
}

func (c *Client) setVcp(ctx context.Context, featureCode FeatureCode, value int64) error {
	if err := ddcutil(ctx, nil, nil, os.Stderr, "setvcp", fmt.Sprintf("%d", featureCode), fmt.Sprintf("%d", value)); err != nil {
		return fmt.Errorf("unable to execute the ddcutil command: %w", err)
	}

	return nil
}

func ddcutil(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "ddcutil", args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
