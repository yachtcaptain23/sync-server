package datastore

import (
	"time"
)

// Client is a struct used to represent records in clients table.
type Client struct {
	ID       string `db:"id"`
	Token    string `db:"token"`
	ExpireAt int64  `db:"expire_at"`
}

// InsertClient create and insert a new client into clients table.
func (pg *Postgres) InsertClient(id string, token string, expireAt int64) error {
	client := Client{ID: id, Token: token, ExpireAt: expireAt}
	stmt := "INSERT INTO clients(id, token, expire_at) VALUES(:id, :token, :expire_at)" +
		"ON CONFLICT (id) DO UPDATE SET token = :token, expire_at = :expire_at"

	_, err := pg.NamedExec(stmt, client)
	return err
}

// GetClient queries the clients table using token, return the clientID if
// the token is valid, otherwise, return empty string.
func (pg *Postgres) GetClient(token string) string {
	var clientID string
	err := pg.Get(&clientID, "SELECT id FROM clients WHERE token = $1 AND expire_at > $2",
		token, time.Now().Unix())
	if err != nil {
		return ""
	}

	return clientID
}
