package coze

import "encoding/json"

// ListWorkspaceReq 列表请求参数
type ListWorkspaceReq struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

func NewListWorkspaceReq() *ListWorkspaceReq {
	return &ListWorkspaceReq{
		PageNum:  1,
		PageSize: 20,
	}
}

// ListWorkspaceResp 列表响应
type ListWorkspaceResp struct {
	TotalCount int          `json:"total_count"`
	Workspaces []*Workspace `json:"workspaces"`
}

// Workspace 工作空间信息
type Workspace struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	IconUrl       string            `json:"icon_url"`
	RoleType      WorkspaceRoleType `json:"role_type"`
	WorkspaceType WorkspaceType     `json:"workspace_type"`
}

// WorkspaceRoleType 工作空间角色类型
type WorkspaceRoleType string

const (
	WorkspaceRoleTypeOwner  WorkspaceRoleType = "owner"
	WorkspaceRoleTypeAdmin  WorkspaceRoleType = "admin"
	WorkspaceRoleTypeMember WorkspaceRoleType = "member"
)

func (t WorkspaceRoleType) String() string {
	return string(t)
}

func (t WorkspaceRoleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *WorkspaceRoleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "owner":
		*t = WorkspaceRoleTypeOwner
	case "admin":
		*t = WorkspaceRoleTypeAdmin
	case "member":
		*t = WorkspaceRoleTypeMember
	default:
		*t = WorkspaceRoleType(s)
	}
	return nil
}

// WorkspaceType 工作空间类型
type WorkspaceType string

const (
	WorkspaceTypePersonal WorkspaceType = "personal"
	WorkspaceTypeTeam     WorkspaceType = "team"
)

func (t WorkspaceType) String() string {
	return string(t)
}

func (t WorkspaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *WorkspaceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "personal":
		*t = WorkspaceTypePersonal
	case "team":
		*t = WorkspaceTypeTeam
	default:
		*t = WorkspaceType(s)
	}
	return nil
}

type workspace struct {
}
