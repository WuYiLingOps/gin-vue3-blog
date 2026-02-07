/*
 * 项目名称：blog-backend
 * 文件名称：role.go
 * 创建时间：2026-02-06 21:57:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户角色相关常量与工具方法（兼容层），统一管理角色字符串，避免魔法字符串散落各处。
 * 说明：
 *   1. 角色常量与判定逻辑的「真实实现」位于 blog-backend/constant/role.go 中。
 *   2. 本文件仅作为兼容层，对外转发 constant 包中的 RoleSuperAdmin / RoleAdmin / RoleUser 以及 IsAdminRole。
 *   3. 这样可以避免 util <-> repository 等包之间产生循环依赖，同时不破坏已有 import "blog-backend/util" 的旧代码。
 *   4. 后续新代码如果不依赖 util 其它工具函数，推荐直接引用 constant 包中的角色常量与方法。
 */
package util

import "blog-backend/constant"

// 兼容保留说明：
//   - 角色常量与工具方法已迁移至 constant 包，以避免 util<->repository 循环依赖。
//   - 这里做 re-export（转导出），避免外部已有代码引用 util.RoleXXX / util.IsAdminRole 时需要立即全量重构。
//   - 建议新代码优先直接依赖 constant 包，以便未来逐步弱化 util 的角色为纯工具包。

const (
	RoleSuperAdmin = constant.RoleSuperAdmin
	RoleAdmin      = constant.RoleAdmin
	RoleUser       = constant.RoleUser
)

func IsAdminRole(role string) bool { return constant.IsAdminRole(role) }
