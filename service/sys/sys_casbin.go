/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 16:09
 */
package sys

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"xz-go-frame/global"

	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

/*
_ "github.com/go-sql-driver/mysql"
这里的_（下划线）是一个特殊的导入方式，被称为“空白导入”（Blank Import）。
当使用这种方式导入包时，该包中的所有init函数会被执行，但是包中的其他公开的函数、类型、变量等不能直接使用。
这种导入方式通常用于仅需要执行包的初始化代码，而不需要直接使用包中的其他内容的情况。
*/

// Casbin info structure
type CasbinInfo struct {
	Path   string `json:"path"`   // 路径
	Method string `json:"method"` // 方法
}

type CasbinService struct{}

func (casbinService *CasbinService) UpdateCasbin(AuthorityID uint, casbinInfos []CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	casbinService.ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, v := range casbinInfos {
		rules = append(rules, []string{authorityId, v.Path, v.Method})
	}
	e := casbinService.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	err := e.InvalidateCache()
	if err != nil {
		return err
	}
	return nil
}

// @description: API更新随动
func (casbinService *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.XZ_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := casbinService.Casbin()
	err = e.InvalidateCache()
	if err != nil {
		return err
	}
	return err
}

// @description: 获取权限列表
func (casbinService *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []CasbinInfo) {
	e := casbinService.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// @description: 清除匹配的权限
func (casbinService *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := casbinService.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

// 持久化到数据库  引入自定义规则
var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (casbinService *CasbinService) CasbinFile() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(global.XZ_DB)
		if err != nil {
			global.Log.Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}
		m, err := model.NewModelFromFile("./conf/rbac_models.conf")
		if err != nil {
			global.Log.Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()

		// 设置用户root的角色为superadmin
		syncedCachedEnforcer.AddRoleForUser("root", "superadmin")
		// 添加自定义函数
		syncedCachedEnforcer.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (interface{}, error) {
			// 获取用户名
			username := arguments[0].(string)
			// 检查用户名的角色是否为superadmin
			return syncedCachedEnforcer.HasRoleForUser(username, "superadmin")
		})

		// 添加自定义函数
		// equals(r.sub, p.sub)
		syncedCachedEnforcer.AddFunction("equals", func(arguments ...interface{}) (interface{}, error) {
			// 获取用户名
			args1 := arguments[0].(string)
			args2 := arguments[1].(string)
			// 检查用户名的角色是否为superadmin
			return strings.EqualFold(args1, args2), nil
		})

	})
	return syncedCachedEnforcer
}

func (casbinService *CasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(global.XZ_DB)
		if err != nil {
			global.Log.Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			global.Log.Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}
