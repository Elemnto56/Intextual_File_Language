package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "./ITX-CLI",
		Short: "Run and manage Intext scripts from the terminal",
		Long:  "ITX-CLI gives you tools to run, validate, and interact with Intext files directly from the terminal. You can control ISEC behavior, inspect file structure, and customize script execution with ease.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to ITX-CLI! Use --help for usage info")
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "See version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`
 __     __   __     ______   ______     __  __     ______  
/\ \   /\ "-.\ \   /\__  _\ /\  ___\   /\_\_\_\   /\__  _\ 
\ \ \  \ \ \-.  \  \/_/\ \/ \ \  __\   \/_/\_\/_  \/_/\ \/ 
 \ \_\  \ \_\\"\_\    \ \_\  \ \_____\   /\_\/\_\    \ \_\ 
  \/_/   \/_/ \/_/     \/_/   \/_____/   \/_/\/_/     \/_/ 
                                                           `)
			fmt.Println(`
			-----------
		Version: v0.8.1									
		Codename: Where's the Logic?
		Developer: Elemnto56 @ Github`)

		},
	}

	var runCmd = &cobra.Command{
		Use:   "run [filename].itx",
		Short: "Execute your Intext files",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			Lexer(filename, nil, true)
			Parser(nil)
			//Validator(nil, nil) UPDATE WHEN CAN
			Interpreter(nil, nil)
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Test stuff",
		Run: func(cmd *cobra.Command, args []string) {
			a, err := EvaluateExpression(`"he"`, nil)
			Check(err)
			fmt.Println(a)
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
