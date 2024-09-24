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
        sendErrorEmbed(s, i, err)
        return
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        sendErrorEmbed(s, i, err)
        return
    }

    type pods struct {
        Name string `json:"Name"`
        ResourceGroup string `json:"ResourceGroup"`
        ServerGUID string `json:"ServerGUID"`
    }

    var podList []pods
    err = json.Unmarshal(respBody, &podList)
    if err != nil {
        fmt.Println(err)
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Pods")
    embed.SetColor(0xffab40)

    podString := ""
    for i, pod := range podList {
        podString += fmt.Sprintf("%d. %s\n", i+1, pod.Name)
    }
    embed.AddField("Pods", podString)

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
        },
    })
}

func deletePod(s *discordgo.Session, i *discordgo.InteractionCreate) {
    podId := i.ApplicationCommandData().Options[0].Options[0].StringValue()

    resp, err := doAPIRequest("DELETE", config.KaminoDeleteEndpoint + "/" + podId, nil)
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

    resp, err := doAPIRequest("POST", config.KaminoBulkDeleteEndpoint, data)
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

    req, err := http.NewRequest(verb, config.KaminoURL + endpoint, body)
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
        Name: "ERROR",
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
        Name: "SUCCESS",
        Value: message,
    })
    
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
    })
}
