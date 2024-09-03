package test

import (
	"fmt"
	"testing"

	"github.com/expr-lang/expr"
)

type CpuPerfData struct {
	User      float64 `expr:"user"`
	System    float64 `expr:"system"`
	Idle      float64 `expr:"idle"`
	Nice      float64 `expr:"nice"`
	Iowait    float64 `expr:"iowait"`
	Irq       float64 `expr:"irq"`
	Softirq   float64 `expr:"softirq"`
	Steal     float64 `expr:"steal"`
	Guest     float64 `expr:"guest"`
	GuestNice float64 `expr:"guestNice"`
}

type MemPerfData struct {
	Used uint64 `expr:"used"`
}

func TestGovaluate(t *testing.T) {
	env0 := map[string]interface{}{
		"cpu":    CpuPerfData{},
		"memory": MemPerfData{},
	}

	env1 := map[string]interface{}{
		"cpu": CpuPerfData{
			User: 10,
		},
	}

	code := `cpu.user > 5`

	program, err := expr.Compile(code, expr.Env(env0), expr.AsBool())
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env1)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}
