package db_test

import (
	"bytes"
	"database/sql"
	"strings"
	"testing"
	"time"

	"reflect"

	"github.com/tomogoma/authms/db"
	"github.com/tomogoma/authms/model"
)

func TestRoach_InsertUserEmailAtomic_nilTx(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	_, err := r.InsertUserEmailAtomic(nil, usr.ID, "test@mailinator.com", false)
	if err == nil {
		t.Errorf("(nil tx) - expected an error, got nil")
	}
}

// TestRoach_InsertUserEmailAtomic shares test cases with TestRoach_InsertUserEmail
// because they use the same underlying implementation.
func TestRoach_InsertUserEmailAtomic(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	r.ExecuteTx(func(tx *sql.Tx) error {
		ret, err := r.InsertUserEmailAtomic(tx, usr.ID, "test@mailinator.com", false)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		if ret == nil {
			t.Fatalf("Got nil group")
		}
		if ret.ID == "" {
			t.Errorf("ID was not assigned")
		}
		if ret.UpdateDate.Before(time.Now().Add(-1 * time.Minute)) {
			t.Errorf("UpdateDate was not assigned")
		}
		if ret.CreateDate.Before(time.Now().Add(-1 * time.Minute)) {
			t.Errorf("CreateDate was not assigned")
		}
		if ret.UserID != usr.ID {
			t.Errorf("User ID mismatch, expect %s, got %s",
				usr.ID, ret.UserID)
		}
		if ret.Address != "test@mailinator.com" {
			t.Errorf("Address mismatch, expect test@mailinator.com, got %s",
				ret.Address)
		}
		if ret.Verified != false {
			t.Errorf("verified mismatch, expect false, got %t", ret.Verified)
		}
		return nil
	})
}

func TestRoach_InsertUserEmail(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	tt := []struct {
		testName string
		usrID    string
		addr     string
		verified bool
		expErr   bool
	}{
		{testName: "valid", usrID: usr.ID, addr: "test@mailinator.com", verified: true, expErr: false},
		{testName: "bad user ID", usrID: "bad id", addr: "test@mailinator.com", verified: true, expErr: true},
		{testName: "empty phone", usrID: usr.ID, addr: "", expErr: true},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			ret, err := r.InsertUserEmail(tc.usrID, tc.addr, tc.verified)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if ret == nil {
				t.Fatalf("Got nil group")
			}
			if ret.ID == "" {
				t.Errorf("ID was not assigned")
			}
			if ret.UpdateDate.Before(time.Now().Add(-1 * time.Minute)) {
				t.Errorf("UpdateDate was not assigned")
			}
			if ret.CreateDate.Before(time.Now().Add(-1 * time.Minute)) {
				t.Errorf("CreateDate was not assigned")
			}
			if ret.UserID != tc.usrID {
				t.Errorf("User ID mismatch, expect %s, got %s",
					tc.usrID, ret.UserID)
			}
			if ret.Address != tc.addr {
				t.Errorf("address mismatch, expect %s, got %s",
					tc.addr, ret.Address)
			}
			if ret.Verified != tc.verified {
				t.Errorf("verified mismatch, expect %t, got %t",
					tc.verified, ret.Verified)
			}
			return
		})
	}
}

