package sqlbuild

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
)

const (
	ItemPerPage  = 10
	DefaultLimit = "10"
	baseQuery    = `SELECT news.id, title, content, pubtime, link
	FROM news`
	countQuery = `SELECT count(*) FROM news`
)

func NewsQuery(url *url.URL, count bool, offset int) string {
	//n лимит записей в выдаче
	limit := url.Query().Get("n")
	//keyword поиск по ключевому слову
	keyword := url.Query().Get("s")
	//page страница
	//page := url.Query().Get("page")

	buf := bytes.NewBuffer([]byte{})

	if !count {
		buf.WriteString(baseQuery)
	} else {
		buf.WriteString(countQuery)
	}

	if keyword != "" {
		pattern := strings.Join([]string{"%", keyword, "%"}, "")
		buf.WriteString(" WHERE title ILIKE ")
		buf.WriteString("'")
		buf.WriteString(pattern)
		buf.WriteString("'")
	}

	if !count {
		buf.WriteString(" ORDER BY id ASC")
	}

	if limit != "" {
		buf.WriteString(" LIMIT ")
		buf.WriteString(limit)
	} else {
		buf.WriteString(" LIMIT ")
		buf.WriteString(DefaultLimit)
	}

	if offset > 0 {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(offset))
	}

	return buf.String()

}
