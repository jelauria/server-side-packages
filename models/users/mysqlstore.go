package users

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type SQLStore struct{ DB *sql.DB }

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db}
}

//GetByID returns the User with the given ID
func (ss *SQLStore) GetByID(id int64) (*User, error) {
	rows, err1 := ss.DB.Query("select id,email,pass_hash,username,first_name,last_name,photo_url from users where id=?", id)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

//GetByEmail returns the User with the given email
func (ss *SQLStore) GetByEmail(email string) (*User, error) {
	rows, err1 := ss.DB.Query("select id,email,pass_hash,username,first_name,last_name,photo_url from users where email=?", email)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

//GetByUserName returns the User with the given Username
func (ss *SQLStore) GetByUserName(username string) (*User, error) {
	rows, err1 := ss.DB.Query("select id,email,pass_hash,username,first_name,last_name,photo_url from users where username=?", username)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.PassHash, &user.UserName,
			&user.FirstName, &user.LastName, &user.PhotoURL); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

//Insert inserts the user into the database, and returns
//the newly-inserted User, complete with the DBMS-assigned ID
func (ss *SQLStore) Insert(user *User) (*User, error) {
	insq := "insert into users(email, pass_hash, username, first_name, last_name, photo_url) values (?,?,?,?,?,?)"
	res, err1 := ss.DB.Exec(insq, user.Email, user.PassHash, user.UserName, user.FirstName, user.LastName, user.PhotoURL)
	if err1 != nil {
		return nil, err1
	}
	id, err2 := res.LastInsertId()
	if err2 != nil {
		return nil, err2
	}
	user.ID = id
	return user, nil
}

//Update applies UserUpdates to the given user ID
//and returns the newly-updated user
func (ss *SQLStore) Update(id int64, updates *Updates) (*User, error) {
	insq := "update users set first_name=?, last_name=? where id=?"
	_, err1 := ss.DB.Exec(insq, updates.FirstName, updates.LastName, id)
	if err1 != nil {
		return nil, err1
	}
	user, err := ss.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Delete deletes the user with the given ID
func (ss *SQLStore) Delete(id int64) error {
	insq := "delete from users where id=?"
	_, err := ss.DB.Exec(insq, id)
	return err
}
