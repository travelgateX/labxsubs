package redis

import (
	"encoding/json"
	"labxsubs/model"
)

func (c RedisClient) SaveSupplier(s model.Supplier) error {
	serialized_supplier, err := json.Marshal(s)
	if err != nil {
		return err
	}

	err = c.Set("s:"+s.ID, string(serialized_supplier), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c RedisClient) AllSuppliers() ([]*model.Supplier, error) {
	var ret []*model.Supplier
	keys := c.Keys("s:*")
	keylist := keys.Val()
	for _, k := range keylist {
		aux := c.Get(k)
		var supplier *model.Supplier
		json.Unmarshal([]byte(aux.Val()), &supplier)
		ret = append(ret, supplier)
	}
	return ret, nil
}
