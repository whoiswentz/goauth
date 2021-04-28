package posts

import "github.com/whoiswentz/goauth/database"

func create(db *database.Database, p Post) (*Post, error) {
	stmt, err := db.Db.Prepare("INSERT INTO posts (title, content) VALUES (?, ?)")
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(p.Title, p.Content)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.Id = id

	return &p, nil
}

func list(db *database.Database) ([]Post, error) {
	rows, err := db.Db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Content); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
