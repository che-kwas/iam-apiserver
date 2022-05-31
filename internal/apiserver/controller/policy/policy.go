// Package policy is the policy controller.
package policy

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/logger"
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
	log *logger.Logger
}

// NewPolicyController creates a policy controller.
func NewPolicyController() *PolicyController {
	return &PolicyController{
		srv: service.NewService(),
		log: logger.L(),
	}
}

// Create creates a new policy.
func (p *PolicyController) Create(c *gin.Context) {
	p.log.X(c).Info("policy create")
	policy := &v1.Policy{}

	if err := c.ShouldBindJSON(policy); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	username := c.GetString(middleware.UsernameKey)
	p.log.Debugf("policy create params: %s, %+v", username, policy)

	err := p.srv.Policies().Create(c, username, policy)
	httputil.WriteResponse(c, err, policy)
}

// Get returns a policy by the policy identifier.
func (p *PolicyController) Get(c *gin.Context) {
	p.log.X(c).Info("policy get")
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	p.log.Debugf("policy get params: %s, %s", username, name)

	policy, err := p.srv.Policies().Get(c, username, name)
	httputil.WriteResponse(c, err, policy)
}

// Update updates the policy by the policy identifier.
func (p *PolicyController) Update(c *gin.Context) {
	p.log.X(c).Info("policy update")
	params := &v1.Policy{}

	if err := c.ShouldBindJSON(params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	p.log.Debugf("policy update params: %s, %s, %+v", username, name, params)

	err := p.srv.Policies().Update(c, username, name, params)
	httputil.WriteResponse(c, err, nil)
}

// List returns all the policies.
func (p *PolicyController) List(c *gin.Context) {
	p.log.X(c).Info("policy list")
	var opts meta.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	username := c.GetString(middleware.UsernameKey)
	p.log.Debugf("policy list params: %s, %v", username, opts)

	policies, err := p.srv.Policies().List(c, username, opts)
	httputil.WriteResponse(c, err, policies)
}

// Delete deletes a policy by the policy identifier.
func (p *PolicyController) Delete(c *gin.Context) {
	p.log.X(c).Info("policy delete")
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	p.log.Debugf("policy delete params: %s, %s", username, name)

	err := p.srv.Policies().Delete(c, username, name)
	httputil.WriteResponse(c, err, nil)
}

// DeleteCollection batch delete policies by username and secretIDs.
func (p *PolicyController) DeleteCollection(c *gin.Context) {
	p.log.X(c).Info("policy delete-collection")
	username := c.GetString(middleware.UsernameKey)
	names := c.QueryArray("name")
	p.log.Debugf("policy delete-collection params: %s, %v", username, names)

	err := p.srv.Policies().DeleteCollection(c, username, names)
	httputil.WriteResponse(c, err, nil)
}
