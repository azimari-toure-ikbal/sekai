package cmd

import (
	"fmt"
	"os"

	"github.com/azimari-toure-ikbal/sekai-core/internal/nextjs"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sekai-core",
	Short: "CLI Tool written in go to make localization easier",
	Long: `
Sekai aims to help developer do localization easier and faster for their project of any size.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		var files []string

		acceptedLang := []string{"fr", "en", "de", "it", "jp", "po", "es"}
		// Helper function to check if the selected language is valid
		isValidLang := func(lang string, acceptedLangs []string) bool {
			for _, l := range acceptedLangs {
				if lang == l {
					return true
				}
			}
			return false
		}

		if !isValidLang(input, acceptedLang) {
			fmt.Printf("You must chose a valid input language in this list %v", acceptedLang)
			return
		}

		if !isValidLang(output, acceptedLang) {
			fmt.Printf("You must chose a valid output language in this list %v", acceptedLang)
			return
		}
	
		if input == output {
			fmt.Printf("The input and output languages cannot be the same")
			return
		}

		switch env {
		case "nextjs":
			err := nextjs.RunForNext(&files, &input, &output)
			if err != nil {
				fmt.Printf("---- Something went wrong: %v ---- \n", err)
				return
			}
		default:
			fmt.Printf("You provided a wrong value of env. We only support --env nextjs for the moment")
	}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sekai.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("env", "e", "nextjs", "The env of your project")
	rootCmd.Flags().StringP("input", "i", "", "The input language which correspond to the language of your project")
	rootCmd.Flags().StringP("output", "o", "", "The desired output language from the localization")
}


