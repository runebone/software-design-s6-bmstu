package main

import (
	"cli/internal/adapter/logger"
	"cli/internal/adapter/service/http"
	"cli/internal/config"
	"cli/internal/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata"

	v1 "cli/internal/usecase/v1"

	"github.com/spf13/cobra"
)

var abspath string = "/home/rukost/University/software-design-s6-bmstu.git/lab5/src/"

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Couldn't set timezone: %v", err)
	}
	time.Local = loc
}

func readTokens(tokens *dto.Tokens, filePath string) error {
	path := abspath + filePath
	file, err := os.Open(path)
	if err != nil {
		// fmt.Printf("Error opening file: %s\n", filePath)
		return err
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(tokens); err != nil {
		// fmt.Printf("Error parsing file: %s\n", filePath)
		return err
	}

	return nil
}

func saveTokens(tokens *dto.Tokens, filePath string) error {
	path := abspath + filePath
	tokensJson, err := json.Marshal(tokens)
	if err != nil {
		fmt.Println("Error marshalling tokens")
		return err
	}

	err = os.WriteFile(path, tokensJson, 0644)
	if err != nil {
		fmt.Printf("Error saving tokens to file: %s\n", filePath)
	}

	return nil
}

func main() {
	cfg, err := config.LoadConfig(abspath + "config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	var Tokens dto.Tokens
	err = readTokens(&Tokens, cfg.Client.TokensPath)
	if err != nil {
		saveTokens(&dto.Tokens{}, cfg.Client.TokensPath)
	}

	logger := logger.NewZapLogger(config.LogConfig{Path: abspath + "cli/cli.log", Level: cfg.Client.Log.Level}) // config.Client.Log)
	ac := cfg.Aggregator
	svc := http.NewAggregatorService(ac.BaseURL, 5*time.Second, logger)

	client := v1.NewClientUseCase(svc)

	rootCmd := &cobra.Command{Use: "todo"}

	// Register command
	registerCmd := &cobra.Command{
		Use:   "register [username] [email] [password]",
		Short: "Register a new user",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := client.Register(context.Background(), args[0], args[1], args[2])
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			saveTokens(tokens, cfg.Client.TokensPath)
		},
	}
	rootCmd.AddCommand(registerCmd)

	// Login command
	loginCmd := &cobra.Command{
		Use:   "login [email] [password]",
		Short: "Login a user",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := client.Login(context.Background(), args[0], args[1])
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Println("Logged in.")
			saveTokens(tokens, cfg.Client.TokensPath)
		},
	}
	rootCmd.AddCommand(loginCmd)

	// Logout command
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout the user",
		Run: func(cmd *cobra.Command, args []string) {
			err := client.Logout(context.Background(), Tokens.RefreshToken)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Println("Logged out.")
			saveTokens(&dto.Tokens{}, cfg.Client.TokensPath)
		},
	}
	rootCmd.AddCommand(logoutCmd)

	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources",
	}

	// Create board command
	createBoardCmd := &cobra.Command{
		Use:   "board [title]",
		Short: "Create a new board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.CreateBoard(ctx, args[0])
		},
	}
	createCmd.AddCommand(createBoardCmd)

	// Create column command
	createColumnCmd := &cobra.Command{
		Use:   "column [board_id] [title]",
		Short: "Create a new column in a board",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.CreateColumn(ctx, args[0], args[1])
		},
	}
	createCmd.AddCommand(createColumnCmd)

	// Create card command
	createCardCmd := &cobra.Command{
		Use:   "card [column_id] [title] [description]",
		Short: "Create a new card in a column",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			var description string
			if len(args) == 3 {
				description = args[2]
			} else {
				description = ""
			}
			client.CreateCard(ctx, args[0], args[1], description)
		},
	}
	createCmd.AddCommand(createCardCmd)
	rootCmd.AddCommand(createCmd)

	// Show command
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show resources",
	}

	// Show boards command
	showBoardsCmd := &cobra.Command{
		Use:   "boards",
		Short: "Show all boards",
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.ShowBoards(ctx)
		},
	}
	showCmd.AddCommand(showBoardsCmd)

	// Show board command
	showBoardCmd := &cobra.Command{
		Use:   "board [board_id]",
		Short: "Show a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.ShowBoard(ctx, args[0])
		},
	}
	showCmd.AddCommand(showBoardCmd)

	// Show column command
	showColumnCmd := &cobra.Command{
		Use:   "column [column_id]",
		Short: "Show a column",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.ShowColumn(ctx, args[0])
		},
	}
	showCmd.AddCommand(showColumnCmd)

	// Show card command
	showCardCmd := &cobra.Command{
		Use:   "card [card_id]",
		Short: "Show a card",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.ShowCard(ctx, args[0])
		},
	}
	showCmd.AddCommand(showCardCmd)
	rootCmd.AddCommand(showCmd)

	// Update command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
	}

	// Update board command
	updateBoardCmd := &cobra.Command{
		Use:   "board",
		Short: "Update a board",
	}

	// Update board title command
	updateBoardTitleCmd := &cobra.Command{
		Use:   "title [board_id] [new_title]",
		Short: "Update board title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.UpdateBoard(ctx, args[0], args[1])
		},
	}
	updateBoardCmd.AddCommand(updateBoardTitleCmd)
	updateCmd.AddCommand(updateBoardCmd)

	// Update column command
	updateColumnCmd := &cobra.Command{
		Use:   "column",
		Short: "Update a column",
	}

	// Update column title command
	updateColumnTitleCmd := &cobra.Command{
		Use:   "title [column_id] [new_title]",
		Short: "Update column title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.UpdateColumn(ctx, args[0], args[1])
		},
	}
	updateColumnCmd.AddCommand(updateColumnTitleCmd)
	updateCmd.AddCommand(updateColumnCmd)

	// Update card command
	updateCardCmd := &cobra.Command{
		Use:   "card",
		Short: "Update a card",
	}

	// Update card title command
	updateCardTitleCmd := &cobra.Command{
		Use:   "title [card_id] [new_title]",
		Short: "Update card title",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.UpdateCardTitle(ctx, args[0], args[1])
		},
	}
	updateCardCmd.AddCommand(updateCardTitleCmd)

	// Update card description command
	updateCardDescriptionCmd := &cobra.Command{
		Use:   "description [card_id] [new_description]",
		Short: "Update card description",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.UpdateCardDescription(ctx, args[0], args[1])
		},
	}
	updateCardCmd.AddCommand(updateCardDescriptionCmd)
	updateCmd.AddCommand(updateCardCmd)
	rootCmd.AddCommand(updateCmd)

	// Move command
	moveCmd := &cobra.Command{
		Use:   "move",
		Short: "Move stuff",
	}

	moveCardCmd := &cobra.Command{
		Use:   "card [card_id] [column_id]",
		Short: "Move card",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.MoveCard(ctx, args[0], args[1])
		},
	}
	moveCmd.AddCommand(moveCardCmd)
	rootCmd.AddCommand(moveCmd)

	// Delete command
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	// Delete board command
	deleteBoardCmd := &cobra.Command{
		Use:   "board [board_id]",
		Short: "Delete a board",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.DeleteBoard(ctx, args[0])
		},
	}
	deleteCmd.AddCommand(deleteBoardCmd)

	// Delete column command
	deleteColumnCmd := &cobra.Command{
		Use:   "column [column_id]",
		Short: "Delete a column",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.DeleteColumn(ctx, args[0])
		},
	}
	deleteCmd.AddCommand(deleteColumnCmd)

	// Delete card command
	deleteCardCmd := &cobra.Command{
		Use:   "card [card_id]",
		Short: "Delete a card",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.DeleteCard(ctx, args[0])
		},
	}
	deleteCmd.AddCommand(deleteCardCmd)
	rootCmd.AddCommand(deleteCmd)

	// Stats command
	statsCmd := &cobra.Command{
		Use:   "stats from [DD-MM-YYYY] to [DD-MM-YYYY]",
		Short: "Show stats for a time period",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx0 := context.WithValue(context.Background(), "tokens", &Tokens)
			ctx := context.WithValue(ctx0, "saveFunc", func(tokens *dto.Tokens) {
				saveTokens(tokens, cfg.Client.TokensPath)
			})
			client.Stats(ctx, args[0], args[1])
		},
	}
	rootCmd.AddCommand(statsCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
