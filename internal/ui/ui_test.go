package ui

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/fatih/color"
)

func TestPrintDashes(t *testing.T) {
	output := captureOutput(PrintDashes)

	expected := "—————————————————————————————————————————————————————————\n"
	if output != expected {
		t.Errorf("PrintDashes() = %v, want %v", output, expected)
	}
}

func TestPrintMsgWithStatus(t *testing.T) {
	type args struct {
		status string
		color  string
		format string
		a      []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Red Color",
			args: args{
				status: "ERROR",
				color:  "red",
				format: "Failed to connect to %s",
				a:      []any{"server"},
			},
			want: color.New(color.FgRed).SprintFunc()("ERROR") + " — Failed to connect to server",
		},
		{
			name: "Test Green Color",
			args: args{
				status: "SUCCESS",
				color:  "green",
				format: "Records %s",
				a:      []any{"accurate"},
			},
			want: color.New(color.FgGreen).SprintFunc()("SUCCESS") + " — Records accurate",
		},
		{
			name: "Test HiBlue Color",
			args: args{
				status: "INFO",
				color:  "hiBlue",
				format: "Loading %s",
				a:      []any{"module"},
			},
			want: color.New(color.FgHiBlue).SprintFunc()("INFO") + " — Loading module",
		},
		{
			name: "Test Non-existent Color",
			args: args{
				status: "UNKNOWN",
				color:  "unknown",
				format: "This should use default format for %s",
				a:      []any{"output"},
			},
			want: "UNKNOWN — This should use default format for output",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				PrintMsgWithStatus(tt.args.status, tt.args.color, tt.args.format, tt.args.a...)
			})
			if output != tt.want {
				t.Errorf("PrintMsgWithStatus() = %v, want %v", output, tt.want)
			}
		})
	}
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	output := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output <- buf.String()
	}()
	f()

	_ = w.Close()
	os.Stdout = stdout

	return <-output
}
