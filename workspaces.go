package coze

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type workspace struct {
	client *httpClient
}

func newWorkspace(client *httpClient) *workspace {
	return &workspace{client: client}
}

func (r *workspace) List(ctx context.Context, req *ListWorkspaceReq) (*NumberPaged[Workspace], error) {
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	return NewNumberPaged[Workspace](
		func(request *PageRequest) (*PageResponse[Workspace], error) {
			uri := "/v1/workspaces"
			resp := &listWorkspaceResp{}
			err := r.client.Request(ctx, http.MethodGet, uri, nil, resp,
				withHTTPQuery("page_num", strconv.Itoa(request.PageNum)),
				withHTTPQuery("page_size", strconv.Itoa(request.PageSize)))
			if err != nil {
				return nil, err
			}
			return &PageResponse[Workspace]{
				Total:   resp.Data.TotalCount,
				HasMore: len(resp.Data.Workspaces) >= request.PageSize,
				Data:    resp.Data.Workspaces,
				LogID:   resp.LogID,
			}, nil
		}, req.PageSize, req.PageNum)
}

// ListWorkspaceReq represents the request parameters for listing workspaces
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

// listWorkspaceResp represents the response for listing workspaces
type listWorkspaceResp struct {
	baseResponse
	Data *ListWorkspaceResp
}

// ListWorkspaceResp represents the response for listing workspaces
type ListWorkspaceResp struct {
	baseModel
	TotalCount int          `json:"total_count"`
	Workspaces []*Workspace `json:"workspaces"`
}

// Workspace represents workspace information
type Workspace struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	IconUrl       string            `json:"icon_url"`
	RoleType      WorkspaceRoleType `json:"role_type"`
	WorkspaceType WorkspaceType     `json:"workspace_type"`
}

// WorkspaceRoleType represents the workspace role type
type WorkspaceRoleType string

const (
	WorkspaceRoleTypeOwner  WorkspaceRoleType = "owner"
	WorkspaceRoleTypeAdmin  WorkspaceRoleType = "admin"
	WorkspaceRoleTypeMember WorkspaceRoleType = "member"
)

func (t WorkspaceRoleType) String() string {
	return string(t)
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

// WorkspaceType represents the workspace type
type WorkspaceType string

const (
	WorkspaceTypePersonal WorkspaceType = "personal"
	WorkspaceTypeTeam     WorkspaceType = "team"
)

func (t WorkspaceType) String() string {
	return string(t)
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
