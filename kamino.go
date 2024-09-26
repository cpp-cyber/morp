package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

func NewKaminoClient() {
	cl = Login()
	if cl == nil {
		fmt.Println("Login failed")
	}
}

func Login() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := &http.Client{
		Jar: jar,
	}

	resp, err := doAPIRequest("POST", config.KaminoLoginEndpoint, map[string]any{
		"username": config.KaminoUser,
		"password": config.KaminoPass,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	client.Jar.SetCookies(resp.Request.URL, resp.Cookies())
	return client
}

func refreshLogin() {
	for {
		<-time.After(30 * time.Minute)
		NewKaminoClient()
	}
}

func getPods(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := doAPIRequest("GET", config.KaminoGetPodsEndpoint, nil)
	if err != nil {
		fmt.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get pods",
				Flags:   64,
			},
		})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	type pods struct {
		Name          string `json:"Name"`
		ResourceGroup string `json:"ResourceGroup"`
		ServerGUID    string `json:"ServerGUID"`
	}

	var podList []pods
	err = json.Unmarshal(respBody, &podList)
	if err != nil {
		fmt.Println(err)
	}

	embed := embed.NewEmbed()
	embed.SetTitle("ㅤ")
	embed.SetColor(0xE0C460)
	embed.SetAuthor("Kamino", "https://kamino.calpolyswift.org/img/bruharmy.0e3831f1.png")

	podString := ""
	for i, pod := range podList {
		podString += fmt.Sprintf("%d. `%s`\n", i+1, pod.Name)
	}

	embed.AddField("Pods", podString)
	embed.SetFooter(fmt.Sprintf("ㅤ\nTotal: %d", len(podList)))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func deletePod(s *discordgo.Session, i *discordgo.InteractionCreate) {
	podId := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	resp, err := doAPIRequest("DELETE", config.KaminoDeleteEndpoint+"/"+podId, nil)
	if err != nil {
		sendErrorEmbed(s, i, err)
		return
	}
	defer resp.Body.Close()

	message := fmt.Sprintf("Pod %s deleted", podId)
	sendSuccessEmbed(s, i, message)
}

func bulkDeletePods(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := make(map[string]any)
	filters := strings.ReplaceAll(i.ApplicationCommandData().Options[0].Options[0].StringValue(), " ", "")
	filtersList := strings.Split(filters, ",")
	data["filters"] = filtersList

	resp, err := doAPIRequest("DELETE", config.KaminoBulkDeleteEndpoint, data)
	if err != nil || resp == nil {
		sendErrorEmbed(s, i, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			sendErrorEmbed(s, i, err)
			return
		}
		msg := make(map[string]string)
		err = json.Unmarshal(message, &msg)
		if err != nil {
			sendErrorEmbed(s, i, err)
			return
		}

		failedPodsError := errors.New(msg["error"])
		sendErrorEmbed(s, i, failedPodsError)
		return
	}

	message := "Pods deleted"
	sendSuccessEmbed(s, i, message)
}

func competitionClone(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options[0].Options) < 2 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please provide a template and number of teams.",
			},
		})
		return
	}

	template := i.ApplicationCommandData().Options[0].Options[0].StringValue()
	count := i.ApplicationCommandData().Options[0].Options[1].IntValue()
	data := map[string]any{
		"template": template,
		"count":    count,
	}

	resp, err := doAPIRequest("POST", config.KaminoCompetitionCloneEndpoint, data)
	if err != nil {
		sendErrorEmbed(s, i, err)
		return
	}
	defer resp.Body.Close()

	type compCloneResponse struct {
		Message string            `json:"message"`
		Users   map[string]string `json:"users"`
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		sendErrorEmbed(s, i, err)
		return
	}

	var compCloneResp compCloneResponse
	err = json.Unmarshal(respBody, &compCloneResp)
	if err != nil {
		sendErrorEmbed(s, i, err)
		return
	}

	embed := embed.NewEmbed()
	embed.SetTitle("ㅤ\nCompetition Clone")
	embed.SetColor(0xE0C460)
	embed.SetDescription(fmt.Sprintf("%s\nㅤ", compCloneResp.Message))
	embed.SetAuthor("Kamino", "https://kamino.calpolyswift.org/img/bruharmy.0e3831f1.png")

	usersString := "```\n"
	for user, password := range compCloneResp.Users {
		usersString += fmt.Sprintf("%s: %s\n", user, password)
	}
	embed.AddField("Users", usersString+"```")

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	})
}

func refreshTemplates(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp, err := doAPIRequest("POST", config.KaminoRefreshTemplatesEndpoint, nil)
	if err != nil {
		sendErrorEmbed(s, i, err)
		return
	}
	defer resp.Body.Close()

	message := "Templates refreshed"
	sendSuccessEmbed(s, i, message)
}

func doAPIRequest(verb, endpoint string, data map[string]any) (*http.Response, error) {
	var body io.Reader
	if data != nil {
		jsonBody, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(verb, config.KaminoURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := cl.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	return resp, nil
}

func sendErrorEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	embed := embed.NewEmbed()
	embed.SetColor(0xff0000)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "ERROR",
		Value: err.Error(),
	})

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	})
}

func sendSuccessEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	embed := embed.NewEmbed()
	embed.SetColor(0x00ff00)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "SUCCESS",
		Value: message,
	})

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	})
}
