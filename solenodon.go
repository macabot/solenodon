package solenodon

import (
	"errors"
)

// see: https://github.com/Jeffail/gabs/blob/master/gabs.go

// Note that encoding/json by default will parse:
// - all number values into float64
// Note that github.com/BurntSushi/toml by default will parse:
// - all integer values into int64
// Note that encoding/xml cannot be mapped to an interface{}

// List of errors
var (
	ErrNotFound = errors.New("solenodon: not found")
)

// Container contains data
type Container struct {
	Data interface{}
}

// Search for a value following the keys
// keys must be of type:
// - string
// - int
// error may be:
// - ErrNotFound
func (c *Container) Search(keys ...interface{}) (*Container, error) {
	data := c.Data
	for _, key := range keys {
		switch w := data.(type) {
		case map[string]interface{}:
			switch v := key.(type) {
			case string:
				var ok bool
				data, ok = w[v]
				if !ok {
					return nil, ErrNotFound
				}
			default:
				return nil, ErrNotFound
			}
		case map[interface{}]interface{}:
			var ok bool
			data, ok = w[key]
			if !ok {
				return nil, ErrNotFound
			}
		case []interface{}:
			switch v := key.(type) {
			case int:
				if v >= len(w) {
					return nil, ErrNotFound
				}
				data = w[v]
			default:
				return nil, ErrNotFound
			}
		default:
			return nil, ErrNotFound
		}
	}
	return &Container{Data: data}, nil
}
