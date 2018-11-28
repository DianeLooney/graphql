package resolver_test

import (
	"testing"

	"github.com/dianelooney/graphql/resolver"
)

type X struct {
	A string
}

func (X) B() string {
	return "B value"
}
func (X) GetC() string {
	return "C value"
}
func (X) D(args map[string]interface{}) string {
	return "D value"
}
func (X) E() (string, error) {
	return "E value", nil
}
func TestReflect(t *testing.T) {
	r := resolver.Reflect{
		Target: X{"A value"},
	}
	tests := map[string]string{
		"A": "A value",
		"B": "B value",
		"C": "C value",
		"D": "D value",
		"E": "E value",
	}
	for k, v := range tests {
		res, err := r.Resolve(k, nil)
		if err != nil {
			t.Errorf("Error returned from Reflect#Resolve(%s, nil): %v", k, err)
		}
		if res != v {
			t.Errorf("Reflect#Resolve(%s, nil) returned '%s' (%T), expected '%s' (%T)", k, res, res, v, v)
		}
	}
}
