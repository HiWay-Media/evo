package evo

import (
	"errors"
	"fmt"

	"github.com/getevo/evo/lib/log"
	"github.com/getevo/evo/lib/validate"
	"gorm.io/gorm"
)

var memoryRolePermissions = cMap{}

func updateRolePermissions() {
	memoryRolePermissions.Init()

	// Load all roles at once
	var roles []Role
	db.Find(&roles)

	// Load all role permissions at once
	var rolePerms []RolePermission
	db.Find(&rolePerms)

	// Load all permissions at once
	var permissions []Permission
	db.Find(&permissions)

	// Create maps for quick lookup
	rolePermMap := make(map[uint][]RolePermission)
	for _, rp := range rolePerms {
		rolePermMap[rp.RoleID] = append(rolePermMap[rp.RoleID], rp)
	}

	permissionMap := make(map[uint]Permission)
	for _, perm := range permissions {
		permissionMap[perm.ID] = perm
	}

	// Process roles and their permissions
	for _, role := range roles {
		var perms Permissions
		for _, rp := range rolePermMap[role.ID] {
			if perm, exists := permissionMap[rp.PermissionID]; exists {
				perms = append(perms, perm)
			} else {
				log.Warning("Roles: found inconsistency, automatically remove permission id %d to fix.", rp.PermissionID)
				db.Delete(&RolePermission{}, "id = ?", rp.ID)
			}
		}
		memoryRolePermissions.Set(role.ID, &perms)
	}
}

// AfterFind after find event
func (r *Role) AfterFind(tx *gorm.DB) (err error) {
	perms := memoryRolePermissions.Get(r.ID)
	if perms != nil {
		r.Permission = perms.(*Permissions)
		r.PermissionSet = []string{}
		for _, item := range *r.Permission {
			r.PermissionSet = append(r.PermissionSet, item.CodeName)
		}
		return
	}
	r.Permission = &Permissions{}
	r.PermissionSet = []string{}
	return
}

// Save save role instance
func (r *Role) Save() error {
	var perms Permissions
	for _, item := range r.PermissionSet {
		var perm Permission
		if errors.Is(db.Where("code_name = ?", item).Take(&perm).Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("permission not found")
		}
		perms = append(perms, perm)
	}
	err := validate.Validate(r)
	if err != nil {
		return err
	}
	temp := Role{}
	if !errors.Is(db.Where("code_name = ?", r.CodeName).Find(&temp).Error, gorm.ErrRecordNotFound) {
		if r.ID == 0 || (r.ID > 0 && r.ID != temp.ID) {
			return fmt.Errorf("codename exist")
		}
	}
	if r.ID > 0 {
		r.SetPermission(perms)
		return db.Save(&r).Error
	} else {
		err := db.Create(&r).Error
		if err != nil {
			return err
		}
		return r.SetPermission(perms)
	}
}

// HasPerm check if role has permission
func (r *Role) HasPerm(v string) bool {
	for _, item := range r.PermissionSet {
		if item == v {
			return true
		}
	}
	return false
}

// SetPermission set role permission
func (r *Role) SetPermission(permissions Permissions) error {

	var listId []uint
	for k, item := range permissions {
		var perm RolePermission
		if errors.Is(db.Where("role_id = ? AND permission_id = ?", r.ID, item.ID).Take(&perm).Error, gorm.ErrRecordNotFound) {
			err := db.Create(&RolePermission{RoleID: r.ID, PermissionID: item.ID}).Error
			if err != nil {
				return err
			}
		}
		listId = append(listId, permissions[k].ID)
	}

	memoryRolePermissions.Set(r.ID, &permissions)
	return db.Delete(RolePermission{}, "role_id = ? AND permission_id NOT IN (?)", r.ID, listId).Error
}
