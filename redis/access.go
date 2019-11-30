package redis

import (
	"encoding/json"
	"labxsubs/model"
)

func (c RedisClient) SaveAccess(m model.Access) error {
	serialized_access, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = c.Set("a:"+m.ID, string(serialized_access), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c RedisClient) AllAccesses() ([]*model.Access, error) {
	var ret []*model.Access
	keys := c.Keys("a:*")
	keyList := keys.Val()

	for _, k := range keyList {
		aux := c.Get(k)
		var access *model.Access
		json.Unmarshal([]byte(aux.Val()), &access)
		ret = append(ret, access)
	}
	return ret, nil
}
