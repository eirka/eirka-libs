package config

import "github.com/eirka/eirka-libs/db"

// GetDatabaseSettings gets limits that are in the database
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

	// akismet has been configured
	if Settings.Akismet.Key != "" {
		Settings.Akismet.Configured = true
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

	err = ps.QueryRow("amazon_id").Scan(&Settings.Amazon.ID)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("amazon_key").Scan(&Settings.Amazon.Key)
	if err != nil {
		panic(err)
	}

	// amazon has been configured
	if Settings.Amazon.ID != "" && Settings.Amazon.Key != "" {
		Settings.Amazon.Configured = true
	}

	err = ps.QueryRow("thread_postsperpage").Scan(&Settings.Limits.PostsPerPage)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("index_threadsperpage").Scan(&Settings.Limits.ThreadsPerPage)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("index_postsperthread").Scan(&Settings.Limits.PostsPerThread)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("prim_js").Scan(&Settings.Prim.JS)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("prim_css").Scan(&Settings.Prim.CSS)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("avatar_minwidth").Scan(&Settings.Limits.AvatarMinWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("avatar_minheight").Scan(&Settings.Limits.AvatarMinHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("avatar_maxwidth").Scan(&Settings.Limits.AvatarMaxWidth)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("avatar_maxheight").Scan(&Settings.Limits.AvatarMaxHeight)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("avatar_maxsize").Scan(&Settings.Limits.AvatarMaxSize)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("password_maxlength").Scan(&Settings.Limits.PasswordMaxLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("password_minlength").Scan(&Settings.Limits.PasswordMinLength)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("cloudflare_email").Scan(&Settings.CloudFlare.Email)
	if err != nil {
		panic(err)
	}

	err = ps.QueryRow("cloudflare_key").Scan(&Settings.CloudFlare.Key)
	if err != nil {
		panic(err)
	}

	// cloudflare has been configured
	if Settings.CloudFlare.Key != "" {
		Settings.CloudFlare.Configured = true
	}

}
