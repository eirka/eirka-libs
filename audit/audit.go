package audit

import (
	"errors"
	"github.com/eirka/eirka-libs/db"
)

type LogType uint

const (
	// what kind of event the audit is
	BoardLog LogType = iota + 1 // 1
	ModLog                      // 2
	UserLog                     // 3
)

var (
	AuditReply           = "Replied"
	AuditNewTag          = "Tag Created"
	AuditAddTag          = "Tag Added"
	AuditUpdateTag       = "Tag Updated"
	AuditDeleteTag       = "Tag Deleted"
	AuditDeleteImageTag  = "Image Tag Deleted"
	AuditNewThread       = "Thread Created"
	AuditCloseThread     = "Thread Closed"
	AuditOpenThread      = "Thread Opened"
	AuditStickyThread    = "Thread Stickied"
	AuditUnstickyThread  = "Thread Unstickied"
	AuditDeleteThread    = "Thread Deleted"
	AuditPurgeThread     = "Thread Purged"
	AuditDeletePost      = "Post Deleted"
	AuditPurgePost       = "Post Purged"
	AuditPurge           = "Deleted Items Purged"
	AuditFlushCache      = "Cache Flushed"
	AuditBanIp           = "IP Banned"
	AuditBanFile         = "File Banned"
	AuditSpam            = "Spam Reported"
	AuditRegister        = "Account Registered"
	AuditChangePassword  = "Password Changed"
	AuditEmailUpdate     = "Email Updated"
	AuditFavoriteRemoved = "Favorite Removed"
	AuditFavoriteAdded   = "Favorite Added"
)

// Audit adds an action to the audit log
type Audit struct {
	User   uint
	Ib     uint
	Type   LogType
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

	if a.Type == 0 {
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

	_, err = dbase.Exec("INSERT INTO audit (user_id,ib_id,audit_type,audit_ip,audit_time,audit_action,audit_info) VALUES (?,?,?,?,NOW(),?,?)",
		a.User, a.Ib, a.Type, a.Ip, a.Action, a.Info)
	if err != nil {
		return
	}

	return
}
