package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/byte4cat/nbx/pkg/clog"
	"github.com/fatih/color"
)

// Run executes a command and prints its output line by line in real-time
func Run(command string, showCommand bool) {
	if showCommand {
		color.HiBlue("$ %v", command)
	}

	// default shell is bash
	cmd := exec.Command("bash", "-c", command)

	// get the stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		clog.Panic("failed to create stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		clog.Panic("failed to create stderr pipe: %v", err)
	}

	// start the command
	err = cmd.Start()
	if err != nil {
		clog.Panic("failed to start command: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2) // two goroutines for stdout and stderr

	// Goroutine reads stdout and prints it line by line
	go func() {
		defer wg.Done()
		defer stdoutPipe.Close() // close the pipe when done

		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // print each line immediately
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			// ignore io.EOF error, it's expected when the command ends
			// clog.Warn("error reading stdout: %v", err)
		}
	}()

	// Goroutine reads stderr and prints it line by line
	go func() {
		defer wg.Done()
		defer stderrPipe.Close() // close the pipe when done

		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			// print the stderr line by line
			clog.Warn("warnings: %v", scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			// ignore io.EOF error, it's expected when the command ends
			// clog.Warn("Error reading stderr: %v", err)
		}
	}()

	// wait for the command to finish
	cmdErr := cmd.Wait()

	// wait for the goroutines to finish processing
	wg.Wait()

	// process the command error, exit code, etc.
	if cmdErr != nil {
		clog.Panic("command execution error: %v", cmdErr)
	}
}

// Return function is trickier with real-time printing AND returning a single string.
// If you need both, you would need to stream and also buffer simultaneously.
// The original Return function's behavior of printing AFTER buffering is incompatible with real-time printing.
// If you need to capture the full output AND print real-time, you'd modify a function like this:
func Return(command string, showCommand bool) string {
	if showCommand {
		color.HiBlue("$ %v", command)
	}

	cmd := exec.Command("bash", "-c", command)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		clog.Panic("failed to create stdout pipe for Return: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		clog.Panic("failed to create stderr pipe for Return: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var stdoutBuf bytes.Buffer // Buffer to collect stdout
	var stderrBuf bytes.Buffer // Buffer to collect stderr (optional, but good for debugging)

	// Start command
	err = cmd.Start()
	if err != nil {
		clog.Panic("failed to start command for Return: %v", err)
	}

	// Goroutine to read stdout, print, and buffer
	go func() {
		defer wg.Done()
		defer stdoutPipe.Close()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)                  // Print line by line immediately
			stdoutBuf.WriteString(line + "\n") // Also write to buffer
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			// ignore io.EOF error, it's expected when the command ends
			// clog.Warn("Error reading stderr: %v", err)
		}
	}()

	// Goroutine to read stderr, print, and buffer
	go func() {
		defer wg.Done()
		defer stderrPipe.Close()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			clog.Warn("warnings: %v", line)    // Print stderr line by line
			stderrBuf.WriteString(line + "\n") // Also write to buffer
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			// ignore io.EOF error, it's expected when the command ends
			// clog.Warn("Error reading stderr: %v", err)
		}
	}()

	// Wait for the command to finish
	cmdErr := cmd.Wait()

	// Wait for reader goroutines to finish processing any remaining output
	wg.Wait()

	// Handle command execution error
	if cmdErr != nil {
		// The stderr goroutine already printed warnings.
		// Here we just check if the command failed.
		clog.Panic("command execution error for Return: %v", cmdErr)
	}

	// The original logic prints the buffer if not empty, then returns it.
	// Since we already printed line by line, this final print might be redundant
	// but we'll keep the original structure's intent of returning the full string.
	return stdoutBuf.String()
}
