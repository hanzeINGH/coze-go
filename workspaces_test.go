package coze

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaces(t *testing.T) {
	// Test List method
	t.Run("List workspaces success", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify request method and path
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/v1/workspaces", req.URL.Path)

				// Verify query parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
				return mockResponse(http.StatusOK, &listWorkspaceResp{
					baseResponse: baseResponse{
						LogID: "test_log_id",
					},
					Data: &ListWorkspaceResp{
						TotalCount: 2,
						Workspaces: []*Workspace{
							{
								ID:            "ws1",
								Name:          "Workspace 1",
								IconUrl:       "https://example.com/icon1.png",
								RoleType:      WorkspaceRoleTypeOwner,
								WorkspaceType: WorkspaceTypePersonal,
							},
							{
								ID:            "ws2",
								Name:          "Workspace 2",
								IconUrl:       "https://example.com/icon2.png",
								RoleType:      WorkspaceRoleTypeAdmin,
								WorkspaceType: WorkspaceTypeTeam,
							},
						},
					},
				})
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, "https://api.coze.com")
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), &ListWorkspaceReq{
			PageNum:  1,
			PageSize: 20,
		})

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		items := paged.Items()
		require.Len(t, items, 2)

		// Verify first workspace
		assert.Equal(t, "ws1", items[0].ID)
		assert.Equal(t, "Workspace 1", items[0].Name)
		assert.Equal(t, "https://example.com/icon1.png", items[0].IconUrl)
		assert.Equal(t, WorkspaceRoleTypeOwner, items[0].RoleType)
		assert.Equal(t, WorkspaceTypePersonal, items[0].WorkspaceType)

		// Verify second workspace
		assert.Equal(t, "ws2", items[1].ID)
		assert.Equal(t, "Workspace 2", items[1].Name)
		assert.Equal(t, "https://example.com/icon2.png", items[1].IconUrl)
		assert.Equal(t, WorkspaceRoleTypeAdmin, items[1].RoleType)
		assert.Equal(t, WorkspaceTypeTeam, items[1].WorkspaceType)
	})

	// Test List method with default pagination
	t.Run("List workspaces with default pagination", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Verify default pagination parameters
				assert.Equal(t, "1", req.URL.Query().Get("page_num"))
				assert.Equal(t, "20", req.URL.Query().Get("page_size"))

				// Return mock response
				return mockResponse(http.StatusOK, &listWorkspaceResp{
					baseResponse: baseResponse{
						LogID: "test_log_id",
					},
					Data: &ListWorkspaceResp{
						TotalCount: 0,
						Workspaces: []*Workspace{},
					},
				})
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, "https://api.coze.com")
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), NewListWorkspaceReq())

		require.NoError(t, err)
		assert.False(t, paged.HasMore())
		assert.Empty(t, paged.Items())
	})

	// Test List method with error
	t.Run("List workspaces with error", func(t *testing.T) {
		mockTransport := &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				// Return error response
				return mockResponse(http.StatusBadRequest, &baseResponse{
					LogID: "test_error_log_id",
				})
			},
		}

		core := newCore(&http.Client{Transport: mockTransport}, "https://api.coze.com")
		workspaces := newWorkspace(core)

		paged, err := workspaces.List(context.Background(), &ListWorkspaceReq{
			PageNum:  1,
			PageSize: 20,
		})

		require.Error(t, err)
		assert.Nil(t, paged)
	})
}

func TestWorkspaceRoleType(t *testing.T) {
	t.Run("WorkspaceRoleType String method", func(t *testing.T) {
		assert.Equal(t, "owner", WorkspaceRoleTypeOwner.String())
		assert.Equal(t, "admin", WorkspaceRoleTypeAdmin.String())
		assert.Equal(t, "member", WorkspaceRoleTypeMember.String())
	})

	t.Run("WorkspaceRoleType UnmarshalJSON", func(t *testing.T) {
		var roleType WorkspaceRoleType

		// Test valid role types
		err := json.Unmarshal([]byte(`"owner"`), &roleType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceRoleTypeOwner, roleType)

		err = json.Unmarshal([]byte(`"admin"`), &roleType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceRoleTypeAdmin, roleType)

		err = json.Unmarshal([]byte(`"member"`), &roleType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceRoleTypeMember, roleType)

		// Test unknown role type
		err = json.Unmarshal([]byte(`"unknown"`), &roleType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceRoleType("unknown"), roleType)

		// Test invalid JSON
		err = json.Unmarshal([]byte(`invalid`), &roleType)
		require.Error(t, err)
	})
}

func TestWorkspaceType(t *testing.T) {
	t.Run("WorkspaceType String method", func(t *testing.T) {
		assert.Equal(t, "personal", WorkspaceTypePersonal.String())
		assert.Equal(t, "team", WorkspaceTypeTeam.String())
	})

	t.Run("WorkspaceType UnmarshalJSON", func(t *testing.T) {
		var workspaceType WorkspaceType

		// Test valid workspace types
		err := json.Unmarshal([]byte(`"personal"`), &workspaceType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceTypePersonal, workspaceType)

		err = json.Unmarshal([]byte(`"team"`), &workspaceType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceTypeTeam, workspaceType)

		// Test unknown workspace type
		err = json.Unmarshal([]byte(`"unknown"`), &workspaceType)
		require.NoError(t, err)
		assert.Equal(t, WorkspaceType("unknown"), workspaceType)

		// Test invalid JSON
		err = json.Unmarshal([]byte(`invalid`), &workspaceType)
		require.Error(t, err)
	})
}
