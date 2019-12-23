package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

func Migrate_000000014(ctx context.Context, db *sqlw.DB, host string) error {
	// Table: notify_mail
	if _, err := db.Exec(
		ctx,
		`CREATE TABLE notify_mail (
			id int(11) NOT NULL AUTO_INCREMENT COMMENT 'AI',
			email varchar(255) NOT NULL COMMENT 'Email address',
			subject varchar(800) NOT NULL COMMENT 'Email subject',
			message text NOT NULL COMMENT 'Email body',
			error text NOT NULL COMMENT 'Send error message',
			datetime datetime NOT NULL COMMENT 'Creation date/time',
			status int(1) NOT NULL COMMENT 'Sending status',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,
	); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(ctx, `ALTER TABLE notify_mail ADD KEY status (status);`); err != nil {
		return err
	}

	return nil
}
