package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/syspulse/common"
	"github.com/syspulse/model"
	"github.com/syspulse/restful/server/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func login(username string, passwd string) (string, *model.User, error) {
	userInfoLst := model.DBSelect(`SELECT 
    u.id 'u_id', u.username, r.identity as 'r_identity'
FROM
    (SELECT 
        id, username
    FROM
        user
    WHERE
        username = ?
            AND passwd = ?
            AND is_active = 1) u
        LEFT JOIN
    user_role ur ON u.id = ur.user_id
        LEFT JOIN
    role r ON ur.role_id = r.id`, username, passwd)

	if len(userInfoLst) < 1 {
		return "", nil, &common.InsightException{Code: 500, Msg: fmt.Sprintf("no user with username '%s'", username)}
	}

	user := new(model.User)
	for _, userInfo := range userInfoLst {
		if user.RoleLst == nil {
			user.ID = userInfo["u_id"].(int64)
			user.Username = string(userInfo["username"].([]uint8))
			user.IsActive = true
			if val, exists := userInfo["r_identity"]; exists && val != nil {
				user.RoleLst = []*model.Role{{
					Identity: string(val.([]uint8)),
				}}
			}

		} else {
			if val, exists := userInfo["r_identity"]; exists && val != nil {
				user.RoleLst = append(user.RoleLst, &model.Role{
					Identity: string(val.([]uint8)),
				})
			}

		}
	}

	token := uuid.NewString()
	bytes, _ := json.Marshal(user)
	model.CacheSet(token, string(bytes), common.SysArgs.Session.Expiration*time.Minute)
	return token, user, nil
}

func UserLogin(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from: ", r)
			common.PrintStackTrace()
		}
	}()

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at user login.", zap.Error(err))
		return
	}
	var info = make(map[string]string)
	err = json.Unmarshal(body, &info)
	if err != nil {
		zap.L().Error("error parse json body at user login.", zap.Error(err))
		return
	}

	username := info["username"]
	passwd := info["passwd"]
	token, user, err := login(username, passwd)

	if err != nil {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusUnauthorized, Msg: err.Error()})
	} else {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: map[string]interface{}{
			"token": token,
			"user":  user,
		}, Msg: "success"})
	}
}

func saveUser(u *model.User) int64 {
	userId := int64(0)
	model.ExecInTransaction(func(tx *sql.Tx) bool {
		sql1 := "insert into user(`username`, `passwd`, `is_active`, `createTimestamp`, `updateTimestamp`) value(?, ?, ?, ?, ?)"
		result, err := tx.Exec(sql1, u.Username, u.Passwd, u.IsActive, u.CreateTimestamp, u.UpdateTimestamp)
		if err != nil {
			zap.L().Error("error insert user info into table user.", zap.Error(err))
			return false
		}
		uid, err := result.LastInsertId()
		if err != nil {
			zap.L().Error("error get user id.", zap.Error(err))
			return false
		}
		sql2 := "insert into user_role(`user_id`, `role_id`) value(?, ?)"
		stmt, err := model.SqlDB.Prepare(sql2)
		if err != nil {
			zap.L().Error("error prepare sql when save relationship between user and role.", zap.Error(err))
		}
		defer stmt.Close()
		for _, role := range u.RoleLst {
			_, err := tx.Stmt(stmt).Exec(uid, role.ID)
			if err != nil {
				zap.L().Error("can't save relationship between user and role: ", zap.Error(err))
			}
		}
		userId = uid

		return true
	})
	return userId
}

func CreateUser(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at user login.", zap.Error(err))
		return
	}
	user := new(model.User)
	err = json.Unmarshal(body, user)
	if err != nil {
		zap.L().Error("error parse json body at user login.", zap.Error(err))
		return
	}
	timestamp := time.Now().UnixMilli()
	user.CreateTimestamp = timestamp
	user.UpdateTimestamp = timestamp
	user.ID = saveUser(user)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: user, Msg: "success"})
}

func saveRole(role *model.Role) int64 {
	var roleId int64
	model.ExecInTransaction(func(tx *sql.Tx) bool {
		query := "insert into role(`name`, `identity`, `createTimestamp`, `updateTimestamp`) values (?, ?, ?, ?)"
		result, err := tx.Exec(query, role.Name, role.Identity, role.CreateTimestamp, role.UpdateTimestamp)
		if err != nil {
			zap.L().Error("error insert role", zap.Error(err))
			return false
		}
		id, err := result.LastInsertId()
		if err != nil {
			zap.L().Error("error get role last insert id", zap.Error(err))
			return false
		}

		sql := "insert into role_permission(`role_id`, `permission_id`) value(?, ?)"
		stmt, err := model.SqlDB.Prepare(sql)
		if err != nil {
			zap.L().Error("can't get statement obj: ", zap.Error(err))
		}
		defer stmt.Close()

		for _, permission := range role.PermissionLst {
			_, err := tx.Stmt(stmt).Exec(id, permission.ID)
			if err != nil {
				zap.L().Error("can't save relationship between role and permission: ", zap.Error(err))
			}
		}

		roleId = id
		return true
	})
	return roleId
}

