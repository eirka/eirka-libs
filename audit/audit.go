package audit

import (
	"errors"
	"github.com/eirka/eirka-libs/db"
)

// Audit adds an action to the audit log
type Audit struct {
	User   uint
	Ib     uint
	Ip     string
	Action string
	Info   string
}

// check struct validity
func (a *Audit) IsValid() bool {

	if a.User == 0 {
		return false
	}

	if a.Ib == 0 {
		return false
	}

	if a.Ip == "" {
		return false
	}

	if a.Action == "" {
		return false
	}

	if a.Info == "" {
		return false
	}

	return true

}

// Submit will insert audit info into the audit log
func (a *Audit) Submit() (err error) {

	if !a.IsValid() {
		return errors.New("Audit not valid")
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	_, err = dbase.Exec("INSERT INTO audit (user_id,ib_id,audit_ip,audit_time,audit_action,audit_info) VALUES (?,?,?,NOW(),?,?)",
		a.User, a.Ib, a.Ip, a.Action, a.Info)
	if err != nil {
		return
	}

	return
}
