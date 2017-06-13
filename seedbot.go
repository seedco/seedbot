package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
	seed "github.com/seedco/seed-go"
)

type Seedbot struct {
	rtm        *slack.RTM
	seedClient *seed.Client
	username   string
	logger     *log.Logger
}

func New(slackToken, seedApiToken string) *Seedbot {
	api := slack.New(slackToken)
	logger := log.New(os.Stdout, "seed-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	if os.Getenv("SEEDBOT_DEBUG") == "true" {
		api.SetDebug(true)
	}
	return &Seedbot{
		rtm:        api.NewRTM(),
		seedClient: seed.New(seedApiToken),
		logger:     logger,
	}
}

func (s *Seedbot) Run() {
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessage(ev); err != nil {
				s.logger.Printf("error when handling message: %v", err)
			}
		case *slack.PresenceChangeEvent:
			if ev.Presence == "active" {
				s.username = ev.User
			}
		case *slack.LatencyReport:
			s.logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			s.logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			s.logger.Print("Invalid credentials")
			return
		default:
			// Ignore other events..
		}
	}
}

func (s *Seedbot) handleMessage(m *slack.MessageEvent) error {
	var err error
	if !strings.HasPrefix(m.Text, fmt.Sprintf("<@%s>", s.username)) {
		return nil
	}

	command := strings.TrimSpace(strings.TrimLeft(m.Text, fmt.Sprintf("<@%s>", s.username)))
	commandSplit := strings.Split(command, " ")
	switch commandSplit[0] {
	case "transactions":
		tReq := seed.TransactionsRequest{
			Client: s.seedClient,
		}

		if len(commandSplit) > 1 {
			var from, to time.Time
			if from, to, err = ProcessDate(strings.Join(commandSplit[1:], " ")); err != nil {
				return fmt.Errorf("error with command %v", err)
			}
			tReq.From = from
			tReq.To = to
		}

		iterator := tReq.Iterator()
		txs, err := iterator.Next()
		if err != nil {
			return fmt.Errorf("error when getting transactions: %v", err)
		}
		if len(txs) == 0 {
			return nil
		}

		tString := constructTransactionMessage(txs)
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(tString, m.Channel))
	case "balance":
		bReq := seed.BalanceRequest{
			Client: s.seedClient,
		}
		b, err := bReq.Get()
		if err != nil {
			return fmt.Errorf("error when geting balance: %v", err)
		}
		bString := constructBalanceMessage(b)
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(bString, m.Channel))
	case "help":
		helpText :=
			`
		I am Seedbot. I can answer questions about your transactions and balances. You can ask me:

help
transactions 6/12/2017
transactions June 12, 2017
transactions June 12
transactions 2017
transactions today
transactions yesterday
transactions this (week|month|year)
transactions last (week|month|year)

balance`
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage(helpText, m.Channel))

	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("I'm sorry I don't know how to do that. I'm still learning", m.Channel))

	}
	return nil
}

func constructTransactionMessage(txs []seed.Transaction) string {
	s := bytes.NewBufferString("```")
	var maxLength int
	for _, tx := range txs {
		if len(tx.Description) > maxLength {
			maxLength = len(tx.Description)
		}
	}

	length := maxLength + 8

	for _, tx := range txs {
		s.WriteString(fmt.Sprintf(fmt.Sprintf("%%s    %%s%%%ds", length-len(tx.Description)), tx.Date.Format("01/02/2006"), tx.Description, CentsToDollarStringWithCommas(tx.Amount)))
		s.WriteString("\n")
	}
	message := strings.TrimRight(s.String(), "\n")
	message = fmt.Sprintf("%s```", message)

	return message
}

func constructBalanceMessage(b seed.Balance) string {
	s := bytes.NewBufferString("```")

	s.WriteString(fmt.Sprintf("Posted: %16s\n", "$"+CentsToDollarStringWithCommas(b.Settled)))
	s.WriteString(fmt.Sprintf("Pending Debits: %7s\n", "$"+CentsToDollarStringWithCommas(int64(b.PendingDebits))))
	s.WriteString(fmt.Sprintf("Pending Credits: %3s\n", "$"+CentsToDollarStringWithCommas(int64(b.PendingCredits))))
	s.WriteString(fmt.Sprintf("Lock Box: %12s\n", "$"+CentsToDollarStringWithCommas(int64(b.Lockbox))))
	s.WriteString(fmt.Sprintf("Available %14s\n", "$"+CentsToDollarStringWithCommas(int64(b.TotalAvailable))))

	return fmt.Sprintf("%s```", s.String())
}
