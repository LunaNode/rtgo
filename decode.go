package rtgo

import "fmt"
import "reflect"
import "strconv"
import "strings"
import "time"

func UnmarshalShort(dataStr string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("cannot unmarshal into type %v, not a pointer", reflect.TypeOf(v))
	}
	el := rv.Elem()
	if el.Kind() != reflect.Struct {
		return fmt.Errorf("cannot unmarshal into non-struct pointer")
	}
	elType := el.Type()

	// map from names to values that we should assign
	values := make(map[string]reflect.Value)

	for i := 0; i < el.NumField(); i++ {
		structField := elType.Field(i)
		name := structField.Tag.Get("rt")
		if name == "" {
			name = structField.Name
		}
		values[name] = el.Field(i)
	}

	// convert data string to string key-value map
	data := make(map[string]string)
	var lastKey string
	for _, line := range strings.Split(dataStr, "\n") {
		if len(line) == 0 {
			continue
		} else if line[0] == ' ' && lastKey != "" {
			// multiline values are prefixed with a space
			data[lastKey] += "\n" + strings.TrimSpace(line)
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			return fmt.Errorf("data has line without colon: %s", line)
		}

		lastKey = parts[0]

		if len(parts[1]) > 0 && parts[1][0] == ' ' {
			data[lastKey] = parts[1][1:]
		} else {
			data[lastKey] = parts[1]
		}
	}

	// convert values to struct field
	for key, value := range data {
		if value == "" || value == "Not set" {
			continue
		}

		rv, ok := values[key]
		if !ok {
			continue
		}
		rvType := rv.Type()

		if rvType.PkgPath() == "time" && rvType.Name() == "Time" {
			t, err := time.Parse("Mon Jan 2 15:04:05 2006", value)
			if err == nil {
				rv.Set(reflect.ValueOf(t))
			} else {
				return fmt.Errorf("failed to decode %s as time: %v", value, err)
			}
		} else if rv.Kind() == reflect.String {
			rv.SetString(value)
		} else if rv.Kind() == reflect.Int {
			x, err := strconv.Atoi(value)
			if err == nil {
				rv.SetInt(int64(x))
			} else {
				return fmt.Errorf("failed to decode %s as int: %v", value, err)
			}
		}
	}

	return nil
}
