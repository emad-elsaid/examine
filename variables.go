package examine

import (
	"fmt"
	"reflect"

	"github.com/go-delve/delve/service/api"
)

type Variables []api.Variable

func (vs Variables) Name(s string) *Variable {
	for _, v := range vs {
		if v.Name == s {
			return &Variable{val: &v}
		}
	}

	return nil
}

type Variable struct {
	val *api.Variable
	err error
}

func (v *Variable) String() string {
	if v.err != nil {
		return ""
	}

	switch v.val.Kind {
	case reflect.String:
		return v.val.Value
	default:
		return v.val.SinglelineString()
	}
}

func (v *Variable) Error() string {
	if v.err == nil {
		return ""
	}

	return v.err.Error()
}

func (v *Variable) Dereference() *Variable {
	if v.err != nil {
		return v
	}

	if v.val.Kind != reflect.Ptr {
		v.err = fmt.Errorf("Variable is not pointer")
		return v
	}

	if len(v.val.Children) == 0 {
		v.err = fmt.Errorf("Pointer doesn't have children")
		return v
	}

	return &Variable{val: &v.val.Children[0]}
}

func (v *Variable) Field(name string) *Variable {
	if v.err != nil {
		return v
	}

	for _, c := range v.val.Children {
		if c.Name == name {
			return &Variable{val: &c}
		}
	}

	fields := make([]string, 0, len(v.val.Children))
	for _, c := range v.val.Children {
		fields = append(fields, c.Name)
	}

	v.err = fmt.Errorf("Field can't be found: %s, Available fields: %s, Var: %v", name, fields, v.val.Kind)

	return v
}
