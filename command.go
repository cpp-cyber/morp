package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with Pong!",
		},

        {
            Name:        "todo",
            Description: "Manage your todos",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "add",
                    Description: "Add a new todo",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "person",
                            Description: "User to assign the todo to",
                            Required:    true,
                        },
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "content",
                            Description: "Content of the todo",
                            Required:    true,
                        },
                    },
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "get",
                    Description: "Get all todos",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "person",
                            Description: "User to get todos for",
                            Required:    false,
                        },
                    },
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "complete",
                    Description: "Complete a todo",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionInteger,
                            Name:        "id",
                            Description: "ID of the todo to complete",
                            Required:    true,
                        },
                    },
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "update",
                    Description: "Update a todo",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionInteger,
                            Name:        "id",
                            Description: "ID of the todo to update",
                            Required:    true,
                        },
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "content",
                            Description: "New content of the todo",
                            Required:    true,
                        },
                    },
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "remove",
                    Description: "Remove a todo",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionInteger,
                            Name:        "id",
                            Description: "ID of the todo to remove",
                            Required:    false,
                        },
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "person",
                            Description: "Person to remove todos for",
                            Required:    false,
                        },
                    },
                },
            },
        },

        {
            Name:        "kamino",
            Description: "Commands for interacting with Kamino",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "get-pods",
                    Description: "Get all pods",
                },

                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "delete-pod",
                    Description: "Delete a pod",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "name",
                            Description: "Name of a single pod to delete.",
                            Required:    true,
                        },
                    },
                },  
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "bulk-delete",
                    Description: "Bulk delete pods",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "filter",
                            Description: "Filter for a bulk delete operation. Takes a comma separated list of filters.",
                            Required:    true,
                        },
                    },
                },  

                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "refresh",
                    Description: "Refresh templates",
                },

                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "competition-clone",
                    Description: "Clone pods for a competition and create respective users.",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "template",
                            Description: "Name of the template to clone pods for.",
                            Required:    true,
                        },
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "count",
                            Description: "Number of teams to clone for.",
                            Required:    true,
                        },
                    },
                },
            },
        },
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":   PingHandler,
        "todo":   TodoHandler,
        "kamino": KaminoHandler,
	}
)
