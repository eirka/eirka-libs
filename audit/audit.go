package audit

import (
	"errors"

	"github.com/eirka/eirka-libs/db"
)

// LogType type
type LogType uint

const (
	// BoardLog is for board specific events
	BoardLog LogType = iota + 1
	// ModLog is for mod events
	ModLog
	// UserLog is for user events
	UserLog
)

var (
	// Board events

	// AuditReply is for new reply events
	AuditReply = "Replied"
	// AuditNewTag is for tag creation events
	AuditNewTag = "Tag Created"
	// AuditAddTag is for adding a tag to an image events
	AuditAddTag = "Tag Added"
	// AuditNewThread is for thread creation events
	AuditNewThread = "Thread Created"

	// Mod events

	// AuditCloseThread is for thread closure events
	AuditCloseThread = "Thread Closed"
	// AuditOpenThread is for thread opening events
	AuditOpenThread = "Thread Opened"
	// AuditStickyThread is for thread sticky events
	AuditStickyThread = "Thread Stickied"
	// AuditUnstickyThread is for thread unsticky events
	AuditUnstickyThread = "Thread Unstickied"
	// AuditDeleteThread is for thread deletion events
	AuditDeleteThread = "Thread Deleted"
	// AuditPurgeThread is for thread purging events
	AuditPurgeThread = "Thread Purged"
	// AuditDeletePost is for post deletion events
	AuditDeletePost = "Post Deleted"
	// AuditPurgePost is for post purging events
	AuditPurgePost = "Post Purged"
	// AuditPurge is for total purging events
	AuditPurge = "Deleted Items Purged"
	// AuditFlushCache is for cache flushing events
	AuditFlushCache = "Cache Flushed"
	// AuditBanIP is for user banning events
	AuditBanIP = "IP Banned"
	// AuditBanFile is for file banning events
	AuditBanFile = "File Banned"
	// AuditSpam is for spam reporting events
	AuditSpam = "Spam Reported"
	// AuditUpdateTag is for admin tag update events
	AuditUpdateTag = "Tag Updated"
	// AuditDeleteTag is for admin tag delete events
	AuditDeleteTag = "Tag Deleted"
	// AuditDeleteImageTag is for admin image tag delete events
	AuditDeleteImageTag = "Image Tag Deleted"

	// User events

	// AuditRegister is for user registration events
	AuditRegister = "Account Registered"
	// AuditChangePassword is for user password update events
	AuditChangePassword = "Password Changed"
	// AuditResetPassword is for user password reset events
	AuditResetPassword = "Password Reset"
	// AuditEmailUpdate is for user email update events
	AuditEmailUpdate = "Email Updated"
	// AuditFavoriteRemoved is for user favorite removal events
	AuditFavoriteRemoved = "Favorite Removed"
	// AuditFavoriteAdded is for user favorite events
	AuditFavoriteAdded = "Favorite Added"
)

// Audit adds an action to the audit log
type Audit struct {
	User   uint
	Ib     uint
	Type   LogType
	IP     string
	Action string
	Info   string
}

// IsValid will check struct validity
func (m *Audit) IsValid() bool {

	if m.User == 0 {
		return false
	}

	if m.Ib == 0 {
		return false
	}

	if m.Type == 0 {
		return false
	}

	if m.IP == "" {
		return false
	}

	if m.Action == "" {
		return false
	}

	if m.Info == "" {
		return false
	}

	return true

}

// Submit will insert audit info into the audit log
func (m *Audit) Submit() (err error) {

	if !m.IsValid() {
		return errors.New("Audit not valid")
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	_, err = dbase.Exec("INSERT INTO audit (user_id,ib_id,audit_type,audit_ip,audit_time,audit_action,audit_info) VALUES (?,?,?,?,NOW(),?,?)",
		m.User, m.Ib, m.Type, m.IP, m.Action, m.Info)
	if err != nil {
		return
	}

	return
}
