package solenodon

// see: https://github.com/Jeffail/gabs/blob/master/gabs.go

// Note that encoding/json by default will parse:
// - all number values into float64
// Note that github.com/BurntSushi/toml by default will parse:
// - all integer values into int64
// Note that encoding/xml cannot be mapped to an interface{}

// Container contains data
type Container struct {
	data   interface{}
	parent *Container
	key    interface{}
}

// NewContainer returns a new Container for the given data
func NewContainer(data interface{}) *Container {
	return &Container{data: data}
}

// NewContainerFromBytes returns a new Container using the given bytes and unmarshal function
func NewContainerFromBytes(b []byte, unmarshal func(b []byte, target interface{}) error) (*Container, error) {
	var data interface{}
	err := unmarshal(b, &data)
	return NewContainer(data), err
}

// NewContainerFromDecoder returns a new Container using the given decode function
func NewContainerFromDecoder(decode func(v interface{}) error) (*Container, error) {
	var data interface{}
	err := decode(&data)
	return NewContainer(data), err
}

// Data returns the data in the Container
func (c *Container) Data() interface{} {
	if c == nil {
		return nil
	}
	return c.data
}

// Get returns a Container containing the value following the path of the given keys
// The returned container will be nil if no result was found
func (c *Container) Get(keys ...interface{}) *Container {
	if c == nil {
		return c
	}
	result := c
	for _, key := range keys {
		var data interface{}
		switch w := result.data.(type) {
		case map[string]interface{}:
			if v, ok := key.(string); ok {
				data, ok = w[v]
				if !ok {
					return nil
				}
			} else {
				return nil
			}
		case map[interface{}]interface{}:
			var ok bool
			data, ok = w[key]
			if !ok {
				return nil
			}
		case []interface{}:
			if v, ok := key.(int); ok {
				if v < 0 || v >= len(w) {
					return nil
				}
				data = w[v]
			} else {
				return nil
			}
		case []map[string]interface{}:
			if v, ok := key.(int); ok {
				if v < 0 || v >= len(w) {
					return nil
				}
				data = w[v]
			}
		default:
			return nil
		}
		result = &Container{
			data:   data,
			parent: result,
			key:    key,
		}
	}
	return result
}

// Has returns true if the Container has a value for the given keys
func (c *Container) Has(keys ...interface{}) bool {
	return c.Get(keys...) != nil
}

// Delete the value at the end of the path of the given keys
func (c *Container) Delete(keys ...interface{}) *Container {
	if c == nil {
		return c
	}
	if len(keys) == 0 {
		c.data = nil
		return c
	}
	parent := c.Get(keys[:len(keys)-1]...)
	if parent == nil {
		return c
	}
	lastKey := keys[len(keys)-1]

	switch w := parent.data.(type) {
	case map[string]interface{}:
		if v, ok := lastKey.(string); ok {
			delete(w, v)
		}
	case map[interface{}]interface{}:
		delete(w, lastKey)
	case []interface{}:
		if v, ok := lastKey.(int); ok && v >= 0 && v < len(w) {
			parent.Replace(append(w[:v], w[v+1:]...))
		}
	case []map[string]interface{}:
		if v, ok := lastKey.(int); ok && v >= 0 && v < len(w) {
			parent.Replace(append(w[:v], w[v+1:]...))
		}
	}
	return c
}

// Replace the data
// If the container has a parent, the parent will reference the replacement
func (c *Container) Replace(with interface{}) *Container {
	// TODO panic if parent does not contains key?
	if c == nil {
		return c
	}
	if c.parent == nil {
		c.data = with
		return c
	}
	switch w := c.parent.data.(type) {
	case map[string]interface{}:
		if v, ok := c.key.(string); ok {
			if _, ok := w[v]; ok {
				w[v] = with
			} else {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		if _, ok := w[c.key]; ok {
			w[c.key] = with
		} else {
			return nil
		}
	case []interface{}:
		if v, ok := c.key.(int); ok && v >= 0 && v < len(w) {
			w[v] = with
		} else {
			return nil
		}
	case []map[string]interface{}:
		if v, ok := c.key.(int); ok && v >= 0 && v < len(w) {
			parentReplacement := make([]interface{}, len(w))
			for i, x := range w {
				parentReplacement[i] = x
			}
			parentReplacement[v] = with
			c.parent.Replace(parentReplacement)
		} else {
			return nil
		}
	default:
		return nil
	}
	c.data = with
	return c
}
