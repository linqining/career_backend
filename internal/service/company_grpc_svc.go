package service

import (
	"career_backend/internal/dal/db/model"
	"career_backend/internal/ecode"
	"career_backend/internal/repository"
	"context"
	"errors"
	"github.com/go-eagle/eagle/pkg/errcode"
	"gorm.io/gorm"

	pb "career_backend/api/company/v1"
)

var (
	_ pb.CompanyServiceServer = (*CompanyServiceServer)(nil)
)

type CompanyServiceServer struct {
	pb.UnimplementedCompanyServiceServer

	repo repository.CompanyRepo
}

func NewCompanyServiceServer(repo repository.CompanyRepo) *CompanyServiceServer {
	return &CompanyServiceServer{repo: repo}
}

func (s *CompanyServiceServer) CreateCompany(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CreateCompanyReply, error) {
	err := req.Validate()
	if err != nil {

	}
	if req.WalletAddress == "" {
		return nil, ecode.ErrInvalidArgument.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": "invalid wallet address",
		})).Status(req).Err()
	}

	var company *model.Company

	company, err = s.repo.GetCompanyByAddress(ctx, req.WalletAddress)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
				"msg": err.Error(),
			})).Status(req).Err()
		}
	}
	if company != nil {
		return nil, ecode.ErrCompanyIsExist.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": "company address exists",
		})).Status(req).Err()
	}

	newCompany := &model.Company{
		Name:          req.Name,
		WalletAddress: req.WalletAddress,
		LogoURL:       req.LogoUrl,
	}

	companyID, err := s.repo.CreateCompany(ctx, newCompany)
	if err != nil {
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}

	return &pb.CreateCompanyReply{
		Id:            companyID,
		Name:          newCompany.Name,
		WalletAddress: newCompany.WalletAddress,
	}, nil
}

func (s *CompanyServiceServer) UpdateCompany(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.UpdateCompanyReply, error) {
	return &pb.UpdateCompanyReply{}, nil
}
func (s *CompanyServiceServer) DeleteCompany(ctx context.Context, req *pb.DeleteCompanyRequest) (*pb.DeleteCompanyReply, error) {
	return &pb.DeleteCompanyReply{}, nil
}
func (s *CompanyServiceServer) GetCompany(ctx context.Context, req *pb.GetCompanyRequest) (*pb.GetCompanyReply, error) {
	company, err := s.repo.GetCompany(ctx, int64(req.Id))
	if err != nil {
		return nil, ecode.ErrCompanyNotExist.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	return &pb.GetCompanyReply{
		Company: &pb.Company{
			Id:            company.ID,
			Name:          company.Name,
			WalletAddress: company.WalletAddress,
			LogoUrl:       company.LogoURL,
		},
	}, nil
}

func (s *CompanyServiceServer) GetCompanyByAddress(ctx context.Context, req *pb.GetCompanyByAddressRequest) (*pb.GetCompanyReply, error) {
	company, err := s.repo.GetCompanyByAddress(ctx, req.WalletAddress)
	if err != nil {
		return nil, ecode.ErrCompanyNotExist.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	return &pb.GetCompanyReply{
		Company: &pb.Company{
			Id:            company.ID,
			Name:          company.Name,
			WalletAddress: company.WalletAddress,
			LogoUrl:       company.LogoURL,
		},
	}, nil
}
func (s *CompanyServiceServer) ListCompany(ctx context.Context, req *pb.ListCompanyRequest) (*pb.ListCompanyReply, error) {
	return &pb.ListCompanyReply{}, nil
}
