// Code generated protoc-gen-go-gin. DO NOT EDIT.
// protoc-gen-go-gin 0.0.14

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
	app "github.com/go-eagle/eagle/pkg/app"
	errcode "github.com/go-eagle/eagle/pkg/errcode"
	"github.com/spf13/cast"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the eagle package it is being compiled against.

// context.
// metadata.
// gin.app.errcode.

type CompanyServiceHTTPServer interface {
	CreateCompany(context.Context, *CreateCompanyRequest) (*CreateCompanyReply, error)
	DeleteCompany(context.Context, *DeleteCompanyRequest) (*DeleteCompanyReply, error)
	GetCompany(context.Context, *GetCompanyRequest) (*GetCompanyReply, error)
	GetCompanyByAddress(context.Context, *GetCompanyByAddressRequest) (*GetCompanyReply, error)
	ListCompany(context.Context, *ListCompanyRequest) (*ListCompanyReply, error)
	UpdateCompany(context.Context, *UpdateCompanyRequest) (*UpdateCompanyReply, error)
}

func RegisterCompanyServiceHTTPServer(r gin.IRouter, srv CompanyServiceHTTPServer) {
	s := CompanyService{
		server: srv,
		router: r,
	}
	s.RegisterService()
}

type CompanyService struct {
	server CompanyServiceHTTPServer
	router gin.IRouter
}

func (s *CompanyService) CreateCompany_0(ctx *gin.Context) {
	var in CreateCompanyRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).CreateCompany(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) UpdateCompany_0(ctx *gin.Context) {
	var in UpdateCompanyRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// make sure the uri include :id
	in.Id = cast.ToInt32(ctx.Param("id"))


	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).UpdateCompany(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) DeleteCompany_0(ctx *gin.Context) {
	var in DeleteCompanyRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// make sure the uri include :id
	in.Id = cast.ToInt32(ctx.Param("id"))

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).DeleteCompany(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) GetCompany_0(ctx *gin.Context) {
	var in GetCompanyRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// make sure the uri include :id
	in.Id = cast.ToInt32(ctx.Param("id"))

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).GetCompany(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) GetCompanyByAddress_0(ctx *gin.Context) {
	var in GetCompanyByAddressRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	// make sure the uri include :id
	in.WalletAddress = ctx.Param("wallet_address")

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).GetCompanyByAddress(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) ListCompany_0(ctx *gin.Context) {
	var in ListCompanyRequest

	if err := ctx.ShouldBindQuery(&in); err != nil {
		app.Error(ctx, errcode.ErrInvalidParam.WithDetails(err.Error()))
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(CompanyServiceHTTPServer).ListCompany(newCtx, &in)
	if err != nil {
		app.Error(ctx, err)
		return
	}

	app.Success(ctx, out)
}

func (s *CompanyService) RegisterService() {
	s.router.Handle("POST", "/v1/companies/create", s.CreateCompany_0)
	s.router.Handle("POST", "/v1/companies/update/:id", s.UpdateCompany_0)
	s.router.Handle("POST", "/v1/companies/delete/:id", s.DeleteCompany_0)
	s.router.Handle("GET", "/v1/companies/:id", s.GetCompany_0)
	s.router.Handle("GET", "/v1/companies/address/:wallet_address", s.GetCompanyByAddress_0)
	s.router.Handle("GET", "/v1/companies", s.ListCompany_0)
}
