package web

import (
    "github.com/jackc/pgx" 
    log15 "gopkg.in/inconshreveable/log15.v2"

    "github.com/SCKelemen/Cassius/data"
    "github.com/SCKelemen/Cassius/mail"
)

type environment struct { 
  user   *data.User 
  pool   *pgx.ConnPool 
  mailer mail.Mailer 
  logger log15.Logger 
} 