package snart

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (b *Bot) Start() error {
	_f := "(*Bot).Start"

	b.Startup = time.Now()
	Log.Infof(_f, "startup at %s", b.Startup)

	err := b.DB.Start()
	if err != nil {
		err = fmt.Errorf("db start: %w", err)
		Log.Error(_f, err)
		return err
	}
	Log.Info(_f, "db started")

	tok, err := b.Token()
	if err != nil {
		err = fmt.Errorf("token: %w", err)
		Log.Error(_f, err)
		return err
	}
	Log.Info(_f, "get token")
	b.Session.Token = tok.Value

	err = b.Session.Open()
	if err != nil {
		err = fmt.Errorf("session open: %w", err)
		Log.Error(_f, err)
		return err
	}
	Log.Info(_f, "session opened")

	b.WaitReady()
	Log.Info(_f, "ready")

	signal.Notify(b.Sig, os.Interrupt)
	signal.Notify(b.Sig, syscall.SIGTERM)
	<-b.Sig
	Log.Info(_f, "interrupt")

	if !b.Session.State.User.Bot {
		err = b.Session.Logout()
		if err != nil {
			err = fmt.Errorf("logout: %w", err)
			Log.Error(_f, err)
			return err
		} else {
			Log.Info(_f, "logged out")
			return nil
		}
	}

	err = b.Session.Close()
	if err != nil {
		err = fmt.Errorf("close: %w", err)
		Log.Error(_f, err)
		return err
	} else {
		Log.Info(_f, "session closed")
		return nil
	}
}

func (b *Bot) Uptime() time.Duration {
	return time.Since(b.Startup)
}
