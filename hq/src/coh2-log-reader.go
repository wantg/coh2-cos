package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
)

type player struct {
	ProfileID      string `json:"profileID"`
	SteamID        string `json:"steamID"`
	NickName       string `json:"nickName"`
	Ranking        int    `json:"ranking"`
	Side           int    `json:"side"`
	Faction        string `json:"faction"`
	InCurrentMatch bool   `json:"inCurrentMatch"`
	Slot           int    `json:"slot"`
}

var players = map[string]player{}

func readLog() {
	// docPath, _ := filepath.Abs(viper.GetString("coh2.doc-path"))
	docPath := viper.GetString("coh2.doc-path")
	sysVarList := []string{"USERPROFILE"}
	for _, v := range sysVarList {
		_v := "%" + v + "%"
		if strings.Contains(docPath, _v) {
			docPath = strings.ReplaceAll(docPath, _v, os.Getenv(v))
		}
	}
	warningsLogPath := filepath.Join(docPath, "warnings.log")
	if _, err := os.Stat(warningsLogPath); os.IsNotExist(err) {
		return
	}
	go func() {
		t, _ := tail.TailFile(warningsLogPath, tail.Config{Follow: true, Poll: true})
		for line := range t.Lines {
			l := strings.TrimSpace(line.Text)
			if strings.Contains(l, "WorldwideAutomatchService::OnStartComplete - detected successful game start") {
				for profileID, player := range players {
					player.InCurrentMatch = false
					players[profileID] = player
				}
			} else if strings.Contains(l, "Match Started") {
				tmp := strings.Split(l, "000:")
				tmp = strings.Split(tmp[1], " ")
				profileIDInt, _ := strconv.ParseInt(tmp[0], 16, 64)
				profileID := strconv.FormatInt(profileIDInt, 10)
				tmp = strings.Split(l, "steam/")
				steamID := strings.Split(tmp[1], "]")[0]
				tmp = strings.Split(l, "=")
				ranking := strings.TrimSpace(tmp[len(tmp)-1])
				_ranking, _ := strconv.Atoi(ranking)
				players[profileID] = player{ProfileID: profileID, SteamID: steamID, Ranking: _ranking}
			} else if strings.Contains(l, "Human Player: ") {
				tmp := strings.Split(l, "Human Player: ")
				tmp = strings.Split(tmp[1], " ")
				slot := tmp[0]
				profileID := tmp[len(tmp)-3]
				side := tmp[len(tmp)-2]
				faction := tmp[len(tmp)-1]
				nickName := strings.Join(tmp[1:len(tmp)-3], " ")
				for _profileID, player := range players {
					if _profileID == profileID {
						_slot, _ := strconv.Atoi(slot)
						_side, _ := strconv.Atoi(side)
						player.Slot = _slot
						player.InCurrentMatch = true
						player.NickName = nickName
						player.Side = _side
						player.Faction = faction
						players[profileID] = player
						break
					}
				}
			}
		}
	}()
}
