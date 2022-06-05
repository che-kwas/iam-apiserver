// Package secret is the secret controller.
package secret

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

// SecretController handles requests for secret resource.
type SecretController struct {
	srv service.Service
	log *logger.Logger
}

// NewSecretController creates a secret controller.
func NewSecretController() *SecretController {
	return &SecretController{
		srv: service.NewService(),
		log: logger.L(),
	}
}

// Create creates a new secret.
func (s *SecretController) Create(c *gin.Context) {
	secret := &v1.Secret{}
	if err := c.ShouldBindJSON(secret); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	s.log.X(c).Infow("secret create param", "secret", secret)

	username := c.GetString(middleware.UsernameKey)
	err := s.srv.Secrets().Create(c, username, secret)
	httputil.WriteResponse(c, err, secret)
}

// Get returns a secret by the secret identifier.
func (s *SecretController) Get(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	s.log.X(c).Infow("secret get params", "name", name)

	secret, err := s.srv.Secrets().Get(c, username, name)

	httputil.WriteResponse(c, err, secret)
}

// Update updates the secret by the secret identifier.
func (s *SecretController) Update(c *gin.Context) {
	params := &v1.Secret{}
	if err := c.ShouldBindJSON(params); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	s.log.X(c).Infow("secret update params", "name", name, "params", params)

	err := s.srv.Secrets().Update(c, username, name, params)
	httputil.WriteResponse(c, err, nil)
}

// List returns all the secrets.
func (s *SecretController) List(c *gin.Context) {
	var opts meta.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	username := c.GetString(middleware.UsernameKey)
	s.log.X(c).Infow("secret list params", "offset", opts.Offset, "limit", opts.Limit)

	secrets, err := s.srv.Secrets().List(c, username, opts)
	httputil.WriteResponse(c, err, secrets)
}

// Delete deletes a secret by the secret identifier.
func (s *SecretController) Delete(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	name := c.Param("name")
	s.log.X(c).Infow("secret delete params", "name", name)

	err := s.srv.Secrets().Delete(c, username, name)
	httputil.WriteResponse(c, err, nil)
}

// DeleteCollection batch delete secrets by username and secretIDs.
func (s *SecretController) DeleteCollection(c *gin.Context) {
	username := c.GetString(middleware.UsernameKey)
	names := c.QueryArray("name")
	s.log.X(c).Infow("secret delete-collection params", "names", names)

	err := s.srv.Secrets().DeleteCollection(c, username, names)
	httputil.WriteResponse(c, err, nil)
}