func CreateRole(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at create role.", zap.Error(err))
		return
	}
	role := new(model.Role)
	err = json.Unmarshal(body, role)
	if err != nil {
		zap.L().Error("error parse json body at create role.", zap.Error(err))
		return
	}
	timestamp := time.Now().UnixMilli()
	role.CreateTimestamp = timestamp
	role.UpdateTimestamp = timestamp
	role.ID = saveRole(role)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: role, Msg: "success"})
}

func ListUserByPage(page, pageSize int) []*model.User {
	first := page * pageSize
	s := `SELECT 
    user.id,
    user.username,
    user.is_active,
    user.createTimestamp,
    user.updateTimestamp,
    role.id 'rid',
    role.name 'rname'
FROM
    (SELECT 
        id,
            username,
            is_active,
            createTimestamp,
            updateTimestamp
    FROM
        user
	WHERE
	    id != 0
    ORDER BY updateTimestamp DESC , id DESC
    LIMIT ? , ?) user
        LEFT JOIN
    user_role ur ON user.id = ur.user_id
        LEFT JOIN
    role ON ur.role_id = role.id`
	userInfoLst := model.DBSelect(s, first, pageSize)
	lst := make([]*model.User, 0)
	for _, userInfo := range userInfoLst {
		userId := userInfo["id"].(int64)
		currSize := len(lst)
		if currSize > 0 && lst[currSize-1].ID == userId {
			last := lst[currSize-1]
			role := &model.Role{
				ID:   userInfo["rid"].(int64),
				Name: string(userInfo["rname"].([]uint8)),
			}
			last.RoleLst = append(last.RoleLst, role)
		} else {
			user := new(model.User)
			user.ID = userId
			user.Username = string(userInfo["username"].([]uint8))
			user.IsActive = userInfo["is_active"].(int64) == 1
			user.CreateTimestamp = userInfo["createTimestamp"].(int64)
			user.UpdateTimestamp = userInfo["updateTimestamp"].(int64)
			if userInfo["rid"] != nil {
				role := &model.Role{
					ID:   userInfo["rid"].(int64),
					Name: string(userInfo["rname"].([]uint8)),
				}
				user.RoleLst = []*model.Role{role}
			}

			lst = append(lst, user)
		}

	}
	return lst
}

func GetUserPage(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("pageNum"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}
	userLst := ListUserByPage(page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: userLst, Msg: "success"})
}

func DeleteUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "user id is not a number."})
		return
	}

	s := "delete from user where id = ?"
	model.DBDelete(s, userId)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func ListPermissionByPage(page, pageSize int) []*model.Permission {
	first := page * pageSize
	lst := make([]*model.Permission, 0)
	s := `select
  p1.id, p1.identity, p1.name, p1.method, p1.url, p1.createTimestamp, p1.updateTimestamp
from 
  permission p1
order by 
  p1.updateTimestamp desc, p1.id desc limit ?, ?`
	permissionLst := model.DBSelect(s, first, pageSize)
	for _, permissionInfo := range permissionLst {
		permission := new(model.Permission)
		permission.ID = permissionInfo["id"].(int64)
		permission.Identity = string(permissionInfo["identity"].([]uint8))
		permission.Name = string(permissionInfo["name"].([]uint8))
		permission.Url = string(permissionInfo["url"].([]uint8))
		permission.Method = string(permissionInfo["method"].([]uint8))
		permission.CreateTimestamp = permissionInfo["createTimestamp"].(int64)
		permission.UpdateTimestamp = permissionInfo["updateTimestamp"].(int64)
		lst = append(lst, permission)
	}
	return lst
}

func GetPermissionLst(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("pageNum"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}
	permissionLst := ListPermissionByPage(page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: permissionLst, Msg: "success"})
}

func savePermission(p *model.Permission) int64 {
	s := "insert into permission(identity, name, method, url, createTimestamp, updateTimestamp) value(?, ?, ?, ?, ?, ?)"
	return model.DBInsert(s, p.Identity, p.Name, p.Method, p.Url, p.CreateTimestamp, p.UpdateTimestamp)
}

