package common

import "context"

type Thread interface {
	Run(ctx context.Context)
}
