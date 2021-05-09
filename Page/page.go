package Page

// page页面
type Page struct {
	// 目标页面Url
	Url string

	// 目标内容ID
	ElementId string

	// true 中文页面 false 英文页面
	PageType bool
}

func NewPage(Url, Id string, Type bool) *Page {
	page := new(Page)
	page.Url = Url
	page.ElementId = Id
	page.PageType = Type
	return page
}
