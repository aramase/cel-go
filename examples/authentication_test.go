package examples

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

func ExampleNestedFunction() {
	env, err := cel.NewEnv(
		// dynType should only really be bytes or string...
		// https://github.com/google/cel-spec/blob/master/doc/langdef.md#string-and-bytes-values
		cel.Variable("claims", cel.DynType),
		ext.Encoders(),
		cel.Function("json",
			cel.Overload("json", []*cel.Type{cel.DynType}, cel.DynType,
				cel.UnaryBinding(
					func(val ref.Val) ref.Val {
						val = val.ConvertToType(types.BytesType)
						data, ok := val.(types.Bytes)
						if !ok {
							return types.MaybeNoSuchOverloadErr(val)
						}
						m := make(map[string]any)
						err := json.Unmarshal(data, &m)
						if err != nil {
							return types.NewErr("json.Unmarshal %v", err.Error())
						}
						return types.DefaultTypeAdapter.NativeToValue(m)
					},
				),
			),
		),
	)

	if err != nil {
		panic(err)
	}

	expr := `json(claims).custom.data.name`
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		panic(issues.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		panic(err)
	}

	out, _, err := prg.Eval(map[string]any{
		"claims": `{"custom":{"data":{"name":"foo"}}}`,
	})

	if err != nil {
		panic(err)
	}

	o := fmt.Sprintf("%v", out)
	log.Print(o)
	fmt.Println(o)
	// Output: foo
}
