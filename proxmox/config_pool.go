package proxmox

import (
	"errors"
	"regexp"
)

func ListPools(c *Client) ([]PoolName, error) {
	raw, err := listPools(c)
	if err != nil {
		return nil, err
	}
	pools := make([]PoolName, len(raw))
	for i, e := range raw {
		pools[i] = PoolName(e.(map[string]interface{})["poolid"].(string))
	}
	return pools, nil
}

func listPools(c *Client) ([]interface{}, error) {
	return c.GetItemListInterfaceArray("/pools")
}

type ConfigPool struct {
	Name    PoolName `json:"name"`
	Comment *string  `json:"comment"`
	Guests  *[]uint  `json:"guests"` // TODO: Change type once we have a type for guestID
}

func (ConfigPool) mapToSDK(params map[string]interface{}) (config ConfigPool) {
	if v, isSet := params["poolid"]; isSet {
		config.Name = PoolName(v.(string))
	}
	if v, isSet := params["comment"]; isSet {
		tmp := v.(string)
		config.Comment = &tmp
	}
	if v, isSet := params["members"]; isSet {
		guests := make([]uint, 0)
		for _, e := range v.([]interface{}) {
			param := e.(map[string]interface{})
			if v, isSet := param["vmid"]; isSet {
				guests = append(guests, uint(v.(float64)))
			}
		}
		if len(guests) > 0 {
			config.Guests = &guests
		}
	}
	return
}

// Same as PoolName.Delete()
func (config ConfigPool) Delete(c *Client) error {
	return config.Name.Delete(c)
}

// Same as PoolName.Exists()
func (config ConfigPool) Exists(c *Client) (bool, error) {
	return config.Name.Exists(c)
}

func (config ConfigPool) Validate() error {
	// TODO: Add validation for Guests and Comment
	return config.Name.Validate()
}

type PoolName string

const (
	PoolName_Error_Characters string = "PoolName may only contain the following characters: a-z, A-Z, 0-9, hyphen (-), and underscore (_)"
	PoolName_Error_Empty      string = "PoolName cannot be empty"
	PoolName_Error_Length     string = "PoolName may not be longer than 1024 characters" // proxmox does not seem to have a max length, so we artificially cap it at 1024
	PoolName_Error_NotExists  string = "Pool doesn't exist"
)

var regex_PoolName = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

func (config PoolName) Delete(c *Client) error {
	if c == nil {
		return errors.New(Client_Error_Nil)
	}
	if err := config.Validate(); err != nil {
		return err
	}
	// TODO: permission check
	exists, err := config.Exists_Unsafe(c)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(PoolName_Error_NotExists)
	}
	return config.Delete_Unsafe(c)
}

func (config PoolName) Delete_Unsafe(c *Client) error {
	return c.Delete("/pools/" + string(config))
}

func (config PoolName) Exists(c *Client) (bool, error) {
	if c == nil {
		return false, errors.New(Client_Error_Nil)
	}
	if err := config.Validate(); err != nil {
		return false, err
	}
	// TODO: permission check
	return config.Exists_Unsafe(c)
}

func (config PoolName) Exists_Unsafe(c *Client) (bool, error) {
	raw, err := listPools(c)
	if err != nil {
		return false, err
	}
	return ItemInKeyOfArray(raw, "poolid", string(config)), nil
}

func (pool PoolName) Get(c *Client) (*ConfigPool, error) {
	if c == nil {
		return nil, errors.New(Client_Error_Nil)
	}
	if err := pool.Validate(); err != nil {
		return nil, err
	}
	// TODO: permission check
	return pool.Get_Unsafe(c)
}

func (pool PoolName) Get_Unsafe(c *Client) (*ConfigPool, error) {
	params, err := c.GetItemConfigMapStringInterface("/pools/"+string(pool), "pool", "CONFIG")
	if err != nil {
		return nil, err
	}
	config := ConfigPool{}.mapToSDK(params)
	return &config, nil
}

func (config PoolName) Validate() error {
	if config == "" {
		return errors.New(PoolName_Error_Empty)
	}
	if len(config) > 1024 {
		return errors.New(PoolName_Error_Length)
	}
	if !regex_PoolName.MatchString(string(config)) {
		return errors.New(PoolName_Error_Characters)
	}
	return nil
}