func CreatePermission(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at user login.", zap.Error(err))
		return
	}
	permission := new(model.Permission)
	err = json.Unmarshal(body, permission)
	if err != nil {
		zap.L().Error("error parse json body at user login.", zap.Error(err))
		return
	}
	timestamp := time.Now().UnixMilli()
	permission.CreateTimestamp = timestamp
	permission.UpdateTimestamp = timestamp
	permission.ID = savePermission(permission)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: permission, Msg: "success"})
}

func DeletePermission(ctx *gin.Context) {
	prmId, err := strconv.ParseInt(ctx.Param("prmId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "permission id is not a number."})
		return
	}

	model.ExecInTransaction(func(tx *sql.Tx) bool {

		sql1 := "delete from role_permission where permission_id = ?"
		_, err := tx.Exec(sql1, prmId)
		if err != nil {
			zap.L().Error("error delete from role_permission.", zap.Error(err))
		}

		sql2 := "delete from menu_permission where permission_id = ?"
		_, err = tx.Exec(sql2, prmId)
		if err != nil {
			zap.L().Error("error delete from role_permission.", zap.Error(err))
		}

		sql3 := "delete from permission where id = ?"
		_, err = tx.Exec(sql3, prmId)
		if err != nil {
			zap.L().Error("error delete from role_permission.", zap.Error(err))
		}

		return true
	})
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func ListRoleByPage(page, pageSize int) []*model.Role {
	first := page * pageSize
	lst := make([]*model.Role, 0)
	s := `SELECT 
    r.id as 'rid',
    r.identity,
    r.name,
    r.createTimestamp,
    r.updateTimestamp,
    p.id as 'pid',
    p.name as 'pname'
FROM
    (SELECT 
        role.id,
            role.identity,
            role.name,
            role.createTimestamp,
            role.updateTimestamp
    FROM
        role
	ORDER BY updateTimestamp DESC , id DESC
    LIMIT ? , ?) AS r
        LEFT JOIN
    role_permission AS rp ON r.id = rp.role_id
        LEFT JOIN
    permission AS p ON rp.permission_id = p.id`
	roleInfoLst := model.DBSelect(s, first, pageSize)
	for _, roleInfo := range roleInfoLst {
		roleId := roleInfo["rid"].(int64)
		currSize := len(lst)
		if currSize > 0 && lst[currSize-1].ID == roleId {
			last := lst[currSize-1]
			permission := &model.Permission{
				ID:   roleInfo["pid"].(int64),
				Name: string(roleInfo["pname"].([]uint8)),
			}
			last.PermissionLst = append(last.PermissionLst, permission)
		} else {
			role := new(model.Role)
			role.ID = roleId
			role.Identity = string(roleInfo["identity"].([]uint8))
			role.Name = string(roleInfo["name"].([]uint8))
			role.CreateTimestamp = roleInfo["createTimestamp"].(int64)
			role.UpdateTimestamp = roleInfo["updateTimestamp"].(int64)
			if roleInfo["pid"] != nil {
				permission := &model.Permission{
					ID:   roleInfo["pid"].(int64),
					Name: string(roleInfo["pname"].([]uint8)),
				}
				role.PermissionLst = []*model.Permission{permission}
			}

			lst = append(lst, role)
		}

	}
	return lst
}

func GetRoleByPage(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("pageNum"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}
	permissionLst := ListRoleByPage(page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: permissionLst, Msg: "success"})
}

func DeleteRole(ctx *gin.Context) {
	roleId, err := strconv.ParseInt(ctx.Param("roleId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "user id is not a number."})
		return
	}

	model.ExecInTransaction(func(tx *sql.Tx) bool {
		s1 := "delete from role_permission where role_id = ?"
		tx.Exec(s1, roleId)
		s := "delete from role where id = ?"
		tx.Exec(s, roleId)
		return true
	})
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func CreateMenu(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		zap.L().Error("error read body at user login.", zap.Error(err))
		return
	}
	menu := new(model.Menu)
	err = json.Unmarshal(body, menu)
	if err != nil {
		zap.L().Error("error parse json body at user login.", zap.Error(err))
		return
	}
	timestamp := time.Now().UnixMilli()
	menu.CreateTimestamp = timestamp
	menu.UpdateTimestamp = timestamp
	menu.ID = saveMenu(menu)
	if menu.ID != 0 {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: menu, Msg: "success"})
	} else {
		ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusInternalServerError, Data: menu, Msg: "error"})
	}
}

