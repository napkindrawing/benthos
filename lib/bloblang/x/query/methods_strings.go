package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"format", true, formatMethod,
)

func formatMethod(target Function, args ...interface{}) (Function, error) {
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, err
		}
		switch t := v.(type) {
		case string:
			return fmt.Sprintf(t, args...), nil
		case []byte:
			return fmt.Sprintf(string(t), args...), nil
		default:
			return nil, fmt.Errorf("expected string value, received %T", v)
		}
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"lowercase", false, lowercaseMethod,
	ExpectNArgs(0),
)

func lowercaseMethod(target Function, _ ...interface{}) (Function, error) {
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, &ErrRecoverable{
				Recovered: "",
				Err:       err,
			}
		}
		switch t := v.(type) {
		case string:
			return strings.ToLower(t), nil
		case []byte:
			return bytes.ToLower(t), nil
		default:
			return nil, &ErrRecoverable{
				Recovered: strings.ToLower(IToString(v)),
				Err:       fmt.Errorf("expected string value, received %T", v),
			}
		}
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"parse_json", false, parseJSONMethod,
	ExpectNArgs(0),
)

func parseJSONMethod(target Function, _ ...interface{}) (Function, error) {
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, err
		}
		var jsonBytes []byte
		switch t := v.(type) {
		case string:
			jsonBytes = []byte(t)
		case []byte:
			jsonBytes = t
		default:
			return nil, fmt.Errorf("expected string value, received %T", v)
		}
		var jObj interface{}
		if err = json.Unmarshal(jsonBytes, &jObj); err != nil {
			return nil, fmt.Errorf("failed to parse value as JSON: %w", err)
		}
		return jObj, nil
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"re_match", true, regexpMatchMethod,
	ExpectNArgs(1),
	ExpectStringArg(0),
)

func regexpMatchMethod(target Function, args ...interface{}) (Function, error) {
	re, err := regexp.Compile(args[0].(string))
	if err != nil {
		return nil, err
	}
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, err
		}
		var result bool
		switch t := v.(type) {
		case string:
			result = re.MatchString(t)
		case []byte:
			result = re.Match(t)
		default:
			return nil, fmt.Errorf("expected string value, received %T", v)
		}
		return result, nil
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"re_replace", true, regexpReplaceMethod,
	ExpectNArgs(2),
	ExpectStringArg(0),
	ExpectStringArg(1),
)

func regexpReplaceMethod(target Function, args ...interface{}) (Function, error) {
	re, err := regexp.Compile(args[0].(string))
	if err != nil {
		return nil, err
	}
	with := args[1].(string)
	withBytes := []byte(with)
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, err
		}
		var result string
		switch t := v.(type) {
		case string:
			result = re.ReplaceAllString(t, with)
		case []byte:
			result = string(re.ReplaceAll(t, withBytes))
		default:
			return nil, fmt.Errorf("expected string value, received %T", v)
		}
		return result, nil
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"slice", true, sliceMethod,
	ExpectAtLeastOneArg(),
	ExpectIntArg(0),
	ExpectIntArg(1),
)

func sliceMethod(target Function, args ...interface{}) (Function, error) {
	low := args[0].(int64)
	var high *int64
	if len(args) > 1 {
		highV := args[1].(int64)
		high = &highV
	}
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, err
		}
		switch t := v.(type) {
		case string:
			highV := int64(len(t))
			if high != nil {
				highV = *high
			}
			return t[low:highV], nil
		case []byte:
			highV := int64(len(t))
			if high != nil {
				highV = *high
			}
			return t[low:highV], nil
		}
		return nil, fmt.Errorf("expected string value, received %T", v)
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"string", false, stringMethod,
	ExpectNArgs(0),
)

func stringMethod(target Function, _ ...interface{}) (Function, error) {
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, &ErrRecoverable{
				Recovered: "",
				Err:       err,
			}
		}
		return IToString(v), nil
	}), nil
}

//------------------------------------------------------------------------------

var _ = RegisterMethod(
	"uppercase", false, uppercaseMethod,
	ExpectNArgs(0),
)

func uppercaseMethod(target Function, _ ...interface{}) (Function, error) {
	return closureFn(func(ctx FunctionContext) (interface{}, error) {
		v, err := target.Exec(ctx)
		if err != nil {
			return nil, &ErrRecoverable{
				Recovered: "",
				Err:       err,
			}
		}
		switch t := v.(type) {
		case string:
			return strings.ToUpper(t), nil
		case []byte:
			return bytes.ToUpper(t), nil
		default:
			return nil, &ErrRecoverable{
				Recovered: strings.ToUpper(IToString(v)),
				Err:       fmt.Errorf("expected string value, received %T", v),
			}
		}
	}), nil
}

//------------------------------------------------------------------------------