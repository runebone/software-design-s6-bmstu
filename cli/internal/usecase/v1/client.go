package v1

import (
	"cli/internal/dto"
	"cli/internal/service"
	"cli/internal/usecase"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ClientUseCase struct {
	svc service.AggregatorService
}

func NewClientUseCase(svc service.AggregatorService) usecase.Client {
	return &ClientUseCase{
		svc: svc,
	}
}

func (uc *ClientUseCase) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.svc.Register(ctx, username, email, password)
	return tokens, err
}

func (uc *ClientUseCase) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.svc.Login(ctx, email, password)
	return tokens, err
}

func (uc *ClientUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	resp, err := uc.svc.Refresh(ctx, refreshToken)
	return resp, err
}

func (uc *ClientUseCase) Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	resp, err := uc.svc.Validate(ctx, token)
	return resp, err
}

func (uc *ClientUseCase) Logout(ctx context.Context, refreshToken string) error {
	err := uc.svc.Logout(ctx, refreshToken)
	return err
}

func (uc *ClientUseCase) ShowBoards(ctx context.Context) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		resp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = resp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	boards, err := uc.svc.ShowBoards(ctx)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for i, board := range boards {
		fmt.Printf("%d. %s\nTitle: %s\n", i+1, board.ID, board.Title)
	}
}

func (uc *ClientUseCase) ShowBoard(ctx context.Context, boardID string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		resp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = resp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	columns, err := uc.svc.ShowBoard(ctx, boardID)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for i, column := range columns {
		fmt.Printf("%d. %s\nTitle: %s\n", i+1, column.ID, column.Title)
	}
}

func (uc *ClientUseCase) ShowColumn(ctx context.Context, columnID string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	cards, err := uc.svc.ShowColumn(ctx, columnID)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for i, card := range cards {
		fmt.Printf("%d. %s\nTitle: %s\n", i+1, card.ID, card.Title)
	}
}

func (uc *ClientUseCase) ShowCard(ctx context.Context, cardID string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	card, err := uc.svc.ShowCard(ctx, cardID)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Title: %s\nDescription: %s\n", card.Title, card.Description)
}

func (uc *ClientUseCase) CreateBoard(ctx context.Context, title string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	resp, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		fmt.Println("access token expired")
		return
	}

	userID, err := uuid.Parse(resp.UserID)
	if err != nil {
		fmt.Println("failed parsing user uuid")
		return
	}

	board := dto.Board{
		UserID: userID,
		Title:  title,
	}

	err = uc.svc.CreateBoard(ctx, board)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Board successfully created.")
}

func (uc *ClientUseCase) CreateColumn(ctx context.Context, boardIDstr, title string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	resp, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		fmt.Println("access token expired")
		return
	}

	userID, err := uuid.Parse(resp.UserID)
	if err != nil {
		fmt.Println("failed parsing user uuid")
		return
	}

	boardID, err := uuid.Parse(boardIDstr)
	if err != nil {
		fmt.Println("failed parsing board uuid")
		return
	}

	column := dto.Column{
		UserID:  userID,
		BoardID: boardID,
		Title:   title,
	}

	err = uc.svc.CreateColumn(ctx, column)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Column successfully created.")
}

func (uc *ClientUseCase) CreateCard(ctx context.Context, columnIDstr, title, description string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	resp, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		fmt.Println("access token expired")
		return
	}

	userID, err := uuid.Parse(resp.UserID)
	if err != nil {
		fmt.Println("failed parsing user uuid")
		return
	}

	columnID, err := uuid.Parse(columnIDstr)
	if err != nil {
		fmt.Println("failed parsing column uuid")
		return
	}

	card := dto.Card{
		UserID:      userID,
		ColumnID:    columnID,
		Title:       title,
		Description: description,
	}

	err = uc.svc.CreateCard(ctx, card)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Card successfully created.")
}

