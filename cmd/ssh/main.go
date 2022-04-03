package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aligator/HideAndShell/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bubbleMiddleware "github.com/charmbracelet/wish/bubbletea"
	loggingMiddleware "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

const host = "0.0.0.0"
const port = 2223

func AllPasswords(ctx ssh.Context, password string) bool {
	return true
}

func AllKeys(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis
func teaHandler() func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		_, _, active := s.Pty()
		if !active {
			fmt.Println("no active terminal, skipping")
			return nil, nil
		}

		g := game.New()

		return g, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
	}
}

func main() {
	rand.Seed(time.Now().UnixMilli())

	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithPublicKeyAuth(AllKeys),
		wish.WithPasswordAuth(AllPasswords),
		wish.WithMiddleware(
			bubbleMiddleware.Middleware(teaHandler()),
			loggingMiddleware.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
