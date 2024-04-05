package examine

import (
	"github.com/go-delve/delve/pkg/proc"
	"github.com/go-delve/delve/service/api"
)

var loadArgs = api.LoadConfig{
	FollowPointers:     true,
	MaxArrayValues:     10,
	MaxStringLen:       100,
	MaxStructFields:    10,
	MaxVariableRecurse: 5,
}

var procLoadArgs = proc.LoadConfig{
	FollowPointers:     true,
	MaxArrayValues:     10,
	MaxStringLen:       100,
	MaxStructFields:    10,
	MaxVariableRecurse: 5,
}
