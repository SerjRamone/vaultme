package client

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/SerjRamone/vaultme/internal/config"
	"github.com/SerjRamone/vaultme/internal/models"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TerminalUI - TUI for VaultMe client app
type TerminalUI struct {
	// tview instance
	tview *tview.Application
	// content pages
	pages *tview.Pages
	// gRPC client
	client *Client
}

// Start - starts VaultMe TUI
func (i *TerminalUI) Start(ctx context.Context, log *zap.Logger, cfg *config.Client) error {
	// chann for stopping application
	exitChan := make(chan error)

	// Create a new rivo/tview app
	i.tview = tview.NewApplication()

	// Create a new tview.Pages for toggle content pages
	i.pages = tview.NewPages()

	client, err := NewClient(log, cfg)
	if err != nil {
		return fmt.Errorf("create client error: %w", err)
	}

	i.client = client
	i.drawLoginForm(ctx)

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(i.pages, 0, 1, true), 0, 1, true)

	i.tview.SetRoot(flex, true).SetFocus(flex).EnableMouse(true)

	// run interface
	go func() {
		exitChan <- i.tview.Run()
	}()

	// wait context Done or error
	select {
	case err := <-exitChan:
		return err
	case <-ctx.Done():
		i.tview.Stop()
		return <-exitChan
	}
}

// drawLoginForm - draws login/registration form
func (i *TerminalUI) drawLoginForm(ctx context.Context) {
	var userDTO models.UserDTO
	form := tview.NewForm().
		AddInputField("Login", "", 50, nil, func(v string) {
			userDTO.Login = v
		}).
		AddPasswordField("Password", "", 50, '*', func(v string) {
			userDTO.Password = v
		}).
		AddButton("Login", func() {
			// get user from server
			u, err := userDTO.GetUser(ctx, i.client)
			if err != nil {
				i.drawErrorModal(err.Error())
				return
			}
			_ = u // TODO: get user's items from server
			i.drawErrorModal("Login successful")
		}).
		AddButton("Register", func() {
			u, err := userDTO.CreateUser(ctx, i.client)
			if err != nil {
				i.drawErrorModal(err.Error())
				return
			}
			_ = u // TODO: get user's items from server
			i.drawErrorModal("Registration successful")
		}).
		AddButton("Exit", i.drawExitModal)

	l1 := "_    __               __ __   __  ___    "
	l2 := "| |  / /____ _ __  __ / // /_ /  |/  /___ "
	l3 := "| | / // __ `// / / // // __// /|_/ // _ \\"
	l4 := "| |/ // /_/ // /_/ // // /_ / /  / //  __/"
	l5 := "|___/ \\__,_/ \\__,_//_/ \\__//_/  /_/ \\___/ "
	frame := tview.NewFrame(form).
		AddText(l1, true, tview.AlignCenter, tcell.ColorWhite).
		AddText(l2, true, tview.AlignCenter, tcell.ColorWhite).
		AddText(l3, true, tview.AlignCenter, tcell.ColorWhite).
		AddText(l4, true, tview.AlignCenter, tcell.ColorWhite).
		AddText(l5, true, tview.AlignCenter, tcell.ColorWhite)

	grid := tview.NewGrid().
		SetRows(-1, -80, 0).
		SetColumns(-10, 60, -10).
		AddItem(frame, 1, 1, 1, 1, 0, 0, true)

	i.pages.AddPage("Login", grid, true, true)
}

// drawErrorModal - draws error in modal window
func (i *TerminalUI) drawErrorModal(text string) {
	name := "Error"

	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "Ok" {
				i.pages.RemovePage(name)
			}
		})

	i.pages.AddPage(name, modal, true, true)
}

// drawExitModal - draws exit in modal dialog
func (i *TerminalUI) drawExitModal() {
	name := "Exit dialog"

	modal := tview.NewModal().
		SetText("Do you really want to close app?").
		AddButtons([]string{"Cancel", "Exit"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "Exit" {
				i.tview.Stop()
			}
			if buttonLabel == "Cancel" {
				i.pages.RemovePage(name)
			}
		})

	i.pages.AddPage(name, modal, true, true)
}
