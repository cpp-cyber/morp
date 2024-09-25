package main

import (
	"fmt"
	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

type todo struct {
	Id        int `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
	User      string
	Task      string
	Completed bool
}

func getDiscordUser(userMention string) (*discordgo.Member, error) {
	// Extract the user ID from the mention using regex
	re := regexp.MustCompile(`<@!?(\d+)>`)
	matches := re.FindStringSubmatch(userMention)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no user ID found in mention")
	}
	userID := matches[1]

	// Get the user by ID
	member, err := s.GuildMember(config.DiscordGuildID, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user from guild: %v", err)
	}

	return member, nil
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
	embed.SetColor(0xE0C460)
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
	embed.SetColor(0xE0C460)
	embed.AddField("Person", i.ApplicationCommandData().Options[0].Options[0].StringValue())
	embed.AddField("Task", i.ApplicationCommandData().Options[0].Options[1].StringValue())

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func getTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todos []todo, numCompleted int64, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get todos: " + err.Error(),
			},
		})
		return
	}

	member, err := getDiscordUser(todos[0].User)
	tasks := []*discordgo.MessageEmbedField{}

	for _, t := range todos {
		tasks = append(tasks, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("ID: %d", t.Id),
			Value:  fmt.Sprintf("```%s```ㅤ", t.Task),
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title: "ㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤ",
		Color: 0xE0C460,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: member.User.AvatarURL("100"),
		},
		Description: fmt.Sprintf("**Total**\n╠ Uncompleted Tasks: **%d**\n╚ Completed Tasks: **%d**\nㅤ", len(todos), numCompleted),
		Fields:      tasks,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
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
	embed.SetTitle("TODOsㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤ\nㅤ")
	embed.SetColor(0xE0C460)
	embed.SetThumbnail("https://cdn-icons-png.flaticon.com/512/9717/9717679.png")
	for _, t := range todos {
		embed.AddField(fmt.Sprintf("ID: %d", t.Id), fmt.Sprintf("%s\n```%s```ㅤ", t.User, t.Task))
	}
	embed.SetFooter(fmt.Sprintf("Total: %d", len(todos)))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
