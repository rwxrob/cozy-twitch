package twitch

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fluffle/goirc/client"
	"github.com/rwxrob/cozy"
)

type Agent struct {
	User    string
	Pass    string
	Server  string
	Channel string
	Chan    chan any

	// api contains the internal http.Client used for all communications
	// with the Twitch REST API. A Keep-Alive is attempted and maintained
	// with a watchdog to protect dropped connections creating new ones
	// dynamically when dropped.
	api *http.Client

	// irc client connection.
	irc *client.Conn
}

var _ cozy.Agent = &Agent{} // TODO remove in final

func (a *Agent) Connect() error {
	a.Chan = make(chan any)

	cfg := client.NewConfig(a.User)
	cfg.Pass = a.Pass
	cfg.Server = "irc.twitch.tv:6667" // TODO only override if empty
	cfg.NewNick = func(n string) string { return n + "^" }
	a.irc = client.Client(cfg)

	// TODO change all the log.* to send messages on chan any

	a.irc.HandleFunc("connected", func(conn *client.Conn, line *client.Line) {
		conn.Join("#" + a.Channel)
		fmt.Printf("Connected to %s\n", line.Args[0])
	})

	a.irc.HandleFunc("disconnected", func(conn *client.Conn, line *client.Line) {
		log.Println("Disconnected from server")
		close(a.Chan)
	})

	a.irc.HandleFunc("join", func(conn *client.Conn, line *client.Line) {
		log.Printf("Joined channel %s\n", line.Args[0])
	})

	a.irc.HandleFunc("privmsg", func(conn *client.Conn, line *client.Line) {
		//nick := line.Nick
		message := line.Args[1]
		fmt.Println(line.Args[0], message)
		//if strings.HasPrefix(message, "!hello") {
		//conn.Privmsg("#"+a.Channel, "Hello, "+nick+"!")
		//}
	})

	err := a.irc.Connect()
	if err != nil {
		return err
	}

	// TODO remove and have caller do the waiting
	<-a.Chan

	return nil
}