func (uc *ClientUseCase) UpdateBoard(ctx context.Context, boardIDstr, title string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	boardID, err := uuid.Parse(boardIDstr)
	if err != nil {
		fmt.Println("failed parsing board uuid")
		return
	}

	board := dto.Board{
		ID:    boardID,
		Title: title,
	}

	err = uc.svc.UpdateBoard(ctx, &board)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Board successfully updated.")
}

func (uc *ClientUseCase) UpdateColumn(ctx context.Context, columnIDstr, title string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	columnID, err := uuid.Parse(columnIDstr)
	if err != nil {
		fmt.Println("failed parsing column uuid")
		return
	}

	column := dto.Column{
		ID:    columnID,
		Title: title,
	}

	err = uc.svc.UpdateColumn(ctx, &column)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Column successfully updated.")
}

func (uc *ClientUseCase) UpdateCardTitle(ctx context.Context, cardIDstr, title string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		fmt.Println("failed parsing card uuid")
		return
	}

	card := dto.Card{
		ID:    cardID,
		Title: title,
	}

	err = uc.svc.UpdateCard(ctx, &card)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Card title successfully updated.")
}

func (uc *ClientUseCase) UpdateCardDescription(ctx context.Context, cardIDstr, description string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		fmt.Println("failed parsing card uuid")
		return
	}

	card := dto.Card{
		ID:          cardID,
		Description: description,
	}

	err = uc.svc.UpdateCard(ctx, &card)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Card description successfully updated.")
}

func (uc *ClientUseCase) MoveCard(ctx context.Context, cardIDstr, columnIDstr string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		fmt.Println("failed parsing card uuid")
		return
	}

	columnID, err := uuid.Parse(columnIDstr)
	if err != nil {
		fmt.Println("failed parsing column uuid")
		return
	}

	card := dto.Card{
		ID:       cardID,
		ColumnID: columnID,
	}

	err = uc.svc.UpdateCard(ctx, &card)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Card successfully moved.")
}

func (uc *ClientUseCase) DeleteBoard(ctx context.Context, id string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	err = uc.svc.DeleteBoard(ctx, id)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Board successfully deleted.")
}

func (uc *ClientUseCase) DeleteColumn(ctx context.Context, id string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	err = uc.svc.DeleteColumn(ctx, id)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Column successfully deleted.")
}

func (uc *ClientUseCase) DeleteCard(ctx context.Context, id string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	err = uc.svc.DeleteCard(ctx, id)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Println("Card successfully deleted.")
}

func (uc *ClientUseCase) Stats(ctx context.Context, from, to string) {
	tokens, ok := ctx.Value("tokens").(*dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
		return
	}

	_, err := uc.svc.Validate(ctx, tokens.AccessToken)
	if err != nil {
		// fmt.Println("Access token expired. Refreshing.")
		refreshResp, err := uc.svc.Refresh(ctx, tokens.RefreshToken)
		if err != nil {
			fmt.Println("Please log in again.")
			return
		}

		tokens.AccessToken = refreshResp.AccessToken

		fn, ok := ctx.Value("saveFunc").(func(*dto.Tokens))
		if !ok {
			fmt.Println("failed to get saveFunc from context")
			return
		}

		fn(tokens)
	}

	stats, err := uc.svc.Stats(ctx, from, to)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	for i, stat := range stats {
		date := stat.Date.Format("02-01-2006")
		fmt.Printf("%d. Date: %s\n", i+1, date)
		fmt.Printf("Users:\n")
		for j, user := range stat.Users {
			fmt.Printf("    %d. %s (%s)\n", j+1, user.Username, user.Email)
		}
		fmt.Printf("Cards:\n")
		for j, card := range stat.Cards {
			fmt.Printf("    %d. %s\n", j+1, card.Title)
		}
		fmt.Printf("Number of cards created by new users: %d\n", stat.NumCardsByNewUsers)
	}
}
