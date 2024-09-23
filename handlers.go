package main

import (
	"github.com/bwmarrin/discordgo"
)

func PingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
}

func TodoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if len(i.ApplicationCommandData().Options) == 0 {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Please provide a subcommand",
            },
        })
        return
    }
    switch i.ApplicationCommandData().Options[0].Name {
    case "add":
        user := i.ApplicationCommandData().Options[0].Options[0].StringValue()
        content := i.ApplicationCommandData().Options[0].Options[1].StringValue()
        err := addTodo(user, content)
        addTodoResponse(s, i, err)
    case "get":
        user := i.ApplicationCommandData().Options[0].Options[0].StringValue()
        todos, err := getTodos(user)
        getTodoResponse(s, i, todos, err)
    case "complete":
        id := i.ApplicationCommandData().Options[0].Options[0].IntValue()
        todo, err := completeTodoById(int(id))
        completeTodoResponse(s, i, &todo, err)
    }
}

func KaminoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
    switch i.ApplicationCommandData().Options[0].Name {
    case "get-pods":
        getPods(s, i)
    case "delete-pod":
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Deleting pod...",
            },
        })
        deletePod(s, i)
    case "bulk-delete":
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Deleting pods...",
            },
        })
        bulkDeletePods(s, i)
    case "refresh":
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Refreshing templates...",
            },
        })
        refreshTemplates(s, i)
    }

}
