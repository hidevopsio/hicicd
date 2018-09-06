package scm

import (
	"time"
)

type GroupMemberInterface interface {
	ListGroupMembers(token, baseUrl string, gid, uid int) (int,  error)
	GetGroupMember(token, baseUrl string, gid, uid int) (*GroupMember, error)
}

type GroupMember struct {
	ID               int        `json:"id"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	Name             string     `json:"name"`
	State            string     `json:"state"`
	CreatedAt        *time.Time `json:"created_at"`
	AccessLevelValue int        `json:"access_level_value"`
}
