package models

import "github.com/romangrechin/anviz-rpc/anviz/errors"

type Capacity struct {
	Users        int32
	Fingerprints int32
	Records      int32
}

func (c *Capacity) Unmarshal(data []byte) error {
	if len(data) != 9 {
		return errors.ErrInvalidCapacityData
	}

	c.Users = toInt(data[0:3])
	c.Fingerprints = toInt(data[3:6])
	c.Records = toInt(data[6:9])
	return nil
}
