package workers

import (
	"context"
	"time"

	"golang-fave/engine/basket"

	"github.com/vladimirok5959/golang-worker/worker"
)

func BasketCleaner(sb *basket.Basket) *worker.Worker {
	return worker.New(func(ctx context.Context, w *worker.Worker, o *[]worker.Iface) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			if sb, ok := (*o)[0].(*basket.Basket); ok {
				sb.Cleanup()
			}
			return
		}
	}, &[]worker.Iface{
		sb,
	})
}
