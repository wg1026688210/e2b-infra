package filesystem

import (
	"errors"
	"os/user"
	"testing"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/e2b-dev/infra/packages/envd/internal/services/spec/filesystem"
)

func TestRemoveEmptyPath(t *testing.T) {
	t.Parallel()

	u, err := user.Current()
	require.NoError(t, err)

	svc := mockService()
	ctx := authn.SetInfo(t.Context(), u)

	req := connect.NewRequest(&filesystem.RemoveRequest{Path: ""})
	_, err = svc.Remove(ctx, req)
	require.Error(t, err)

	var connectErr *connect.Error
	ok := errors.As(err, &connectErr)
	assert.True(t, ok, "expected error to be of type *connect.Error")
	assert.Equal(t, connect.CodeInvalidArgument, connectErr.Code())
}
