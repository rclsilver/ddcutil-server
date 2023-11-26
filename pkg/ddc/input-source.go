package ddc

import "context"

type InputSource int64
type InputSourceName string

const (
	inputSource FeatureCode = 60
)

func (c *Client) GetInputSource(ctx context.Context) (InputSource, error) {
	v, err := c.getVcp(ctx, inputSource)
	if err != nil {
		return 0, err
	}
	return InputSource(v), nil
}

func (c *Client) SetInputSource(ctx context.Context, src InputSource) error {
	return c.setVcp(ctx, inputSource, int64(src))
}
