package solenodon

// see: https://github.com/Jeffail/gabs/blob/master/gabs.go

// Note that encoding/json by default will parse:
// - all number values into float64
// Note that github.com/BurntSushi/toml by default will parse:
// - all integer values into int64
// Note that encoding/xml cannot be mapped to an interface{}

// Container contains data
type Container struct {
	Data interface{}
}

// Search for a value following the keys
// keys must be of type:
// - string
// - int
func (c *Container) Search(keys ...interface{}) *Container {
	if c == nil {
		return nil
	}
	data := c.Data
	for _, key := range keys {
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
	return &Container{Data: data}
}
