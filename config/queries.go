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

	err = ps.QueryRow("antispam_key").Scan(&config.Settings.Antispam.AntispamKey)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("antispam_cookiename").Scan(&config.Settings.Antispam.CookieName)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("antispam_cookievalue").Scan(&config.Settings.Antispam.CookieValue)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_minwidth").Scan(&config.Settings.Limits.ImageMinWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_minheight").Scan(&config.Settings.Limits.ImageMinHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxwidth").Scan(&config.Settings.Limits.ImageMaxWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxheight").Scan(&config.Settings.Limits.ImageMaxHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("image_maxsize").Scan(&config.Settings.Limits.ImageMaxSize)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("webm_maxlength").Scan(&config.Settings.Limits.WebmMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thread_postsmax").Scan(&config.Settings.Limits.PostsMax)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("comment_maxlength").Scan(&config.Settings.Limits.CommentMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("comment_minlength").Scan(&config.Settings.Limits.CommentMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("title_maxlength").Scan(&config.Settings.Limits.TitleMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("title_minlength").Scan(&config.Settings.Limits.TitleMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("name_maxlength").Scan(&config.Settings.Limits.NameMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("name_minlength").Scan(&config.Settings.Limits.NameMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("tag_maxlength").Scan(&config.Settings.Limits.TagMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("tag_minlength").Scan(&config.Settings.Limits.TagMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thumbnail_maxwidth").Scan(&config.Settings.Limits.ThumbnailMaxWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("thumbnail_maxheight").Scan(&config.Settings.Limits.ThumbnailMaxHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("param_maxsize").Scan(&config.Settings.Limits.ParamMaxSize)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("guest_posting").Scan(&config.Settings.General.GuestPosting)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("auto_registration").Scan(&config.Settings.General.AutoRegistration)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("akismet_key").Scan(&config.Settings.Akismet.Key)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("akismet_host").Scan(&config.Settings.Akismet.Host)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("sfs_confidence").Scan(&config.Settings.StopForumSpam.Confidence)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_region").Scan(&config.Settings.Amazon.Region)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_bucket").Scan(&config.Settings.Amazon.Bucket)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_id").Scan(&config.Settings.Amazon.Id)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_key").Scan(&config.Settings.Amazon.Key)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_auth").Scan(&config.Settings.Google.Auth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_bucket").Scan(&config.Settings.Google.Bucket)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("google_key").Scan(&config.Settings.Google.Key)
	if err != nil {
		panic(err)
	}

	return

}
