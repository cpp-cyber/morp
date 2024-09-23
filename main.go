package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var (
    ctx = context.Background()
    db *gorm.DB
    s *discordgo.Session
    config = Config{}
    cl *http.Client
)

func init() {
    ReadConfig(&config)

    var err error
    s, err = discordgo.New("Bot " + config.DiscordToken)
    if err != nil {
        log.Fatalf("Cannot create a new session: %v", err)
    }

    cl = &http.Client{}
    NewKaminoClient()
    go refreshLogin()

    db = ConnectDB()
    db.AutoMigrate(&todo{})

    s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
        switch i.Type {
        case discordgo.InteractionApplicationCommand:
            if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
                h(s, i)
            }
        }
    })
}

func main() {
    s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
        log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
    })

    err := s.Open()
    if err != nil {
        log.Fatalf("Cannot open a new session: %v", err)
    }

    s.Identify.Intents = discordgo.IntentsGuildMessages

    fmt.Println(s.State)

    log.Println("Adding commands...")
    registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
    for i, v := range commands {
        cmd, err := s.ApplicationCommandCreate(s.State.User.ID, config.DiscordGuildID, v)
        if err != nil {
            log.Panicf("Cannot create '%v' command: %v", v.Name, err)
        }
        registeredCommands[i] = cmd
        log.Printf("Added \"%s\"", v.Name)
    }

    defer s.Close()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    log.Println("Press Ctrl+C to exit")
    <-stop

    log.Println("Gracefully shutting down.")
}
