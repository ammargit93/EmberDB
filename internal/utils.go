package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func StringifyValue(v Value) string {
	encoded := base64.StdEncoding.EncodeToString(v.Data)
	return fmt.Sprintf("%d:%s", v.Type, encoded)
}

func ParseValue(t string, v interface{}) (Value, error) {
	switch t {
	case "string":
		s, ok := v.(string)
		if !ok {
			return Value{}, errors.New("value must be string")
		}
		return StringValue(s), nil

	case "int":
		// JSON numbers decode as float64
		f, ok := v.(float64)
		if !ok {
			return Value{}, errors.New("value must be number")
		}
		return IntValue(int64(f)), nil

	case "float":
		f, ok := v.(float64)
		if !ok {
			return Value{}, errors.New("value must be number")
		}
		return FloatValue(f), nil

	case "bool":
		b, ok := v.(bool)
		if !ok {
			return Value{}, errors.New("value must be bool")
		}
		return BoolValue(b), nil

	case "file":
		s, ok := v.(string)
		if !ok {
			return Value{}, errors.New("file must be base64 string")
		}
		data, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return Value{}, err
		}
		return FileValue(data), nil

	default:
		return Value{}, errors.New("unsupported type")
	}
}
