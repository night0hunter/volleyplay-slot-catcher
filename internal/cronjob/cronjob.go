package cronjob

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type CatchSlotCronHandler interface {
	CatchCron(ctx context.Context) error
}

type CatchSlotCron struct {
	handler CatchSlotCronHandler
}

func NewCatchSlotCron(handler CatchSlotCronHandler) *CatchSlotCron {
	return &CatchSlotCron{
		handler: handler,
	}
}

func (c *CatchSlotCron) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.handler.CatchCron(ctx)
	if err != nil {
		fmt.Println(errors.Wrap(err, "handler.Cron"))
	}
}
