package migrate

import (
	"golang-fave/engine/sqlw"
)

func Migrate_000000008(db *sqlw.DB, host string) error {
	// Changes
	if _, err := db.Exec(`ALTER TABLE blog_posts ADD COLUMN category INT(11) NOT NULL DEFAULT 1 AFTER alias;`); err != nil {
		return err
	}

	// Indexes
	if _, err := db.Exec(`ALTER TABLE blog_posts ADD KEY FK_blog_posts_category (category);`); err != nil {
		return err
	}

	// References
	if _, err := db.Exec(`
		ALTER TABLE blog_posts ADD CONSTRAINT FK_blog_posts_category
		FOREIGN KEY (category) REFERENCES blog_cats (id) ON DELETE RESTRICT;
	`); err != nil {
		return err
	}

	// Remove default
	if _, err := db.Exec(`ALTER TABLE blog_posts ALTER category DROP DEFAULT;`); err != nil {
		return err
	}

	return nil
}
