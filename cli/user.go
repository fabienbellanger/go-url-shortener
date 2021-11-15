package cli

import (
	"log"

	"github.com/spf13/cobra"
)

type userRegister struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

var (
	userEmail    string
	userPassword string
)

func init() {
	userCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	userCmd.Flags().StringVarP(&userPassword, "password", "p", "", "user password")

	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "register",
	Short: "User creation",
	Long:  `User creation`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, db, err := initConfigLoggerDatabase()
		if err != nil {
			log.Fatalln(err)
		}
		println(logger, db)
		user := userRegister{
			Email:    userEmail,
			Password: userPassword,
		}
		log.Printf("%v\n", user)

		// Validate data
		// -------------

		// Register new user
		// -----------------

		// Display result
		// --------------
	},
}
