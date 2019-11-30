package labx_subs

import (
	"context"
	"fmt"
	"labxsubs/model"
	"labxsubs/redis"
	"sort"
	"strconv"

	"github.com/markbates/going/randx"
)

type accessResolver struct{ *Resolver }

func (r *accessResolver) Supplier(ctx context.Context, obj *model.Access) (*model.Supplier, error) {
	suppliers, err := redis.GetClient().AllSuppliers()

	exist := false
	var supplier *model.Supplier

	for _, sup := range suppliers {
		if sup.ID == obj.SupplierID {
			exist = true
			supplier = sup
		}
	}
	if !exist {
		return nil, fmt.Errorf("supplier not found")
	}

	return supplier, err
}

// *************** Accesses Query/Mutations ************************** //
func (r *queryResolver) Accesses(ctx context.Context) ([]*model.Access, error) {
	res, err := redis.GetClient().AllAccesses()
	sort.Slice(res, func(i, j int) bool {
		id1, _ := strconv.Atoi(res[i].ID)
		id2, _ := strconv.Atoi(res[j].ID)
		return id1 < id2
	})
	return res, err
}

func (r *mutationResolver) CreateAccess(ctx context.Context, input *AccessInput) (*model.Access, error) {
	// validations
	suppliers, err := redis.GetClient().AllSuppliers()
	id := randx.String(8)
	exist := false
	for _, sup := range suppliers {
		if sup.ID == input.Supplier {
			exist = true
			sup.AccessesID = append(sup.AccessesID, id)
			redis.GetClient().SaveSupplier(*sup)
		}
	}
	if !exist {
		return nil, fmt.Errorf("cannot create an access for a supplier that does not exists")
	}

	a := model.Access{
		ID:         id,
		Name:       input.Name,
		URL:        input.URL,
		User:       input.User,
		SupplierID: input.Supplier,
	}
	redis.GetClient().SaveAccess(a)

	return &a, err
}

func (r *mutationResolver) UpdateAccess(ctx context.Context, id string, input *AccessInput) (*model.Access, error) {
	acceses, err := redis.GetClient().AllAccesses()
	var ret *model.Access

	exist := false
	for _, value := range acceses {
		if value.ID == id {
			exist = true
			ret := value

			ret.Name = input.Name
			ret.SupplierID = input.Supplier
			ret.URL = input.URL
			ret.User = input.User

			redis.GetClient().SaveAccess(*ret)
		}
	}
	if !exist {
		return nil, fmt.Errorf("access not found")
	}
	return ret, err
}
