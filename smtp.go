package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang-fave/engine/mysqlpool"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/wrapper/config"
	"golang-fave/utils"
)

func smtp_send(host, port, user, pass, subject, msg string, receivers []string) error {
	return utils.SMTPSend(host, port, user, pass, subject, msg, receivers)
}

func smtp_prepare(db *sqlw.DB, conf *config.Config) {
	rows, err := db.Query(
		`SELECT
			id,
			email,
			subject,
			message
		FROM
			notify_mail
		WHERE
			status = 2
		ORDER BY
			id ASC
		;`,
	)
	if err == nil {
		defer rows.Close()
		values := make([]string, 4)
		scan := make([]interface{}, len(values))
		for i := range values {
			scan[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scan...)
			if err == nil {
				if _, err := db.Exec(
					`UPDATE notify_mail SET status = 3 WHERE id = ?;`,
					utils.StrToInt(string(values[0])),
				); err == nil {
					go func(db *sqlw.DB, conf *config.Config, id int, subject, msg string, receivers []string) {
						if err := smtp_send(
							(*conf).SMTP.Host,
							utils.IntToStr((*conf).SMTP.Port),
							(*conf).SMTP.Login,
							(*conf).SMTP.Password,
							subject,
							msg,
							receivers,
						); err == nil {
							if _, err := db.Exec(
								`UPDATE notify_mail SET status = 1 WHERE id = ?;`,
								id,
							); err != nil {
								fmt.Printf("Smtp send error (sql, success): %v\n", err)
							}
						} else {
							if _, err := db.Exec(
								`UPDATE notify_mail SET error = ?, status = 0 WHERE id = ?;`,
								err.Error(),
								id,
							); err != nil {
								fmt.Printf("Smtp send error (sql, error): %v\n", err)
							}
						}
					}(
						db,
						conf,
						utils.StrToInt(string(values[0])),
						html.EscapeString(string(values[2])),
						string(values[3]),
						[]string{html.EscapeString(string(values[1]))},
					)
				} else {
					fmt.Printf("Smtp send error (sql, update): %v\n", err)
				}
			}
		}
	}
}

func smtp_process(dir, host string, mp *mysqlpool.MySqlPool) {
	db := mp.Get(host)
	if db != nil {
		conf := config.ConfigNew()
		if err := conf.ConfigRead(strings.Join([]string{dir, "config", "config.json"}, string(os.PathSeparator))); err == nil {
			if !((*conf).SMTP.Host == "" || (*conf).SMTP.Login == "" && (*conf).SMTP.Password == "") {
				if err := db.Ping(); err == nil {
					smtp_prepare(db, conf)
				}
			}
		} else {
			fmt.Printf("Smtp error (config): %v\n", err)
		}
	}
}

func smtp_loop(www_dir string, stop chan bool, mp *mysqlpool.MySqlPool) {
	dirs, err := ioutil.ReadDir(www_dir)
	if err == nil {
		for _, dir := range dirs {
			select {
			case <-stop:
				break
			default:
				if mp != nil {
					target_dir := strings.Join([]string{www_dir, dir.Name()}, string(os.PathSeparator))
					if utils.IsDirExists(target_dir) {
						smtp_process(target_dir, dir.Name(), mp)
					}
				}
			}
		}
	}
}

func smtp_start(www_dir string, mp *mysqlpool.MySqlPool) (chan bool, chan bool) {
	ch := make(chan bool)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(5 * time.Second):
				// Run every 5 seconds
				smtp_loop(www_dir, stop, mp)
			case <-ch:
				ch <- true
				return
			}
		}
	}()
	return ch, stop
}

func smtp_stop(ch, stop chan bool) {
	for {
		select {
		case stop <- true:
		case ch <- true:
			<-ch
			return
		case <-time.After(3 * time.Second):
			fmt.Println("Smtp error: force exit by timeout after 3 seconds")
			return
		}
	}
}
