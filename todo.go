package main

import (
	"fmt"

	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

type todo struct {
    Id int `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
    User string
    Task string
    Completed bool
}

func completeTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todo *todo, err error) {
    if err != nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Failed to complete todo: " + err.Error(),
            },
        })
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todo Completed")
    embed.SetColor(0xffab40)
    embed.AddField("Person", todo.User)
    embed.AddField("Task", todo.Task)

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
        },
    })
}

func addTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
    if err != nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Failed to add todo: " + err.Error(),
            },
        })
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todo Added")
    embed.SetColor(0xffab40)
    embed.AddField("Person", i.ApplicationCommandData().Options[0].Options[0].StringValue())
    embed.AddField("Task", i.ApplicationCommandData().Options[0].Options[1].StringValue())

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
        },
    })
}

func getTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todos []todo, err error) {
        if err != nil {
            s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
                Type: discordgo.InteractionResponseChannelMessageWithSource,
                Data: &discordgo.InteractionResponseData{
                    Content: "Failed to get todos: " + err.Error(),
                },
            })
            return
        }

        embed := embed.NewEmbed()
        embed.SetTitle("Todos for " + todos[0].User)
        embed.SetColor(0xffab40)
        for _, t := range todos {
            embed.AddField(fmt.Sprintf("Id: %d", t.Id), fmt.Sprintf("Task: %s\n Completed: %t", t.Task, t.Completed))
        }
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
            },
        })
}

func getAllTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todos []todo, err error) {
    if err != nil {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Failed to get todos: " + err.Error(),
            },
        })
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todos")
    embed.SetColor(0xffab40)
    for _, t := range todos {
        embed.AddField(fmt.Sprintf("Id: %d", t.Id), fmt.Sprintf("Person: %s\n Task: %s\n Completed: %t", t.User, t.Task, t.Completed))
    }
    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
        },
    })
}