func TestRoach_UpdateUserEmail(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	expMail := insertEmail(t, r, usr.ID)
	expMail.Address = "new@addr.mail"
	expMail.Verified = true
	tt := []struct {
		name         string
		userID       string
		newAddr      string
		newVerStatus bool
		expNotFound  bool
	}{
		{name: "valid", userID: usr.ID, newAddr: expMail.Address, newVerStatus: expMail.Verified, expNotFound: false},
		{name: "not found", userID: "123", newAddr: expMail.Address, newVerStatus: expMail.Verified, expNotFound: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			updMail, err := r.UpdateUserEmail(tc.userID, tc.newAddr, tc.newVerStatus)
			if tc.expNotFound {
				if !r.IsNotFoundError(err) {
					t.Errorf("Expected IsNotFound, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if updMail.UpdateDate.Equal(expMail.CreateDate) || updMail.UpdateDate.Before(expMail.CreateDate) {
				t.Errorf("Update date not set correctly before/equal to create date")
			}
			expMail.UpdateDate = updMail.UpdateDate
			if !reflect.DeepEqual(expMail, updMail) {
				t.Errorf("Email mismatch:\nExpect:\t%+v\nGot:\t%+v",
					expMail, updMail)
			}
		})
	}
}

func TestRoach_UpdateUserEmailAtomic_nilTx(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	mail := insertEmail(t, r, usr.ID)
	_, err := r.UpdateUserEmailAtomic(nil, usr.ID, mail.ID, true)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

func TestRoach_InsertEmailTokenAtomic_nilTx(t *testing.T) {
	setupTime := time.Now()
	dbt := []byte(strings.Repeat("x", 57))
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	email := insertEmail(t, r, usr.ID)
	_, err := r.InsertEmailTokenAtomic(nil, usr.ID, email.Address, dbt, false, setupTime)
	if err == nil {
		t.Errorf("(nil tx) - expected an error, got nil")
	}
}

// TestRoach_InsertEmailTokenAtomic shares test cases with TestRoach_InsertEmailToken
// because they use the same underlying implementation.
func TestRoach_InsertEmailTokenAtomic(t *testing.T) {
	setupTime := time.Now()
	dbt := []byte(strings.Repeat("x", 57))
	isUsed := false
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	email := insertEmail(t, r, usr.ID)
	r.ExecuteTx(func(tx *sql.Tx) error {
		ret, err := r.InsertEmailTokenAtomic(tx, usr.ID, email.Address, dbt, isUsed, setupTime)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		if ret == nil {
			t.Fatalf("Got nil group")
		}
		if ret.ID == "" {
			t.Errorf("ID was not assigned")
		}
		if ret.IssueDate.Before(setupTime) {
			t.Errorf("Issue date was not assigned")
		}
		if ret.UserID != usr.ID {
			t.Errorf("User ID mismatch, expect %s, got %s",
				usr.ID, ret.UserID)
		}
		if ret.Address != email.Address {
			t.Errorf("Invalid email: expect %s, got %s", email.Address, ret.Address)
		}
		if !bytes.Equal(ret.Token, dbt) {
			t.Errorf("Invalid db token: expect %s, got %s", dbt, ret.Token)
		}
		if ret.IsUsed != isUsed {
			t.Errorf("Invalid used val: expect %t, got %t", isUsed, ret.IsUsed)
		}
		// TODO truncate tc.expiry appropriately before testing
		//if ret.ExpiryDate != setupTime {
		//	t.Errorf("Invalid expiry: expect %v, got %v", setupTime, ret.ExpiryDate)
		//}
		return nil
	})
}

func TestRoach_InsertEmailToken(t *testing.T) {
	setUpTime := time.Now()
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	email := insertEmail(t, r, usr.ID)
	validDBT := []byte(strings.Repeat("x", 57))
	tt := []struct {
		testName string
		usrID    string
		addr     string
		dbt      []byte
		isUsed   bool
		expiry   time.Time
		expErr   bool
	}{
		{
			testName: "valid",
			usrID:    usr.ID,
			addr:     email.Address,
			dbt:      validDBT,
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   false,
		},
		{
			testName: "bad user id",
			usrID:    "bad id",
			addr:     email.Address,
			dbt:      validDBT,
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   true,
		},
		{
			testName: "empty emaili",
			usrID:    usr.ID,
			addr:     "",
			dbt:      validDBT,
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   true,
		},
		{
			testName: "email not exists",
			usrID:    usr.ID,
			addr:     "not exists email",
			dbt:      validDBT,
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   true,
		},
		{
			testName: "empty dbt",
			usrID:    usr.ID,
			addr:     email.Address,
			dbt:      []byte{},
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   true,
		},
		{
			testName: "nil dbt",
			usrID:    usr.ID,
			addr:     email.Address,
			dbt:      nil,
			isUsed:   false,
			expiry:   setUpTime.Add(5 * time.Minute),
			expErr:   true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			ret, err := r.InsertEmailToken(tc.usrID, tc.addr, tc.dbt, tc.isUsed, tc.expiry)
			if tc.expErr {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if ret == nil {
				t.Fatalf("Got nil group")
			}
			if ret.ID == "" {
				t.Errorf("ID was not assigned")
			}
			if ret.IssueDate.Before(setUpTime) {
				t.Errorf("Issue date was not assigned")
			}
			if ret.UserID != tc.usrID {
				t.Errorf("User ID mismatch, expect %s, got %s",
					tc.usrID, ret.UserID)
			}
			if ret.Address != tc.addr {
				t.Errorf("Invalid email: expect %s, got %s", tc.addr, ret.Address)
			}
			if !bytes.Equal(ret.Token, tc.dbt) {
				t.Errorf("Invalid db token: expect %s, got %s", tc.dbt, ret.Token)
			}
			if ret.IsUsed != tc.isUsed {
				t.Errorf("Invalid used val: expect %t, got %t", tc.isUsed, ret.IsUsed)
			}
			// TODO truncate tc.expiry appropriately before testing
			//if ret.ExpiryDate != tc.expiry {
			//	t.Errorf("Invalid expiry: expect %v, got %v", tc.expiry, ret.ExpiryDate)
			//}
			return
		})
	}
}

func TestRoach_EmailTokens(t *testing.T) {
	conf := setup(t)
	defer tearDown(t, conf)
	r := newRoach(t, conf)
	usr := insertUser(t, r)
	usrNoTkns := insertUser(t, r)
	email := insertEmail(t, r, usr.ID)
	dbt1 := insertEmailToken(t, r, usr.ID, email.Address)
	dbt2 := insertEmailToken(t, r, usr.ID, email.Address)
	expDBTs := []model.DBToken{*dbt2, *dbt1}
	tt := []struct {
		name        string
		userID      string
		expNotFound bool
	}{
		{name: "found", userID: usr.ID, expNotFound: false},
		{name: "not found", userID: usrNoTkns.ID, expNotFound: true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actDBTs, err := r.EmailTokens(tc.userID, 0, 2)
			if tc.expNotFound {
				if !r.IsNotFoundError(err) {
					t.Fatalf("Expected not found error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got error: %v", err)
			}
			if !reflect.DeepEqual(expDBTs, actDBTs) {
				t.Errorf("DBTokens mismatch:\nExpect:\t%+v\nGot:\t%+v",
					expDBTs, actDBTs)
			}
		})
	}
}

func insertEmail(t *testing.T, r *db.Roach, usrID string) *model.VerifLogin {
	m, err := r.InsertUserEmail(usrID, "test@mailinator.com", false)
	if err != nil {
		t.Fatalf("Error setting up: insert email: %v", err)
	}
	return m
}

func insertEmailToken(t *testing.T, r *db.Roach, usrID, email string) *model.DBToken {
	tkn, err := r.InsertEmailToken(
		usrID,
		email,
		[]byte(strings.Repeat("x", 56)),
		false,
		time.Now().Add(5*time.Minute),
	)
	if err != nil {
		t.Fatalf("Error setting up: insert email token: %v", err)
	}
	return tkn
}
