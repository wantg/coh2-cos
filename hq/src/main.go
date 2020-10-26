package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"coh2-cos/templates"
)

//go:generate go run ./cmd/html-templates.go

const contentTypeHTML = "text/html;charset=UTF-8"
const contentTypeJSON = "application/json;charset=utf-8"

func loadConfig() {
	configPtr := flag.String("c", "", "config file path")
	flag.Parse()
	configPath, _ := filepath.Abs(*configPtr)
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	loadConfig()
	readLog()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	rootPath, _ := os.Executable()
	r.Static("/assets", filepath.Join(filepath.Dir(rootPath), "assets"))

	r.GET("/match", func(c *gin.Context) {
		c.Data(http.StatusOK, contentTypeHTML, templates.Data["match.html"])
	})

	r.GET("/match/players", func(c *gin.Context) {
		_players := []player{}
		for _, player := range players {
			if player.InCurrentMatch {
				_players = append(_players, player)
			}
		}
		sort.Slice(_players, func(a, b int) bool {
			if _players[a].Side == _players[b].Side {
				return _players[a].Slot < _players[b].Slot
			} else {
				return _players[a].Side < _players[b].Side
			}
		})
		c.PureJSON(http.StatusOK, _players)
	})

	r.GET("/player/:steamID/avatar", func(c *gin.Context) {
		steamID := c.Param("steamID")
		resp, err := http.Get(`https://coh2-api.reliclink.com/community/external/proxysteamuserrequest?title=coh2&profileNames=["/steam/` + steamID + `"]&request=/ISteamUser/GetPlayerSummaries/v0002/`)
		if err != nil {
			c.PureJSON(http.StatusOK, gin.H{"result": 1})
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.PureJSON(http.StatusOK, gin.H{"result": 2})
			return
		}
		metadata := gjson.ParseBytes(body)
		players := metadata.Get("steamResults.response.players").Array()
		if len(players) == 0 {
			c.PureJSON(http.StatusOK, gin.H{"result": 3})
			return
		}
		c.Redirect(http.StatusMovedPermanently, players[0].Get("avatarfull").String())
	})

	r.GET("/player/:steamID/ranking", func(c *gin.Context) {
		steamID := c.Param("steamID")
		resp, err := http.Get(`https://coh2-api.reliclink.com/community/leaderboard/GetPersonalStat?title=coh2&profile_names=["/steam/` + steamID + `"]`)
		if err != nil {
			c.PureJSON(http.StatusOK, gin.H{"result": 1})
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.PureJSON(http.StatusOK, gin.H{"result": 2})
			return
		}
		metadata := gjson.ParseBytes(body)
		statGroups := metadata.Get("statGroups").Array()
		leaderboardStats := metadata.Get("leaderboardStats").Array()
		var groupID int64
		for _, g := range statGroups {
			members := g.Get("members").Array()
			if len(members) == 1 {
				groupID = g.Get("id").Int()
				break
			}
		}
		rankData := map[int64]interface{}{}
		for _, s := range leaderboardStats {
			if s.Get("statGroup_id").Int() == groupID {
				rankData[s.Get("leaderboard_id").Int()] = s.Value()
			}
		}
		c.PureJSON(http.StatusOK, gin.H{"result": 0, "data": rankData})
	})

	listen := viper.GetString("app.listen")
	logrus.Info("listen on http://" + listen)
	r.Run(listen)
}
