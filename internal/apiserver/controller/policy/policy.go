// Package policy implements the policy handler.
package policy

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/meta"
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	v1 "iam-apiserver/api/apiserver/v1"
	"iam-apiserver/internal/apiserver/service"
)

// PolicyController handles requests for policy resource.
type PolicyController struct {
	srv service.Service
}

// NewPolicyController creates a policy handler.
func NewPolicyController() *PolicyController {
	return &PolicyController{
		srv: service.NewService(),
	}
}

// Create creates a new policy.
func (p *PolicyController) Create(c *gin.Context) {
	policy := &v1.Policy{}

	if err := c.ShouldBindJSON(policy); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	username := c.GetString(middleware.UsernameKey)
	err := p.srv.Policies().Create(c, username, policy)
	httputil.WriteResponse(c, err, policy)
}

// Get returns a policy by the policy identifier.
func (p *PolicyController) Get(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	policy, err := p.srv.Policies().Get(c, username, c.Param("name"))

	httputil.WriteResponse(c, err, policy)
}

// Update updates the policy by the policy identifier.
func (p *PolicyController) Update(c *gin.Context) {
	params := &v1.Policy{}

	if err := c.ShouldBindJSON(params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	username := c.GetString(middleware.UsernameKey)
	err := p.srv.Policies().Update(c, username, c.Param("name"), params)
	httputil.WriteResponse(c, err, nil)
}

// List returns all the policies.
func (p *PolicyController) List(c *gin.Context) {
	var opts meta.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}

	username := c.GetString(middleware.UsernameKey)
	policies, err := p.srv.Policies().List(c, username, opts)
	httputil.WriteResponse(c, err, policies)
}

// Delete deletes a policy by the policy identifier.
func (p *PolicyController) Delete(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	err := p.srv.Policies().Delete(c, username, c.Param("name"))
	httputil.WriteResponse(c, err, nil)
}

// DeleteCollection batch delete policies by username and secretIDs.
func (p *PolicyController) DeleteCollection(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	secretIDs := c.QueryArray("name")
	err := p.srv.Policies().DeleteCollection(c, username, secretIDs)
	httputil.WriteResponse(c, err, nil)
}
