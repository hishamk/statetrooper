package statetrooper

import "fmt"

func stringable(t interface{}) bool {
	if _, ok := t.(string); ok {
		return true
	}

	if _, ok := t.(fmt.Stringer); ok {
		return true
	}

	return false
}

// function to convert any type to a string
func toString(t interface{}) string {
	if str, ok := t.(string); ok {
		return str
	}

	if strer, ok := t.(fmt.Stringer); ok {
		return strer.String()
	}

	return fmt.Sprintf("%v", t)
}
