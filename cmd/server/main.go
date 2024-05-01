package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/osousa/resticky/config"
	"github.com/osousa/resticky/internal"
	"go.uber.org/zap"
)

// ╔════════════════════════════════╗
// ║                API             ║
// ╚════════════════════════════════╝
type App struct {
	CLI    internal.CLI
	Log    *zap.Logger
	Ctx    context.Context
	Server internal.Server
	signal chan os.Signal
}

func ProvideApp(
	ctx context.Context,
	log *zap.Logger,
	srv internal.Server,
) *App {
	return &App{
		Ctx:    ctx,
		Log:    log,
		Server: srv,
		signal: make(chan os.Signal, 1),
	}
}

// @title			smid API
// @version		1.0.0
// @BasePath		/
// @description	Token generation and validation API
func main() {
	// Main context for the application
	ctx := context.Background()
	// Provide custom configuration
	cfg := ProvideAppConfig()
	// Load configuration file
	cnf := config.Read(ctx, cfg)
	// Create the application instance
	app := NewApp(ctx, cnf)
	// We have liftoff!
	app.Start()
}

// NewApp function is used to initialize the application.
// It uses wire to inject all dependencies. (see wire.go)
func NewApp(ctx context.Context, cfg *AppConfig) *App {
	app, err := InitializeApp(
		ctx,
		cfg.Server,
		cfg.ConnMan,
	)

	// If any error occurs, panic
	if err != nil {
		panic(err)
	}

	return app
}

// Start function is used to start the application.
// Server.Start is non-blocking, runs in a goroutine.
func (app *App) Start() {
	// Starting the application
	app.Log.Info("Starting App")
	// Start the gRPC server in a goroutine
	app.Server.Start()

	// Reg signal handler for graceful stop
	signal.Notify(app.signal,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT)

	// Start the gRPC server in a goroutine
	app.Server.Start()

	// Wait for the application to finish
	app.Wait()
}

// The application will finish when signal channel's closed
// or receives a signal from the OS (linux only)
func (app *App) Wait() {
	if sig := <-app.signal; sig != nil {
		app.Log.Info(
			"Received signal",
			zap.String("signal", sig.String()),
		)
		//TODO: Add graceful exit code here
	}

	app.Stop()
}

// All the connections will be closed. Don't Close
// other resources here, only the http server.
func (app *App) Stop() {
	app.Log.Info("Stopping App")
	app.Server.Stop()
}
