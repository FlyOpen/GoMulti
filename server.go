package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"net/http"
	"strconv"
	"time"
)

const WIDTH = 300
const HEIGHT = 300

type Player struct {
	Id string
	X  int64
	Y  int64
	Vx int64
	Vy int64
}

func newPlayer() *Player {

	rb := make([]byte, 32)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}

	id := base64.URLEncoding.EncodeToString(rb)

	return &Player{Id: id, X: 0, Y: 0, Vx: 1, Vy: 1}
}

type Game struct {
	players map[*melody.Session]*Player
}

func newGame() *Game {
	return &Game{players: make(map[*melody.Session]*Player)}
}

func (this *Game) AddPlayer(s *melody.Session) {
	p := newPlayer()
	this.players[s] = p

	for z, _ := range this.players {
		z.Write([]byte("add:" + p.Id + ":" + strconv.FormatInt(p.X, 10) + "," + strconv.FormatInt(p.Y, 10)))
	}
}

func (this *Game) RemovePlayer(s *melody.Session) {
	delete(this.players, s)
}

func (this *Game) run() {
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			<-ticker.C
			for _, p := range this.players {
				p.X += p.Vx
				p.Y += p.Vy

				if p.X > WIDTH {
					p.X = WIDTH
					p.Vx = 0
				}
				if p.Y > HEIGHT {
					p.Y = HEIGHT
					p.Vy = 0
				}

				for s, _ := range this.players {
					fmt.Println("player:" + p.Id + ":" + strconv.FormatInt(p.X, 10) + "," + strconv.FormatInt(p.Y, 10))

					s.Write([]byte("player:" + p.Id + ":" + strconv.FormatInt(p.X, 10) + "," + strconv.FormatInt(p.Y, 10)))
				}

			}

		}
	}()
}

func main() {
	r := gin.New()
	m := melody.New()

	size := 65536
	m.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  size,
		WriteBufferSize: size,
	}
	m.Config.MaxMessageSize = int64(size)
	m.Config.MessageBufferSize = 2048

	game := newGame()

	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		fmt.Println("Debut de connection s")
		m.HandleRequest(c.Writer, c.Request)
	})

	//var mutex sync.Mutex

	m.HandleConnect(func(s *melody.Session) {
		game.AddPlayer(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("message")
		fmt.Println(msg)

	})

	m.HandleDisconnect(func(s *melody.Session) {
		game.RemovePlayer(s)
	})

	game.run()

	r.Run(":5000")
}
