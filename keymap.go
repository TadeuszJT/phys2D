package phys2D

type Key int

type keyMap struct {
	keyToIndex           []int
	indexToKey, usedKeys []Key
}

func (k *keyMap) popUsedKey() Key {
	key := Key(k.usedKeys[len(k.usedKeys)-1])
	k.usedKeys = k.usedKeys[:len(k.usedKeys)-1]
	return key
}

func (k *keyMap) pushUsedKey(key Key) {
	k.usedKeys = append(k.usedKeys, key)
}

func (k *keyMap) Append() Key {
	/* Index of new entity */
	index := len(k.indexToKey)
	k.indexToKey = append(k.indexToKey, 0)

	/* Used key available */
	for len(k.usedKeys) > 0 {
		key := k.popUsedKey()

		// discard if end key
		end := Key(len(k.keyToIndex) - 1)
		if key == end {
			k.keyToIndex = k.keyToIndex[:end]
			continue
		}

		k.indexToKey[index] = key
		k.keyToIndex[key] = index
		return key
	}

	/* New key */
	key := Key(len(k.keyToIndex))
	k.keyToIndex = append(k.keyToIndex, index)
	k.indexToKey[index] = key
	return key
}

func (k *keyMap) Delete(key Key) {
	index := k.keyToIndex[key]

	/* indexToKey */
	end := len(k.indexToKey) - 1
	if index != end { // swap indexToKey
		endKey := k.indexToKey[end]
		k.indexToKey[index] = endKey
		k.keyToIndex[endKey] = index
	}
	k.indexToKey = k.indexToKey[:end]

	/* keyToIndex */
	end = len(k.keyToIndex) - 1
	if int(key) != end { // swap keyToIndex
		k.pushUsedKey(key)
		k.keyToIndex[key] = -1
	} else {
		k.keyToIndex = k.keyToIndex[:end]
	}
}
