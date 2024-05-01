package internal

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/osousa/resticky/internal/protobuf/resticky"
	"go.uber.org/zap"
)

type CLI interface {
	Start() error
	Stop() (*resticky.RestickyResponse, error)
}

type LockCmd struct{}

func (r *LockCmd) Run(ctx context.Context, handlers *Handlers, log zap.Logger) error {
	log.Info("LockCmd called - locking databases")
	// connect to the databases
	var (
		resp *resticky.RestickyResponse
		req  *resticky.RestickyRequest
		err  error
	)

	req = &resticky.RestickyRequest{Id: "1"}
	resp, err = handlers.client.LockAll(ctx, req)
	if err != nil {
		err := fmt.Errorf("Error: %v - descriptor: %v", err, resp.Error)
		log.Error("could not lock databases", zap.Error(err))
		return err
	}

	return nil
}

type UnlockCmd struct {
	Timeout int `arg:"" name:"timeout" help:"Timeout for the lock."`
}

func (r *UnlockCmd) Run(ctx context.Context, handlers *Handlers, log zap.Logger) error {
	log.Info("UnlockCmd called - unlocking databases")
	// connect to the databases
	var (
		resp *resticky.RestickyResponse
		req  *resticky.RestickyRequest
		err  error
	)

	req = &resticky.RestickyRequest{Id: "1"}
	resp, err = handlers.client.UnlockAll(ctx, req)
	if err != nil {
		err := fmt.Errorf("Error: %v - descriptor: %v", err, resp.Error)
		log.Error("could not lock databases", zap.Error(err))
		return err
	}

	return nil
}

type Context struct {
	client resticky.RestickyServiceClient
}

type cmdList struct {
	Debug  bool      `help:"Enable debug mode."`
	Lock   LockCmd   `cmd:"" help:"Lock all databases"`
	Unlock UnlockCmd `cmd:"" help:"Unlock all databases"`
}

type Handlers struct {
	client resticky.RestickyServiceClient
}

type cli struct {
	Ctx      context.Context
	Handlers *Handlers
	Log      *zap.Logger
	Cmd      cmdList
}

func ProvideCLI(ctx context.Context, log *zap.Logger, rsc resticky.RestickyServiceClient) CLI {
	cmdList := cmdList{}
	lockCmd := LockCmd{}
	unlockCmd := UnlockCmd{}
	cmdList.Lock = lockCmd
	cmdList.Unlock = unlockCmd
	Hanlders := &Handlers{
		client: rsc,
	}

	return &cli{
		Ctx:      ctx,
		Log:      log,
		Cmd:      cmdList,
		Handlers: Hanlders,
	}
}

// CLI is the command line interface for the application
func (c *cli) Start() error {
	cli := kong.Parse(&c.Cmd)
	c.Log.Info("Starting CLI")
	cli.BindTo(c.Ctx, (*context.Context)(nil))
	err := cli.Run(c.Handlers, *c.Log)

	if err != nil {
		c.Log.Error("Error running CLI", zap.Error(err))
		return err
	}

	return nil
}

// stop CLI
func (c *cli) Stop() (*resticky.RestickyResponse, error) {
	c.Log.Info("Stopping CLI")
	var (
		resp *resticky.RestickyResponse
		req  *resticky.RestickyRequest
		err  error
	)

	resp, err = c.Handlers.client.UnlockAll(c.Ctx, req)
	if err != nil {
		err := fmt.Errorf("Error: %v - descriptor: %v", err, resp.Error)
		c.Log.Error("could not unlock databases", zap.Error(err))
		return nil, err
	}

	return resp, nil
}
