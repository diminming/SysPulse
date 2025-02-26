package server

import (
	"strings"

	"github.com/syspulse/model"
	"go.uber.org/zap"
)

func loadAuthConfig() {
	sql1 := "select prms.method, prms.url, prms.identity from permission prms"
	prmsInfoLst := model.DBSelect(sql1)
	for _, prmsInfo := range prmsInfoLst {
		key := "__req__" + string(prmsInfo["method"].([]uint8))
		model.CacheHSet(key, string(prmsInfo["url"].([]uint8)), string(prmsInfo["identity"].([]uint8)))
	}
	sql2 := "select role.identity r_identity, permission.identity p_identity from role left join role_permission on role.id = role_permission.role_id left join permission on role_permission.permission_id = permission.id"
	mappingLst := model.DBSelect(sql2)
	cache := make(map[string][]string)
	for _, mapping := range mappingLst {

		pIdentity, exists := mapping["p_identity"]
		if !exists || pIdentity == nil {
			continue
		}

		roleId := string(mapping["r_identity"].([]uint8))
		prmsId := string(mapping["p_identity"].([]uint8))

		val, exists := cache[roleId]
		if exists {
			cache[roleId] = append(val, prmsId)
		} else {
			cache[roleId] = []string{prmsId}
		}

	}

	for k, v := range cache {
		key := "__role__" + k
		value := strings.Join(v, ",")
		model.CacheSet(key, value, 0)
		zap.L().Info("mapping between role and peremission", zap.String("role identity", key), zap.String("permission identity", value))
	}

}

func BuildAuthCache() {
	loadAuthConfig()
}

func CheckAuth(url, method string, roleLst []string) bool {
	key := "__req__" + method
	cfg := model.CacheHGetAll(key)
	for urlPattern, permission := range cfg {
		result := (urlPattern == url)

		if result {
			for _, role := range roleLst {
				permissionString := model.CacheGet("__role__" + role)
				if strings.Contains(permissionString, permission) {
					return true
				}
			}
		}
	}
	return false
}
