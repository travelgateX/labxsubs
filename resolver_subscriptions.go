package labx_subs

import (
	"context"
	"labxsubs/model"
	"sync"

	"labxsubs/redis"

	"github.com/markbates/going/randx"
)

type subscriptionResolver struct{ *Resolver }

var mutex sync.Mutex

func (r *subscriptionResolver) SupplierCreated(ctx context.Context) (<-chan *model.Supplier, error) {
	id := randx.String(8)
	createdSupplierEvent := make(chan *model.Supplier, 1)

	mutex.Lock()
	redis.CreatedSupplierChannel[id] = createdSupplierEvent
	mutex.Unlock()

	go func() {
		<-ctx.Done()
		mutex.Lock()
		delete(redis.CreatedSupplierChannel, id)
		mutex.Unlock()
	}()

	return createdSupplierEvent, nil
}
