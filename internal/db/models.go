package db

import (
	"fmt"
)

func GetOrInsertID(table, column, value string) int {
	Mutex.Lock()
	defer Mutex.Unlock()

	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE %s = ?", table, column)
	err := DB.QueryRow(query, value).Scan(&id)
	if err == nil {
		return id
	}

	insert := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", table, column)
	res, err := DB.Exec(insert, value)
	if err != nil {
		return 0
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0
	}
	return int(lastID)
}
