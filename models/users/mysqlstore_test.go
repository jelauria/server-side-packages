package users

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestGetByID is a test function for the SQLStore's GetByID
func TestGetByID(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		idToGet      int64
		expectError  bool
	}{
		{
			"User Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"User Not Found",
			&User{},
			2,
			true,
		},
		{
			"User With Large ID Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1234567890,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,pass_hash,username,first_name,last_name,photo_url from users where id=?"

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnError(ErrUserNotFound)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.idToGet).WillReturnRows(row)

			// Test GetByID()
			user, err := mainSQLStore.GetByID(c.idToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByEmail is a test function for the SQLStore's GetByEmail
func TestGetByEmail(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		emailToGet   string
		expectError  bool
	}{
		{
			"User with Plain Email Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"test@test.com",
			false,
		},
		{
			"User Not Found",
			&User{},
			"test@test.com",
			true,
		},
		{
			"User With Mixed Case Email Found",
			&User{
				1234567890,
				"Testing2@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"Testing2@test.com",
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,pass_hash,username,first_name,last_name,photo_url from users where email=?"

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnError(ErrUserNotFound)

			// Test GetByEmail()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.emailToGet).WillReturnRows(row)

			// Test GetByEmail()
			user, err := mainSQLStore.GetByEmail(c.emailToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestGetByUserName is a test function for the SQLStore's GetByUserName
func TestGetByUserName(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		expectedUser *User
		unameToGet   string
		expectError  bool
	}{
		{
			"User with Plain Username Found",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			"username",
			false,
		},
		{
			"User Not Found",
			&User{},
			"username2",
			true,
		},
		{
			"User With Mixed Case Username Found",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"uSeRnaMe1",
				"firstname",
				"lastname",
				"photourl",
			},
			"uSeRnaMe1",
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create an expected row to the mock DB
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "select id,email,pass_hash,username,first_name,last_name,photo_url from users where username=?"

		if c.expectError {
			// Set up expected query that will expect an error
			mock.ExpectQuery(query).WithArgs(c.unameToGet).WillReturnError(ErrUserNotFound)

			// Test GetByUserName()
			user, err := mainSQLStore.GetByUserName(c.unameToGet)
			if user != nil || err == nil {
				t.Errorf("Expected error [%v] but got [%v] instead", ErrUserNotFound, err)
			}
		} else {
			// Set up an expected query with the expected row from the mock DB
			mock.ExpectQuery(query).WithArgs(c.unameToGet).WillReturnRows(row)

			// Test GetByUserName()
			user, err := mainSQLStore.GetByUserName(c.unameToGet)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(user, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

// TestUpdate is a test function for the SQLStore's Update
func TestUpdate(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		ogUser       *User
		expectedUser *User
		givenUpdates *Updates
		submittedID  int64
		expectError  bool
	}{
		{
			"Normal Case",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname1",
				"lastname1",
				"photourl",
			},
			&Updates{
				"firstname1",
				"lastname1",
			},
			1,
			false,
		},
		{
			"Incorrect User ID",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			&User{},
			&Updates{
				"firstname1",
				"lastname1",
			},
			2,
			true,
		},
	}

	for _, c := range cases {
		//Generate Existing database to test updates
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		// Create row detailing the original user
		mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.ogUser.ID,
			c.ogUser.Email,
			c.ogUser.PassHash,
			c.ogUser.UserName,
			c.ogUser.FirstName,
			c.ogUser.LastName,
			c.ogUser.PhotoURL,
		)
		// Create row detailing the updated user
		row := mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.expectedUser.ID,
			c.expectedUser.Email,
			c.expectedUser.PassHash,
			c.expectedUser.UserName,
			c.expectedUser.FirstName,
			c.expectedUser.LastName,
			c.expectedUser.PhotoURL,
		)

		query := "update users set first_name=?, last_name=? where id=?"

		if c.expectError {
			mock.ExpectPrepare("update").ExpectExec().WithArgs(
				c.givenUpdates.FirstName,
				c.givenUpdates.LastName,
				c.submittedID,
			)
			db.Prepare(query)
			// Test Update()
			result, err := mainSQLStore.Update(c.submittedID, c.givenUpdates)
			if result != nil || err == nil {
				t.Errorf("Expected error but did not get one.")
			}
		} else {
			mock.ExpectPrepare("update").ExpectExec().WithArgs(
				c.givenUpdates.FirstName,
				c.givenUpdates.LastName,
				c.submittedID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("select id,email,pass_hash,username,first_name,last_name,photo_url from users where id=?").WithArgs(c.submittedID).WillReturnRows(row)
			db.Prepare(query)
			// Test Update()
			result, err := mainSQLStore.Update(c.submittedID, c.givenUpdates)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(result, c.expectedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

//TestInsert is a test function for the SQLStore's Insert
func TestInsert(t *testing.T) {
	// Create a slice of test cases
	cases := []struct {
		name         string
		insertedUser *User
		expectError  bool
	}{
		{
			"Normal User",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			false,
		},
		{
			"Blank User",
			&User{},
			true,
		},
		{
			"User With Mixed Case Params",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"uSeRnaMe1",
				"firstname",
				"lastname",
				"photourl",
			},
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		query := "insert into users(email, pass_hash, username, first_name, last_name, photo_url) values (?,?,?,?,?,?)"

		if c.expectError {
			mock.ExpectPrepare("insert into").ExpectExec().WithArgs(
				c.insertedUser.Email,
				c.insertedUser.PassHash,
				c.insertedUser.UserName,
				c.insertedUser.FirstName,
				c.insertedUser.LastName,
				c.insertedUser.PhotoURL,
			)
			db.Prepare(query)
			// Test Insert()
			result, err := mainSQLStore.Insert(c.insertedUser)
			if result != nil || err == nil {
				t.Errorf("Expected error but did not get one.")
			}
		} else {
			mock.ExpectPrepare("insert into").ExpectExec().WithArgs(
				c.insertedUser.Email,
				c.insertedUser.PassHash,
				c.insertedUser.UserName,
				c.insertedUser.FirstName,
				c.insertedUser.LastName,
				c.insertedUser.PhotoURL,
			).WillReturnResult(sqlmock.NewResult(1, 1))
			db.Prepare(query)
			// Test Insert()
			result, err := mainSQLStore.Insert(c.insertedUser)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
			if !reflect.DeepEqual(result, c.insertedUser) {
				t.Errorf("Error, invalid match in test [%s]", c.name)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}

//TestDelete is a test function for the SQLStore's Delete
func TestDelete(t *testing.T) {

	// Create a slice of test cases
	cases := []struct {
		name        string
		ogUser      *User
		submittedID int64
		expectError bool
	}{
		{
			"Normal User",
			&User{
				1,
				"test@test.com",
				[]byte("passhash123"),
				"username",
				"firstname",
				"lastname",
				"photourl",
			},
			1,
			false,
		},
		{
			"Incorrect User ID",
			&User{
				1234567890,
				"test@test.com",
				[]byte("passhash123"),
				"uSeRnaMe1",
				"firstname",
				"lastname",
				"photourl",
			},
			2,
			false,
		},
	}

	for _, c := range cases {
		// Create a new mock database for each case
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("There was a problem opening a database connection: [%v]", err)
		}
		defer db.Close()

		mainSQLStore := NewSQLStore(db)

		mock.NewRows([]string{
			"ID",
			"Email",
			"PassHash",
			"UserName",
			"FirstName",
			"LastName",
			"PhotoURL"},
		).AddRow(
			c.ogUser.ID,
			c.ogUser.Email,
			c.ogUser.PassHash,
			c.ogUser.UserName,
			c.ogUser.FirstName,
			c.ogUser.LastName,
			c.ogUser.PhotoURL,
		)

		query := "delete from users where id=?"

		if c.expectError {
			mock.ExpectPrepare("delete").ExpectExec().WithArgs(c.submittedID)
			db.Prepare(query)
			// Test Delete()
			err := mainSQLStore.Delete(c.submittedID)
			if err == nil {
				t.Errorf("Expected error but did not get one.")
			}
		} else {
			mock.ExpectPrepare("delete").ExpectExec().WithArgs(c.submittedID).WillReturnResult(sqlmock.NewResult(1, 1))
			db.Prepare(query)
			// Test Delete()
			err := mainSQLStore.Delete(c.submittedID)
			if err != nil {
				t.Errorf("Unexpected error on successful test [%s]: %v", c.name, err)
			}
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unfulfilled expectations: %s", err)
		}

	}
}
