package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
)

type RouteInfo struct {
	Path      string      `json:"path"`
	Name      string      `json:"name"`
	Component string      `json:"component"`
	Redirect  string      `json:"redirect"`
	Hidden    bool        `json:"hidden"`
	Meta      MetaInfo    `json:"meta"`
	Children  []RouteInfo `json:"children,omitempty"`
}

type MetaInfo struct {
	Title string   `json:"title"`
	Icon  string   `json:"icon,omitempty"`
	Affix bool     `json:"affix,omitempty"`
	Roles []string `json:"roles,omitempty"` // 可选：限制访问角色
}

var fullRoutes = []RouteInfo{
	{
		Path:      "/",
		Name:      "常用",
		Component: "layout/index",
		Redirect:  "dashboard",
		Meta:      MetaInfo{Roles: []string{"admin", "teacher", "student"}},
		Children: []RouteInfo{
			{
				Path:      "dashboard",
				Name:      "工作台",
				Component: "views/dashboard/index",
				Meta:      MetaInfo{Title: "工作台", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
			},
			{
				Path:      "users",
				Name:      "用户管理",
				Component: "views/users/index",
				Meta:      MetaInfo{Title: "用户管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
			},
			{
				Path:      "course",
				Name:      "课程管理",
				Component: "views/course/index",
				Meta:      MetaInfo{Title: "课程管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
			},
			{
				Path:      "classes",
				Name:      "班级管理",
				Component: "views/classes/index",
				Meta:      MetaInfo{Title: "班级管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
			},
		},
	},
	{
		Path:      "dashboard",
		Name:      "工作台",
		Component: "views/dashboard/index",
		Meta:      MetaInfo{Title: "工作台", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
	},
	{
		Path:      "users",
		Name:      "用户管理",
		Component: "views/users/index",
		Meta:      MetaInfo{Title: "用户管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
	},
	{
		Path:      "course",
		Name:      "课程管理",
		Component: "views/course/index",
		Meta:      MetaInfo{Title: "课程管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
	},
	{
		Path:      "classes",
		Name:      "班级管理",
		Component: "views/classes/index",
		Meta:      MetaInfo{Title: "班级管理", Icon: "dashboard", Affix: true, Roles: []string{"admin", "teacher", "student"}},
	},
	{
		Path:      "/login",
		Name:      "登录",
		Component: "views/login/index",
		Hidden:    true,
		Meta:      MetaInfo{Roles: []string{"admin", "teacher", "student"}},
	},
	{
		Path:      "*",
		Name:      "404",
		Component: "views/404/index",
		Hidden:    true,
		Meta:      MetaInfo{Roles: []string{"admin", "teacher", "student"}},
	},
}

// GetRoutesHandler 下发路由表
// @Summary 获取路由配置
// @Tags User
// @Param roles query string true "角色数组，JSON格式"
// @Success 200 {object} basic.Resp{data=[]RouteInfo}
// @Router /api/v1/user/routes [get]
func GetRoutesHandler(c *gin.Context) {
	roleStr := c.Query("roles")
	if roleStr == "" {
		basic.RequestFailure(c, "缺少角色参数")
		return
	}

	var roles []string
	if err := json.Unmarshal([]byte(roleStr), &roles); err != nil {
		basic.RequestFailure(c, "角色参数格式错误，应为 JSON 数组")
		return
	}

	filtered := filterRoutesByRoles(fullRoutes, roles)
	basic.Success(c, filtered)
}
func filterRoutesByRoles(routes []RouteInfo, roles []string) []RouteInfo {
	roleMap := make(map[string]bool)
	for _, r := range roles {
		roleMap[r] = true
	}

	var result []RouteInfo
	for _, route := range routes {
		// 没有限制的，或者包含当前角色
		if len(route.Meta.Roles) == 0 || hasCommonRole(route.Meta.Roles, roleMap) {
			// 递归过滤 children
			newRoute := route
			if len(route.Children) > 0 {
				newRoute.Children = filterRoutesByRoles(route.Children, roles)
			}
			result = append(result, newRoute)
		}
	}
	return result
}

func hasCommonRole(routeRoles []string, roleMap map[string]bool) bool {
	for _, r := range routeRoles {
		if roleMap[r] {
			return true
		}
	}
	return false
}
