package solenodon

type c interface {
	Get(keys ...interface{}) *Container
	Set(keys []interface{}, value interface{})
	Del(keys ...interface{})
	Search(value interface{})
	Replace(with interface{}) *Container
}

// see: https://github.com/Jeffail/gabs/blob/master/gabs.go

// Note that encoding/json by default will parse:
// - all number values into float64
// Note that github.com/BurntSushi/toml by default will parse:
// - all integer values into int64
// Note that encoding/xml cannot be mapped to an interface{}

// Container contains data
type Container struct {
	Data interface{}

	parent interface{}
	key    interface{}
}

// Search for a value following the keys
// keys must be of type:
// - string
// - int
// The returned container will be nil if no result was found
func (c *Container) Search(keys ...interface{}) *Container {
	if c == nil {
		return nil
	}
	if len(keys) == 0 {
		return c
	}
	data := c.Data
	var parent, key interface{}
	for _, key = range keys {
		parent = data
		switch w := data.(type) {
		case map[string]interface{}:
			switch v := key.(type) {
			case string:
				var ok bool
				data, ok = w[v]
				if !ok {
					return nil
				}
			default:
				return nil
			}
		case map[interface{}]interface{}:
			var ok bool
			data, ok = w[key]
			if !ok {
				return nil
			}
		case []interface{}:
			switch v := key.(type) {
			case int:
				if v >= len(w) {
					return nil
				}
				data = w[v]
			default:
				return nil
			}
		default:
			return nil
		}
	}
	return &Container{
		Data:   data,
		parent: parent,
		key:    key,
	}
}

// Replace the data
// If the container has a parent, the parent will reference the replacement
func (c *Container) Replace(with interface{}) *Container {
	if c == nil {
		return nil
	}
	if c.parent == nil {
		c.Data = with
		return c
	}
	switch w := c.parent.(type) {
	case map[string]interface{}:
		switch v := c.key.(type) {
		case string:
			if _, ok := w[v]; ok {
				w[v] = with
			} else {
				return nil
			}
		default:
			return nil
		}
	case map[interface{}]interface{}:
		if _, ok := w[c.key]; ok {
			w[c.key] = with
		} else {
			return nil
		}
	case []interface{}:
		switch v := c.key.(type) {
		case int:
			if v >= len(w) {
				return nil
			}
			w[v] = with
		default:
			return nil
		}
	default:
		return nil
	}
	return &Container{
		Data:   with,
		parent: c.parent,
		key:    c.key,
	}
}
