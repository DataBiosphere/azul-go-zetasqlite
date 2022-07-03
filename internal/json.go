package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-zetasql/types"
)

func JSONFromZetaSQLValue(v types.Value) string {
	value := jsonFromZetaSQLValue(v)
	switch v.Type().Kind() {
	case types.DATE:
		return toDateValueFromString(value)
	case types.ARRAY:
		return toArrayValueFromJSONString(value)
	case types.STRUCT:
		return toStructValueFromJSONString(value)
	}
	return value
}

func jsonFromZetaSQLValue(v types.Value) string {
	switch v.Type().Kind() {
	case types.DATE:
		return toDateValueFromInt64(v.ToInt64())
	case types.ARRAY:
		elems := []string{}
		if v.IsNull() {
			return "null"
		}
		for i := 0; i < v.NumElements(); i++ {
			elem := v.Element(i)
			elems = append(elems, JSONFromZetaSQLValue(elem))
		}
		return fmt.Sprintf("[%s]", strings.Join(elems, ","))
	case types.STRUCT:
		fields := []string{}
		structType := v.Type().AsStruct()
		for i := 0; i < v.NumFields(); i++ {
			field := v.Field(i)
			name := structType.Field(i).Name()
			val := JSONFromZetaSQLValue(field)
			fields = append(
				fields,
				fmt.Sprintf("%s:%s", strconv.Quote(name), string(val)),
			)
		}
		return fmt.Sprintf("{%s}", strings.Join(fields, ","))
	default:
		vv := v.SQLLiteral(0)
		if vv == "NULL" {
			return "null"
		}
		return vv
	}
}
