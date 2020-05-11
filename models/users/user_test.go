package users

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

//TODO: add tests for the various functions in user.go, as described in the assignment.
//use `go test -cover` to ensure that you are covering all or nearly all of your code paths.
func TestValidate(t *testing.T) {
	cases := []struct {
		name      string
		testUser  *NewUser
		expectErr bool
	}{
		{
			"Invalid Email Address",
			&NewUser{"blindBanditskajdbsjb.edu", "TwinkleToes1337", "TwinkleToes1337", "TheBlindBandit", "Toph", "Beifong"},
			true,
		},
		{
			"Short Password",
			&NewUser{"jelauria@uw.edu", "Rocks", "Rocks", "TheBlindBandit", "Joyce", "Elauria"},
			true,
		},
		{
			"Non-matching Passwords",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks133", "TheBlindBandit", "Joyce", "Elauria"},
			true,
		},
		{
			"Blank Password",
			&NewUser{"jelauria@uw.edu", "", "Rocks", "TheBlindBandit", "Joyce", "Elauria"},
			true,
		},
		{
			"Blank Username",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks1337", "", "Joyce", "Elauria"},
			true,
		},
		{
			"Username with Spaces",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks1337", "The Blind Bandit", "Joyce", "Elauria"},
			true,
		},
		{
			"Valid User",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
	}

	for _, c := range cases {
		err := c.testUser.Validate()
		if c.expectErr && err == nil {
			t.Errorf("case %s: expected error but did not get one", c.name)
		}
		if !c.expectErr && err != nil {
			t.Errorf("case %s: unexpected error '%v'", c.name, err)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct {
		name     string
		testUser *User
		expected string
	}{
		{
			"First Name Only",
			&User{0, "jelauria@uw.edu", nil, "TheBlindBandit", "Joyce", "", "test.com"},
			"Joyce",
		},
		{
			"Last Name Only",
			&User{0, "jelauria@uw.edu", nil, "TheBlindBandit", "", "Elauria", "test.com"},
			"Elauria",
		},
		{
			"No Names",
			&User{0, "jelauria@uw.edu", nil, "TheBlindBandit", "", "", "test.com"},
			"",
		},
		{
			"Both Names Filled",
			&User{0, "jelauria@uw.edu", nil, "TheBlindBandit", "Joyce", "Elauria", "test.com"},
			"Joyce Elauria",
		},
	}
	for _, c := range cases {
		result := c.testUser.FullName()
		if result != c.expected {
			t.Errorf("case %s: expected '%s' but instead got '%s'", c.name, c.expected, result)
		}
	}
}

func TestToUser(t *testing.T) {
	cases := []struct {
		name      string
		testUser  *NewUser
		expectErr bool
	}{
		{
			"Valid User",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
		{
			"Invalid User",
			&NewUser{"jelauria@uw.edu", "TophRocks1337", "TophRocks137", "The Blind Bandit", "Joyce", "Elauria"},
			true,
		},
		{
			"Email with Uppercase Characters",
			&NewUser{"JElauria@uw.edu", "TophRocks1337", "TophRocks1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
		{
			"Email with Spaces",
			&NewUser{"jelauria@uw.edu ", "TophRocks1337", "TophRocks1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
		{
			"Email with Uppercase and Spaces",
			&NewUser{" JElauria@uw.edu ", "TophRocks1337", "TophRocks1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
		{
			"Password with Symbols",
			&NewUser{"jelauria@uw.edu", "Toph!Rocks!#1337", "Toph!Rocks!#1337", "TheBlindBandit", "Joyce", "Elauria"},
			false,
		},
	}
	expectImgUrl := "https://www.gravatar.com/avatar/c95dc22c4910b86b78391d834f120947"
	for _, c := range cases {
		result, err := c.testUser.ToUser()
		if c.expectErr && err == nil {
			t.Errorf("case %s: expected error but did not get one", c.name)
		}
		if !c.expectErr && err != nil {
			t.Errorf("case %s: unexpected error '%v'", c.name, err)
		}
		if err == nil && !c.expectErr {
			if result.PhotoURL != expectImgUrl {
				t.Errorf("case %s: incorrect photo url\nEXPECTED:%s\nACTUAL:%s", c.name, expectImgUrl, result.PhotoURL)
			}
			passErr := bcrypt.CompareHashAndPassword(result.PassHash, []byte(c.testUser.Password))
			if passErr != nil {
				t.Errorf("case %s: password was hashed improperly, received error: %v", c.name, passErr)
			}
		}
	}
}

func TestAuthenticate(t *testing.T) {
	nUser := &NewUser{"jelauria@uw.edu", "Toph!Rocks!#1337", "Toph!Rocks!#1337", "TheBlindBandit", "Joyce", "Elauria"}
	testUser, _ := nUser.ToUser()
	cases := []struct {
		name          string
		submittedPass string
		expectErr     bool
	}{
		{
			"Correct Password",
			"Toph!Rocks!#1337",
			false,
		},
		{
			"Incorrect Password",
			"Toh!Rocks!#1337",
			true,
		},
		{
			"No Password Submitted",
			"",
			true,
		},
	}
	for _, c := range cases {
		err := testUser.Authenticate(c.submittedPass)
		if c.expectErr && err == nil {
			t.Errorf("case %s: expected error but did not get one", c.name)
		}
		if !c.expectErr && err != nil {
			t.Errorf("case %s: unexpected error '%v'", c.name, err)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name          string
		updates       *Updates
		expectedFname string
		expectedLname string
	}{
		{
			"Update First Name",
			&Updates{"Toph", "Elauria"},
			"Toph",
			"Elauria",
		},
		{
			"Update Last Name",
			&Updates{"Joyce", "Beifong"},
			"Joyce",
			"Beifong",
		},
		{
			"Update Both Names",
			&Updates{"Toph", "Beifong"},
			"Toph",
			"Beifong",
		},
		{
			"Reduce to Last Name",
			&Updates{"", "Elauria"},
			"",
			"Elauria",
		},
		{
			"Reduce to First Name",
			&Updates{"Joyce", ""},
			"Joyce",
			"",
		},
		{
			"Remove Both Names",
			&Updates{"", ""},
			"",
			"",
		},
	}
	nUser := &NewUser{"jelauria@uw.edu", "Toph!Rocks!#1337", "Toph!Rocks!#1337", "TheBlindBandit", "Joyce", "Elauria"}
	for _, c := range cases {
		testUser, _ := nUser.ToUser()
		err := testUser.ApplyUpdates(c.updates)
		if err != nil {
			t.Errorf("case %s: unexpected error '%v'", c.name, err)
		}
		if testUser.FirstName != c.expectedFname {
			t.Errorf("case %s: error with first name following update\nEXPECTED:%s\nACTUAL:%s", c.name, c.expectedFname, testUser.FirstName)
		}
		if testUser.LastName != c.expectedLname {
			t.Errorf("case %s: error with last name following update\nEXPECTED:%s\nACTUAL:%s", c.name, c.expectedLname, testUser.LastName)
		}
	}
}
