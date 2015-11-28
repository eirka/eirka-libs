package perms

import (
	"database/sql"

	"github.com/techjanitor/pram-libs/db"
	e "github.com/techjanitor/pram-libs/errors"
)

// get the user info from id
func Check(uid, ib uint) (allowed bool, err error) {

	// check for invalid stuff
	if uid == 0 || uid == 1 || ib == 0 {
		err = e.ErrInvalidParam
		return
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	// holds our role
	var role uint

	// get data from users table
	err = dbase.QueryRow(`SELECT COALESCE((SELECT MAX(role_id) FROM user_ib_role_map WHERE user_ib_role_map.user_id = users.user_id AND ib_id = ?),user_role_map.role_id) as role
    FROM users
    INNER JOIN user_role_map ON (user_role_map.user_id = users.user_id)
    WHERE users.user_id = ?`, u.Id).Scan(&role)
	if err != nil {
		return
	}

	switch role {
	case 3:
		allowed = true
	case 4:
		allowed = true
	default:
		allowed = false
	}

	return

}
