package service

import (
	"career_backend/internal/dal/db/model"
	"career_backend/internal/ecode"
	"career_backend/internal/repository"
	"context"
	"github.com/go-eagle/eagle/pkg/errcode"

	pb "career_backend/api/company/v1"
)

var (
	_ pb.CompanyServiceServer = (*CompanyServiceServer)(nil)
)

type CompanyServiceServer struct {
	pb.UnimplementedCompanyServiceServer

	repo repository.CompanyRepo
}

func NewCompanyServiceServer() *CompanyServiceServer {
	return &CompanyServiceServer{}
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
		return nil, ecode.ErrInternalError.WithDetails(errcode.NewDetails(map[string]interface{}{
			"msg": err.Error(),
		})).Status(req).Err()
	}
	if company != nil && company.ID > 0 {
		return nil, ecode.ErrUserIsExist.Status(req).Err()
	}

	newCompany := &model.Company{
		Name:          req.Name,
		WalletAddress: req.WalletAddress,
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
	return &pb.GetCompanyReply{}, nil
}
func (s *CompanyServiceServer) GetCompanyByAddress(ctx context.Context, req *pb.GetCompanyByAddressRequest) (*pb.GetCompanyReply, error) {
	return &pb.GetCompanyReply{}, nil
}
func (s *CompanyServiceServer) ListCompany(ctx context.Context, req *pb.ListCompanyRequest) (*pb.ListCompanyReply, error) {
	return &pb.ListCompanyReply{}, nil
}
