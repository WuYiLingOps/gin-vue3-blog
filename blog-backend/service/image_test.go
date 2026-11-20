/*
 * 项目名称：blog-backend
 * 文件名称：image_test.go
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片管理单元测试，覆盖 settings 值的图片 URL 提取与外链误判修复
 */
package service

import "testing"

func TestExtractJSONArrayStrings(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  int
	}{
		{"正常数组", `["https://a.com/1.png","https://a.com/2.jpg"]`, 2},
		{"非数组", "https://gitee.com/WuYiLingOps", 0},
		{"空数组", `[]`, 0},
		{"非法 JSON", `[broken`, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := extractJSONArrayStrings(c.value)
			if len(got) != c.want {
				t.Errorf("extractJSONArrayStrings(%q) = %v, want %d elements", c.value, got, c.want)
			}
		})
	}
}

func TestLooksLikeImageURL(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"https://a.com/img/1.png", true},
		{"https://a.com/img/1.PNG", true},
		{"https://a.com/1.webp?x-oss-process=style/resize", true},
		{"https://gitee.com/WuYiLingOps", false},
		{"http://localhost:3000", false},
		{"https://a.com/page", false},
		{"https://a.com/", false},
		{"", false},
	}
	for _, c := range cases {
		if got := looksLikeImageURL(c.url); got != c.want {
			t.Errorf("looksLikeImageURL(%q) = %v, want %v", c.url, got, c.want)
		}
	}
}

func TestExtractJSONImageFields(t *testing.T) {
	// friendlink_info 场景：avatar/screenshot 字段应提取，url/rss/desc/name 不应提取
	friendlink := `{"name":"我的博客","desc":"描述","url":"https://blog.example.com","avatar":"https://pic.imgdb.cn/item/64.png","screenshot":"https://pic.imgdb.cn/item/65.jpg","rss":"https://blog.example.com/rss"}`
	got := extractJSONImageFields(friendlink)
	if len(got) != 2 {
		t.Errorf("friendlink_info 提取数 = %d (%v), want 2", len(got), got)
	}

	// 社交链接配置（截图误判场景）：social_gitee 不是图片字段
	social := `{"social_gitee":"https://gitee.com/WuYiLingOps","social_github":"https://github.com/x"}`
	if got := extractJSONImageFields(social); len(got) != 0 {
		t.Errorf("社交配置不应提取图片, got %v", got)
	}

	// about_personality 场景：img 字段应提取
	personality := `{"type":"INTJ","img":"https://cdn.example.com/intj.png"}`
	if got := extractJSONImageFields(personality); len(got) != 1 {
		t.Errorf("about img 提取数 = %d, want 1", len(got))
	}

	// 非 URL 值（图标名）：不应提取
	badValue := `{"icon":"compass-outline","img":"compass-outline"}`
	if got := extractJSONImageFields(badValue); len(got) != 0 {
		t.Errorf("非 URL 值不应提取, got %v", got)
	}

	// 数组值被忽略（由 extractJSONArrayStrings 单独处理）
	withArray := `{"images":["https://a.com/1.png"],"cover":"https://a.com/2.png"}`
	if got := extractJSONImageFields(withArray); len(got) != 1 {
		t.Errorf("对象内嵌数组应忽略, got %v", got)
	}
}

func TestExtractPlainURL(t *testing.T) {
	// 纯图片 URL 应提取
	if got := extractPlainURL("https://a.com/bg.png"); got == "" {
		t.Error("纯图片 URL 应被提取")
	}
	// 非图片 URL 仍会被 extractPlainURL 返回，但外层会做 looksLikeImageURL 过滤（截图误判场景的最终防线在外层）
	if got := extractPlainURL("http://localhost:3000"); got == "" {
		t.Error("纯 URL 应被 extractPlainURL 返回，过滤由外层 looksLikeImageURL 负责")
	}
	// 带空格的不是纯 URL
	if got := extractPlainURL("hello world https://a.com"); got != "" {
		t.Error("含空格不应提取")
	}
	// 本地路径
	if got := extractPlainURL("/uploads/siteCovers/1.png"); got == "" {
		t.Error("/uploads/ 路径应被提取")
	}
}
