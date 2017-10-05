package db

import (
	"database/sql"
	"reflect"

	"github.com/tomogoma/authms/model"
	"github.com/tomogoma/go-commons/errors"
)

// InsertUserType inserts into the database returning calculated values.
func (r *Roach) InsertUserName(userID, username string) (*model.Username, error) {
	return insertUserName(r.db, userID, username)
}

// InsertUserType inserts through tx returning calculated values.
func (r *Roach) InsertUserNameAtomic(tx *sql.Tx, userID, username string) (*model.Username, error) {
	return insertUserName(tx, userID, username)
}

func (r *Roach) UpdateUsername(userID, username string) (*model.Username, error) {
	return nil, errors.NewNotImplemented()
}

func insertUserName(tx inserter, userID, username string) (*model.Username, error) {
	if tx == nil || reflect.ValueOf(tx).IsNil() {
		return nil, errorNilTx
	}
	un := model.Username{UserID: userID, Value: username}
	insCols := ColDesc(ColUserID, ColUserName, ColUpdateDate)
	retCols := ColDesc(ColID, ColCreateDate, ColUpdateDate)
	q := `
	INSERT INTO ` + TblUserNameIDs + ` (` + insCols + `)
		VALUES ($1,$2,CURRENT_TIMESTAMP)
		RETURNING ` + retCols
	err := tx.QueryRow(q, userID, username).Scan(&un.ID, &un.CreateDate, &un.UpdateDate)
	if err != nil {
		return nil, err
	}
	return &un, nil
}
