package labx_subs

import (
	"context"
	"encoding/json"
	"errors"
	"labxsubs/model"
	"labxsubs/redis"

	"github.com/markbates/going/randx"
)

type supplierResolver struct{ *Resolver }

func (r *supplierResolver) Accesses(ctx context.Context, obj *model.Supplier) ([]*model.Access, error) {
	var accesses []*model.Access

	redisAccesses, err := redis.GetClient().AllAccesses()
	var mapAccess = make(map[string]model.Access)

	for _, redAcc := range redisAccesses {
		mapAccess[redAcc.ID] = *redAcc
	}

	for id := range obj.AccessesID {
		strid := obj.AccessesID[id]
		access, exist := mapAccess[strid]
		if !exist {
			return nil, errors.New("access with id " + strid + " not found")
		}
		accesses = append(accesses, &access)
	}
	return accesses, err
}

// *************** Suppliers Query/Mutations ************************** //
func (r *queryResolver) Suppliers(ctx context.Context) ([]*model.Supplier, error) {
	suppliers, err := redis.GetClient().AllSuppliers()
	return suppliers, err
}

func (r *mutationResolver) CreateSupplier(ctx context.Context, input *SupplierInput) (*model.Supplier, error) {
	id := randx.String(8)
	s := model.Supplier{
		ID:   id,
		Name: input.Name,
		Code: input.Code,
	}

	err := redis.GetClient().SaveSupplier(s)
	if err != nil {
		return nil, err
	}
	sup, _ := json.Marshal(s)
	redis.GetClient().Publish("suppliers", sup)

	return &s, nil
}

func (r *mutationResolver) UpdateSupplier(ctx context.Context, id string, input *SupplierInput) (*model.Supplier, error) {
	s := model.Supplier{
		ID:   id,
		Name: input.Name,
		Code: input.Code,
	}

	err := redis.GetClient().SaveSupplier(s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
