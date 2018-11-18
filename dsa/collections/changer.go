// Tideland Go Library - DSA - Collections - Changer
//
// Copyright (C) 2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collections

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.one/go/trace/errors"
)

//--------------------
// CHANGER
//--------------------

// changer implements the Changer interface.
type changer struct {
	node *node
	err  error
}

// Value implements the Changer interface.
func (c *changer) Value() (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.node.content.value(), nil
}

// SetValue implements the Changer interface.
func (c *changer) SetValue(v interface{}) (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	oldValue := c.node.content.value()
	newValue := justValue{v}
	if !c.node.isAllowed(newValue, true) {
		return nil, errors.New(ErrDuplicate, "setting duplicate value is not allowed")
	}
	c.node.content = newValue
	return oldValue, nil
}

// Add implements the Changer interface.
func (c *changer) Add(v interface{}) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(justValue{v})
	return err
}

// Remove implements the Changer interface.
func (c *changer) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List implements the Changer interface.
func (c *changer) List() ([]interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []interface{}
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, cn.content.value())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error implements the Changer interface.
func (c *changer) Error() error {
	return c.err
}

//--------------------
// STRING CHANGER
//--------------------

// stringChanger implements the StringChanger interface.
type stringChanger struct {
	node *node
	err  error
}

// Value implements the StringChanger interface.
func (c *stringChanger) Value() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.node.content.value() == nil {
		return "", nil
	}
	return c.node.content.value().(string), nil
}

// SetValue implements the StringChanger interface.
func (c *stringChanger) SetValue(v string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	oldValue := c.node.content.value().(string)
	newValue := justValue{v}
	if !c.node.isAllowed(newValue, true) {
		return "", errors.New(ErrDuplicate, "setting duplicate string value is not allowed")
	}
	c.node.content = newValue
	return oldValue, nil
}

// Add implements the StringChanger interface.
func (c *stringChanger) Add(v string) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(justValue{v})
	return err
}

// Remove implements the StringChanger interface.
func (c *stringChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List implements the StringChanger interface.
func (c *stringChanger) List() ([]string, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []string
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, cn.content.value().(string))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error implements the StringChanger interface.
func (c *stringChanger) Error() error {
	return c.err
}

//--------------------
// KEY/VALUE CHANGER
//--------------------

// keyValueChanger implements the KeyValueChanger interface.
type keyValueChanger struct {
	node *node
	err  error
}

// Key implements the KeyValueChanger interface.
func (c *keyValueChanger) Key() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	return c.node.content.key().(string), nil
}

// SetKey implements the KeyValueChanger interface.
func (c *keyValueChanger) SetKey(key string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if !c.node.container.duplicates {
		if c.node.hasDuplicateSibling(key) {
			return "", errors.New(ErrDuplicate, "setting duplicate key is not allowed")
		}
	}
	current := c.node.content.key().(string)
	c.node.content = keyValue{key, c.node.content.value()}
	return current, nil
}

// Value implements the KeyValueChanger interface.
func (c *keyValueChanger) Value() (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.node.content.value(), nil
}

// SetValue implements the KeyValueChanger interface.
func (c *keyValueChanger) SetValue(value interface{}) (interface{}, error) {
	if c.err != nil {
		return nil, c.err
	}
	current := c.node.content.value()
	c.node.content = keyValue{c.node.content.key(), value}
	return current, nil
}

// Add implements the KeyValueChanger interface.
func (c *keyValueChanger) Add(k string, v interface{}) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(keyValue{k, v})
	return err
}

// Remove implements the KeyValueChanger interface.
func (c *keyValueChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List implements the KeyValueChanger interface.
func (c *keyValueChanger) List() ([]KeyValue, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []KeyValue
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, KeyValue{cn.content.key().(string), cn.content.value()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error implements the KeyValueChanger interface.
func (c *keyValueChanger) Error() error {
	return c.err
}

//--------------------
// KEY/STRING VALUE CHANGER
//--------------------

// keyStringValueChanger implements the KeyStringValueChanger interface.
type keyStringValueChanger struct {
	node *node
	err  error
}

// Key implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) Key() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	return c.node.content.key().(string), nil
}

// SetKey implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) SetKey(key string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if !c.node.container.duplicates {
		if c.node.hasDuplicateSibling(key) {
			return "", errors.New(ErrDuplicate, "setting duplicate key is not allowed")
		}
	}
	current := c.node.content.key().(string)
	c.node.content = keyValue{key, c.node.content.value()}
	return current, nil
}

// Value implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) Value() (string, error) {
	if c.err != nil {
		return "", c.err
	}
	if c.node.content.value() == nil {
		return "", nil
	}
	return c.node.content.value().(string), nil
}

// SetValue implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) SetValue(value string) (string, error) {
	if c.err != nil {
		return "", c.err
	}
	current := c.node.content.value().(string)
	c.node.content = keyValue{c.node.content.key(), value}
	return current, nil
}

// Add implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) Add(k, v string) error {
	if c.err != nil {
		return c.err
	}
	_, err := c.node.addChild(keyValue{k, v})
	return err
}

// Remove implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) Remove() error {
	if c.err != nil {
		return c.err
	}
	return c.node.remove()
}

// List implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) List() ([]KeyStringValue, error) {
	if c.err != nil {
		return nil, c.err
	}
	var list []KeyStringValue
	err := c.node.doChildren(func(cn *node) error {
		list = append(list, KeyStringValue{cn.content.key().(string), cn.content.value().(string)})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Error implements the KeyStringValueChanger interface.
func (c *keyStringValueChanger) Error() error {
	return c.err
}

// EOF
