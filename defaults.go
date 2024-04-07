package examine

import (
	"github.com/go-delve/delve/service/api"
)

var loadArgs = api.LoadConfig{
	FollowPointers:     false,
	MaxArrayValues:     0,
	MaxStringLen:       0,
	MaxStructFields:    0,
	MaxVariableRecurse: 0,
}
