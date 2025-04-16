package executors

import "fmt"

type PrintExec struct {
}

func NewPrintExec() *PrintExec {
	return &PrintExec{}
}

func (p *PrintExec) DoChangeWithNewValues(configMap map[string]map[string]interface{}, valuesNames string) error {
	execStr := "run app with "
	for k, v := range configMap[valuesNames] {
		execStr += fmt.Sprintf(" -e %s=%s", k, v)
	}

	fmt.Println(execStr)
	return nil
}
