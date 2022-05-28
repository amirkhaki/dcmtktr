package main

import (
	"log"
	"os"
	"os/signal"
	"net/http"
	"syscall"
	"strings"
	"encoding/json"
	"bytes"
	"github.com/bwmarrin/discordgo"
)
var token string
var endpoint string
func readData(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(data), "\t\n ")
}
func init() {
	if len(os.Args) < 3{
		log.Fatal("usage: ./dcmtktr /path/to/token /path/to/endpoint")
	}
	token = readData(os.Args[1])
	endpoint = readData(os.Args[2])
}

func sendInviteToEndpoint(inviterID, targetUser, guildID string) error {
	data := struct{
		InID string `json:"inviter_id"`
		TID string `json:"target_user_id"`
		GID string `json:"guild_id"`
	}{
		InID: inviterID,
		TID: targetUser,
		GID: guildID,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil

}
func inviteCreate(s *discordgo.Session, i *discordgo.InviteCreate) {
	if err := sendInviteToEndpoint(i.Inviter.ID, i.TargetUser.ID, i.GuildID); err != nil {
		log.Println(err)
	}
	//invite created, send request to endpoint
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	dg.AddHandler(inviteCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildInvites
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
