package command

//
// import (
// 	"bytes"
// 	"fmt"
// 	"os/exec"
//
// 	"github.com/fatih/color"
// 	"github.com/yimincai/nbx/pkg/clog"
// )
//
// // By default, shell will be using bash
// func commandOut(command string, shell string) (string, string, error) {
// 	var stdout bytes.Buffer
// 	var stderr bytes.Buffer
//
// 	if shell == "" {
// 		shell = "bash"
// 	}
//
// 	cmd := exec.Command(shell, "-c", command)
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
// 	err := cmd.Run()
// 	return stdout.String(), stderr.String(), err
// }
//
// // Run executes a command and prints it to stdout
// func Run(command string, showCommand bool) {
// 	if showCommand {
// 		color.HiBlue("$ %v", command)
// 	}
// 	out, errOut, err := commandOut(command, "")
// 	if err != nil || errOut != "" {
// 		if errOut != "" {
// 			clog.Warn("warning: %v", errOut)
// 		}
// 		if err != nil {
// 			clog.Panic("error: %v", err)
// 		}
// 	}
// 	if out != "" {
// 		fmt.Println(out)
// 	}
// }
//
// // Return executes a command and returns the output
// func Return(command string, showCommand bool) string {
// 	if showCommand {
// 		color.HiBlue("$ %v", command)
// 	}
// 	out, errOut, err := commandOut(command, "")
// 	if err != nil || errOut != "" {
// 		if errOut != "" {
// 			clog.Warn("warning: %v", errOut)
// 		}
// 		if err != nil {
// 			clog.Panic("error: %v", err)
// 		}
// 	}
// 	if out != "" {
// 		fmt.Println(out)
// 	}
//
// 	return out
// }
