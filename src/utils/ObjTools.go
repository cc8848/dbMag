package utils

import (
	"bytes"
	"encoding/gob"
)

/*
深度Copy对象
*/
func DeepCopy(dst, src interface{}) error {

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {

		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