func saveMenu(menu *model.Menu) int64 {
	var menuId int64
	model.ExecInTransaction(func(tx *sql.Tx) bool {
		query := "insert into menu(`title`, `identity`, `type`, `url`, `index`, `parent`, `createTimestamp`, `updateTimestamp`) values (?, ?, ?, ?, ?, ?, ?, ?)"
		result, err := tx.Exec(query, menu.Title, menu.Identity, menu.Type, menu.Url, menu.Index, menu.ParentId, menu.CreateTimestamp, menu.UpdateTimestamp)
		if err != nil {
			zap.L().Error("error insert menu", zap.Error(err))
			return false
		}
		id, err := result.LastInsertId()
		if err != nil {
			zap.L().Error("error get menu last insert id", zap.Error(err))
			return false
		}

		sql := "insert into menu_permission(`menu_id`, `permission_id`) value(?, ?)"
		stmt, err := model.SqlDB.Prepare(sql)
		if err != nil {
			zap.L().Error("can't get statement obj: ", zap.Error(err))
		}
		defer stmt.Close()

		for _, permission := range menu.PermissionLst {
			_, err := tx.Stmt(stmt).Exec(id, permission.ID)
			if err != nil {
				zap.L().Error("can't save relationship between menu and permission: ", zap.Error(err))
			}
		}

		menuId = id
		return true
	})
	return menuId
}

func DeleteMenu(ctx *gin.Context) {
	menuId, err := strconv.ParseInt(ctx.Param("menuId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "user id is not a number."})
		return
	}

	model.ExecInTransaction(func(tx *sql.Tx) bool {
		s1 := "delete from menu_permission where menu_id = ?"
		tx.Exec(s1, menuId)
		s := "delete from menu where id = ?"
		tx.Exec(s, menuId)
		return true
	})
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Msg: "success"})
}

func ListMenuByPage(page, pageSize int) []*model.Menu {
	first := page * pageSize
	lst := make([]*model.Menu, 0)
	s := `SELECT 
    m.id AS 'mid',
    m.identity,
    m.title AS 'm_title',
    m.url,
    m.type,
    m.index,
    m.parentId,
    m.parentTitle,
    m.createTimestamp,
    m.updateTimestamp,
    p.id AS 'pid',
    p.name
FROM
    (SELECT 
        menu.id,
            menu.identity,
            menu.title,
            menu.url,
            menu.type,
            menu.index,
            menu.parent,
            menu.createTimestamp,
            menu.updateTimestamp,
            menu2.id "parentId",
            menu2.title "parentTitle"
    FROM
        menu menu
    LEFT JOIN menu menu2 ON menu.parent = menu2.id
    ORDER BY menu.updateTimestamp DESC , menu.id DESC
    LIMIT ? , ?) AS m
        LEFT JOIN
    menu_permission AS rp ON m.id = rp.menu_id
        LEFT JOIN
    permission AS p ON rp.permission_id = p.id`
	menuInfoLst := model.DBSelect(s, first, pageSize)
	for _, menuInfo := range menuInfoLst {
		menuId := menuInfo["mid"].(int64)
		currSize := len(lst)
		if currSize > 0 && lst[currSize-1].ID == menuId {
			last := lst[currSize-1]
			permission := &model.Permission{
				ID:   menuInfo["pid"].(int64),
				Name: string(menuInfo["name"].([]uint8)),
			}
			last.PermissionLst = append(last.PermissionLst, permission)
		} else {
			menu := new(model.Menu)
			menu.ID = menuId
			menu.Identity = string(menuInfo["identity"].([]uint8))
			menu.Title = string(menuInfo["m_title"].([]uint8))
			menu.Url = string(menuInfo["url"].([]uint8))
			menu.Type = string(menuInfo["type"].([]uint8))
			menu.Index = uint8(menuInfo["index"].(int64))
			parentId, exists := menuInfo["parentId"]
			if exists && parentId != nil {
				menu.ParentId = parentId.(int64)
				menu.ParentTitle = string(menuInfo["parentTitle"].([]uint8))
			}
			menu.CreateTimestamp = menuInfo["createTimestamp"].(int64)
			menu.UpdateTimestamp = menuInfo["updateTimestamp"].(int64)
			if menuInfo["pid"] != nil {
				permission := &model.Permission{
					ID:   menuInfo["pid"].(int64),
					Name: string(menuInfo["name"].([]uint8)),
				}
				menu.PermissionLst = []*model.Permission{permission}
			}

			lst = append(lst, menu)
		}

	}
	return lst
}

func GetMenuByPage(ctx *gin.Context) {
	values := ctx.Request.URL.Query()
	page, err := strconv.Atoi(values.Get("pageNum"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "page is not a number."})
		return
	}
	pageSize, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.JsonResponse{Status: http.StatusBadRequest, Msg: "pageSize is not a number."})
		return
	}
	menuLst := ListMenuByPage(page, pageSize)
	ctx.JSON(http.StatusOK, response.JsonResponse{Status: http.StatusOK, Data: menuLst, Msg: "success"})
}
