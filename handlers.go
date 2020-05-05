package pastry

import "context"

// @todo think about readers the same way github.com/uhthomas/pastry does

type DeliverHandler func(ctx context.Context, key []byte) error
type ForwardHandler func(ctx context.Context, next []byte, key []byte) error
