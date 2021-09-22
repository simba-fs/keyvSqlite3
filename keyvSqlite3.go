package keyvsqlite3

import (
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/simba-fs/keyv"
)

const createTable = `CREATE TABLE IF NOT EXISTS keyv (
	key string PRIMARY KEY UNIQUE,
	value string
);`

type keyvData struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type adapter struct {
	db *sqlx.DB
}

func (a adapter) Has(key string) bool {
	data := keyvData{}
	err := a.db.Get(&data, `SELECT * FROM keyv WHERE key = $1`, key)
	if err != nil {
		return false
	}

	if data.Key == "" && data.Value == "" {
		return false
	}
	return true
}

func (a adapter) Get(key string) (string, error) {
	data := keyvData{}
	err := a.db.Get(&data, `SELECT * FROM keyv WHERE key = $1`, key)
	if err != nil {
		return "", err
	}

	return data.Value, nil
}

func (a adapter) Set(key string, val string) error {
	data := keyvData{
		Key:   key,
		Value: val,
	}

	_, err := a.db.NamedExec(`INSERT OR REPLACE INTO keyv VALUES (:key, :value)`, data)
	return err
}

func (a adapter) Remove(key string) error {
	_, err := a.db.Exec(`DELETE FROM keyv WHERE key = $1`, key)

	return err
}

func (a adapter) Clear(prefix string) error {
	// get all keys
	keys, err := a.Keys()
	if err != nil {
		return err
	}

	stmt, err := a.db.Preparex(`DELETE FROM keyv WHERE key=$1`)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			_, err := stmt.Exec(key)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (a adapter) Keys() ([]string, error) {
	data := []keyvData{}

	err := a.db.Select(&data, `SELECT key FROM keyv`)
	if err != nil {
		return make([]string, 0), err
	}

	keys := make([]string, len(data))
	for index, value := range data {
		keys[index] = value.Key
	}

	return keys, nil
}

type adapterNewer struct{}

func (a adapterNewer) Connect(uri string) (keyv.Adapter, error) {
	uri = strings.SplitN(uri, "://", 2)[1]

	db, err := sqlx.Connect("sqlite3", uri)
	if err != nil {
		return nil, err
	}

	// create table
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	adapter := adapter{
		db: db,
	}

	return adapter, nil
}

func init() {
	adapterNewer := adapterNewer{}
	keyv.Register("sqlite3", adapterNewer)
}
