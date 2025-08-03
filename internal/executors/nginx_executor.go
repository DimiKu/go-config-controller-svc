package executors

import (
	"fmt"
	"os/exec"
)

type NginxExec struct {
}

func NewNginxExec() *NginxExec {
	return &NginxExec{}
}

func (n *NginxExec) DoChangeWithNewValues(configMap map[string]map[string]interface{}, valuesNames string) error {
	if err := n.deleteExistContainer(); err != nil {
		return err
	}

	execStr := "docker run -d -p 8080:80 "
	for k, v := range configMap[valuesNames] {
		execStr += fmt.Sprintf(" -e %s=%s", k, v)
	}

	execStr += " my-nginx"

	cmd := exec.Command("bash", "-c", execStr)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

func (n *NginxExec) DoRestart() error {
	execStr := "docker restart my-nginx"

	cmd := exec.Command("bash", "-c", execStr)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil

}

func (n *NginxExec) deleteExistContainer() error {
	execStr := "docker ps | grep my-nginx | awk '{print$1}' | xargs docker rm  --force || true"
	cmd := exec.Command("bash", "-c", execStr)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
