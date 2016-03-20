package audit

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
)

func TestAudit(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	mock.ExpectExec(`INSERT INTO audit \(user_id,ib_id,audit_type,audit_ip,audit_time,audit_action,audit_info\)`).
		WithArgs(1, 1, UserLog, "10.0.0.1", AuditEmailUpdate, "meta info").
		WillReturnResult(sqlmock.NewResult(1, 1))

	audit := Audit{
		User:   1,
		Ib:     1,
		Type:   UserLog,
		Ip:     "10.0.0.1",
		Action: AuditEmailUpdate,
		Info:   "meta info",
	}

	// submit audit
	err = audit.Submit()
	assert.NoError(t, err, "An error was not expected")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestAuditInvalid(t *testing.T) {

	var err error

	audit := Audit{
		User:   0,
		Ib:     1,
		Type:   UserLog,
		Ip:     "10.0.0.1",
		Action: AuditEmailUpdate,
		Info:   "meta info",
	}

	// submit audit
	err = audit.Submit()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, errors.New("Audit not valid"), "Error should match")
	}

}
