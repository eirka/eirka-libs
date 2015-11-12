package config

import "github.com/techjanitor/pram-libs/db"

// Get limits that are in the database
func GetDatabaseSettings() {

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		panic(err)
	}

	ps, err := dbase.Prepare("SELECT settings_value FROM settings WHERE settings_key = ? LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer ps.Close()

	err = ps.QueryRow("antispam_key").Scan(&Settings.Antispam.AntispamKey)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("antispam_cookiename").Scan(&Settings.Antispam.CookieName)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("antispam_cookievalue").Scan(&Settings.Antispam.CookieValue)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_minwidth").Scan(&Settings.Limits.ImageMinWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_minheight").Scan(&Settings.Limits.ImageMinHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxwidth").Scan(&Settings.Limits.ImageMaxWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxheight").Scan(&Settings.Limits.ImageMaxHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxsize").Scan(&Settings.Limits.ImageMaxSize)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("webm_maxlength").Scan(&Settings.Limits.WebmMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thread_postsmax").Scan(&Settings.Limits.PostsMax)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("comment_maxlength").Scan(&Settings.Limits.CommentMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("comment_minlength").Scan(&Settings.Limits.CommentMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("title_maxlength").Scan(&Settings.Limits.TitleMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("title_minlength").Scan(&Settings.Limits.TitleMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("name_maxlength").Scan(&Settings.Limits.NameMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("name_minlength").Scan(&Settings.Limits.NameMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("tag_maxlength").Scan(&Settings.Limits.TagMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("tag_minlength").Scan(&Settings.Limits.TagMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thumbnail_maxwidth").Scan(&Settings.Limits.ThumbnailMaxWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thumbnail_maxheight").Scan(&Settings.Limits.ThumbnailMaxHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("param_maxsize").Scan(&Settings.Limits.ParamMaxSize)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("guest_posting").Scan(&Settings.General.GuestPosting)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("auto_registration").Scan(&Settings.General.AutoRegistration)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("akismet_key").Scan(&Settings.Akismet.Key)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("akismet_host").Scan(&Settings.Akismet.Host)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("sfs_confidence").Scan(&Settings.StopForumSpam.Confidence)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_region").Scan(&Settings.Amazon.Region)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_bucket").Scan(&Settings.Amazon.Bucket)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_id").Scan(&Settings.Amazon.Id)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_key").Scan(&Settings.Amazon.Key)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_auth").Scan(&Settings.Google.Auth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_bucket").Scan(&Settings.Google.Bucket)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_key").Scan(&Settings.Google.Key)
	if err != nil {
		panic(err)
	}

	return

}
