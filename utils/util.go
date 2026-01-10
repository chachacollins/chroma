package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

func MakeDir(dirName string) {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create a dir: %s", err)
	}
}

func WriteFile(filename string, data string) {
	file, err := os.Create(filename)
	if err != nil {
		PrintError("Error: could not create file")
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(data)
	if err != nil {
		PrintError("Error: could not write to file")
	}
	err = writer.Flush()

	if err != nil {
		PrintError("could not flush buffer")
	}
}

func PrintSuccess(msg string) {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	fmt.Println(style.Render(msg))
}

func PrintError(msg string) {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	fmt.Println(style.Render(msg))
}

func PrintLog(msg string) {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	fmt.Println(style.Render(msg))
}

func ExtractCodeBlocks(md string) string {
	lines := strings.Split(md, "\n")
	var result strings.Builder
	counter := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "```") && counter == 0 {
			counter += 1
		} else {
			if(len(lines) - 1 != i) {
				result.WriteString(line)
				result.WriteString("\n")
			}
		}
	}
	return result.String()
}
