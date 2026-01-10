package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/chachacollins/chroma/cfg"
	"github.com/chachacollins/chroma/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chroma",
	Short: "chroma is a command line documentation generator",
	Long: `A Fast and Flexible documentation generator built with go and powered by Mistral ai Codestral
               `,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var markdownCmd = &cobra.Command{
	Use:   "md [input-file] [output-file]",
	Short: "Generate markdown documentation from source files",
	Long: `Generates comprehensive markdown documentation by analyzing the specified input file.
The generated documentation is saved to the specified output file. This command supports
various file formats and produces well-structured markdown output.

If the output file already exists, its contents will be overwritten. Make sure to backup
any existing documentation before running this command.

Example:
  chroma md input.go api.md
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			utils.PrintError("Error: please provide  [input-file] [output-file]")
			os.Exit(1)
		}
		fileName := args[0]
		writeFilePath := args[1]
		split := strings.Split(writeFilePath, ".")
		if split[1] != "md" {
			utils.PrintError("Error: please provide a filename with a markdown extension")
			os.Exit(1)
		}
		file := readFile(fileName)
		generateMd(writeFilePath, file)
	},
}

var inlineCmd = &cobra.Command{
	Use:   "il [file]",
	Short: "Add AI-generated documentation comments to source code",
	Long: `Analyzes the specified source file and automatically generates meaningful inline documentation.
The AI will add descriptive comments for functions, types, and important code blocks directly
within your source code. The original file will be updated with the new documentation.

CAUTION: This command will modify your source file. Please ensure you have a backup or
version control in place before proceeding.

Example:
  chroma il main.go
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			utils.PrintError("Error: please provide a <code path> ")
			os.Exit(1)
		}
		fileName := args[0]
		file := readFile(fileName)
		inlineComm(fileName, file)
	},
}

var starCmd = &cobra.Command{
	Use:   "star [source-file]",
	Short: "Generate Starlight documentation site for your project",
	Long: `Creates a comprehensive Astro Starlight documentation site from your project's source code.
This command will:
- Generate API documentation from your source files
- Create a well-structured Starlight navigation sidebar
- Set up customizable theme and styling
- Add example usage and code snippets
- Organize content into logical sections

The documentation will be generated in a 'docs' directory within your project.
If the docs directory already exists, its contents will be overwritten.

Example:
  chroma star ./[file-name]

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			utils.PrintError("Error: please provide a <code path> or serve command ")
			os.Exit(1)
		}
		fileName := args[0]
		file := readFile(fileName)
		getDocs()
		starLight(fileName, file)
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure API credentials for documentation generation",
	Long: `Set up your environment with the necessary API credentials for AI-powered documentation.

This command will:
- Guide you through entering your API key
- Validate the provided credentials
- Store them securely in your environment
- Create a configuration file if needed

Example:
  chroma init

The command will prompt you for your API key and handle the secure storage of your credentials.
This is required before using other documentation generation commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Init()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve ",
	Short: "Start a local server for your Starlight documentation",
	Long: `Launch a development server to preview your generated Starlight documentation site.

This command will:
- Start a local development server
- Watch for file changes and update in real-time
- Open your documentation in the default browser

The server defaults to port 4321 if not specified.

Note: Ensure you have generated your documentation using the 'star' command before serving.`,
	Run: func(cmd *cobra.Command, args []string) {
		serveDocs()
	},
}

func init() {
	rootCmd.AddCommand(markdownCmd)
	rootCmd.AddCommand(inlineCmd)
	rootCmd.AddCommand(starCmd)
	rootCmd.AddCommand(initCmd)
	starCmd.AddCommand(serveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
