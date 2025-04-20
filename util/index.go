package util

import (
	"context"
)

func Init(ctx context.Context) {
	InitGitHubFileUploader(ctx)
	InitJwtC(ctx)
}
