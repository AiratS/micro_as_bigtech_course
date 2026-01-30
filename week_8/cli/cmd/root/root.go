package root

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-app",
	Short: "My cli APP",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create smth",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete smth",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create new user",
	Run: func(cmd *cobra.Command, args []string) {
		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("no username: %v", err)
		}
		log.Printf("username: %s", usernameStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete user",
	Run: func(cmd *cobra.Command, args []string) {
		usernameStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("no username: %v", err)
		}
		log.Printf("username: %s", usernameStr)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "User name")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flas as required: %s", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "User name")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flas as required: %s", err.Error())
	}
}
