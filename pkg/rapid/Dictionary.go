// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

type Dictionary struct {
	keys        []interface{}
	boolKeys    map[bool]int
	stringKeys  map[string]int
	intKeys     map[int]int
	float64Keys map[float64]int
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		keys:        []interface{}{},
		boolKeys:    map[bool]int{},
		stringKeys:  map[string]int{},
		intKeys:     map[int]int{},
		float64Keys: map[float64]int{},
	}
}

func (d *Dictionary) Cap() int {
	return cap(d.keys)
}

func (d *Dictionary) Reset() {
	d.keys = d.keys[:0]
	//d.boolKeys = map[bool]int{}
	for k := range d.boolKeys {
		delete(d.boolKeys, k)
	}
	//d.stringKeys = map[string]int{}
	for k := range d.stringKeys {
		delete(d.stringKeys, k)
	}
	//d.intKeys = map[int]int{}
	for k := range d.intKeys {
		delete(d.intKeys, k)
	}
	//d.float64Keys = map[float64]int{}
	for k := range d.float64Keys {
		delete(d.float64Keys, k)
	}
}

func (d *Dictionary) Keys() []interface{} {
	return d.keys
}

func (d *Dictionary) addBoolKey(key bool) int {
	d.keys = append(d.keys, key)
	i := len(d.keys) - 1
	d.boolKeys[key] = i
	return i
}

func (d *Dictionary) addStringKey(key string) int {
	d.keys = append(d.keys, key)
	i := len(d.keys) - 1
	d.stringKeys[key] = i
	return i
}

func (d *Dictionary) addIntKey(key int) int {
	d.keys = append(d.keys, key)
	i := len(d.keys) - 1
	d.intKeys[key] = i
	return i
}

func (d *Dictionary) addFloat64Key(key float64) int {
	d.keys = append(d.keys, key)
	i := len(d.keys) - 1
	d.float64Keys[key] = i
	return i
}

func (d *Dictionary) addArrayKey(key []interface{}) int {
	d.keys = append(d.keys, key)
	i := len(d.keys) - 1
	// array keys are not indexed
	return i
}

func (d *Dictionary) AddKey(key interface{}) int {
	if b, ok := key.(bool); ok {
		return d.addBoolKey(b)
	}
	if str, ok := key.(string); ok {
		return d.addStringKey(str)
	}
	if i, ok := key.(int); ok {
		return d.addIntKey(i)
	}
	if f, ok := key.(float64); ok {
		return d.addFloat64Key(f)
	}
	if arr, ok := key.([]interface{}); ok {
		return d.addArrayKey(arr)
	}
	return -1
}

func (d *Dictionary) AddChain(chain ...interface{}) []int {
	indicies := make([]int, 0, len(chain))
	for _, v := range chain {
		if i := d.GetIndex(v); i != -1 {
			indicies = append(indicies, i)
		} else {
			i := d.AddKey(v)
			indicies = append(indicies, i)
		}
	}
	return indicies
}

func (d *Dictionary) hasBoolKey(key bool) bool {
	_, ok := d.boolKeys[key]
	return ok
}

func (d *Dictionary) hasStringKey(key string) bool {
	_, ok := d.stringKeys[key]
	return ok
}

func (d *Dictionary) hasIntKey(key int) bool {
	_, ok := d.intKeys[key]
	return ok
}

func (d *Dictionary) hasFloat64Key(key float64) bool {
	_, ok := d.float64Keys[key]
	return ok
}

func (d *Dictionary) HasKey(key interface{}) bool {
	if b, ok := key.(bool); ok {
		return d.hasBoolKey(b)
	}
	if str, ok := key.(string); ok {
		return d.hasStringKey(str)
	}
	if i, ok := key.(int); ok {
		return d.hasIntKey(i)
	}
	if f, ok := key.(float64); ok {
		return d.hasFloat64Key(f)
	}
	return false
}

func (d *Dictionary) getBoolIndex(key bool) int {
	i, ok := d.boolKeys[key]
	if !ok {
		return -1
	}
	return i
}

func (d *Dictionary) getStringIndex(key string) int {
	i, ok := d.stringKeys[key]
	if !ok {
		return -1
	}
	return i
}

func (d *Dictionary) getIntIndex(key int) int {
	i, ok := d.intKeys[key]
	if !ok {
		return -1
	}
	return i
}

func (d *Dictionary) getFloat64Index(key float64) int {
	i, ok := d.float64Keys[key]
	if !ok {
		return -1
	}
	return i
}

func (d *Dictionary) GetIndex(key interface{}) int {
	if b, ok := key.(bool); ok {
		return d.getBoolIndex(b)
	}
	if str, ok := key.(string); ok {
		return d.getStringIndex(str)
	}
	if i, ok := key.(int); ok {
		return d.getIntIndex(i)
	}
	if f, ok := key.(float64); ok {
		return d.getFloat64Index(f)
	}
	return -1
}

func (d *Dictionary) GetChain(indicies []int) ([]interface{}, bool) {
	chain := make([]interface{}, 0, len(indicies))
	for _, i := range indicies {
		if i >= len(d.keys) {
			return chain, false
		}
		chain = append(chain, d.keys[i])
	}
	return chain, true
}

func (d *Dictionary) GetKey(i int) (interface{}, bool) {
	if i >= len(d.keys) {
		return -1, false
	}
	return d.keys[i], true
}

func (d *Dictionary) Set(i int, v interface{}) bool {
	if i >= len(d.keys) {
		return false
	}
	d.keys[i] = v
	return true
}
