/*
 * 项目名称：blog-backend
 * 文件名称：role.go
 * 创建时间：2026-02-06 21:57:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户角色相关常量与工具方法（独立包），用于避免 util<->repository 等循环依赖。
 */
package constant

// 用户角色常量定义
const (
	// RoleSuperAdmin 超级管理员（系统拥有者），拥有所有权限
	RoleSuperAdmin = "super_admin"
	// RoleAdmin 普通管理员，具备后台管理和内容管理权限
	RoleAdmin = "admin"
	// RoleUser 普通用户（默认角色）
	RoleUser = "user"
)

// IsAdminRole 判断是否为具备管理员权限的角色
// 包含 super_admin 和 admin 两种角色
func IsAdminRole(role string) bool {
	return role == RoleAdmin || role == RoleSuperAdmin
}
