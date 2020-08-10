// Package bot contains the general workings of a Snart Bot.
package bot

import (
	"time"

	dg "github.com/bwmarrin/discordgo"

	"github.com/go-snart/snart/db"
	"github.com/go-snart/snart/route"
)

const _p = "snart/bot"

// Bot holds all the internal workings of a Snart bot.
type Bot struct {
	DB      *db.DB
	Session *dg.Session

	Router *route.Router
	Gamers []Gamer

	Interrupt chan Interrupt
	Startup   time.Time

	Ready bool
}

// New creates a Bot.
func New() *Bot {
	return &Bot{
		DB:      db.New(),
		Session: nil,

		Router: route.NewRouter(),
		Gamers: []Gamer{GamerUptime},

		Interrupt: make(chan Interrupt),
		Startup:   time.Now(),
	}
}
