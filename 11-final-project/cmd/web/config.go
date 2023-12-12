package main

import (
	"database/sql"
	"final-project/cmd/data"
	"github.com/alexedwards/scs/v2"
	"log"
	"sync"
)

type Config struct {
	Session       *scs.SessionManager
	DB            *sql.DB
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Wait          *sync.WaitGroup
	Models        data.Models
	Mailer        Mail
	ErrorChan     chan error
	ErrorChanDone chan bool
}
