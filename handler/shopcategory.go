package handler

import (
	"context"
	"github.com/wsjcko/shopcategory/common"
	"github.com/wsjcko/shopcategory/domain/model"
	"github.com/wsjcko/shopcategory/domain/service"
	pb "github.com/wsjcko/shopcategory/protobuf/pb"
	log "go-micro.dev/v4/logger"
)

type ShopCategory struct {
	CategoryDataService service.ICategoryDataService
}

func (c *ShopCategory) Init(categoryDataService service.ICategoryDataService) {
	c.CategoryDataService = categoryDataService
}

// CreateCategory 提供创建分类的服务
func (c *ShopCategory) CreateCategory(ctx context.Context, request *pb.CategoryRequest, response *pb.CreateCategoryResponse) error {
	category := &model.Category{}
	//赋值
	err := common.SwapTo(request, category)
	if err != nil {
		return err
	}
	categoryId, err := c.CategoryDataService.AddCategory(category)
	if err != nil {
		return err
	}
	response.Message = "分类添加成功"
	response.CategoryId = categoryId
	return nil
}

// UpdateCategory 提供分类更新服务
func (c *ShopCategory) UpdateCategory(ctx context.Context, request *pb.CategoryRequest, response *pb.UpdateCategoryResponse) error {
	category := &model.Category{}
	err := common.SwapTo(request, category)
	if err != nil {
		return err
	}
	err = c.CategoryDataService.UpdateCategory(category)
	if err != nil {
		return err
	}
	response.Message = "分类更新成功"
	return nil
}

// DeleteCategory 提供分类删除服务
func (c *ShopCategory) DeleteCategory(ctx context.Context, request *pb.DeleteCategoryRequest, response *pb.DeleteCategoryResponse) error {
	err := c.CategoryDataService.DeleteCategory(request.CategoryId)
	if err != nil {
		return nil
	}
	response.Message = "删除成功"
	return nil
}

// FindCategoryByName 根据分类名称查找分类
func (c *ShopCategory) FindCategoryByName(ctx context.Context, request *pb.FindByNameRequest, response *pb.CategoryResponse) error {
	category, err := c.CategoryDataService.FindCategoryByName(request.CategoryName)
	if err != nil {
		return err
	}
	return common.SwapTo(category, response)
}

// FindCategoryByID 根据分类ID查找分类
func (c *ShopCategory) FindCategoryByID(ctx context.Context, request *pb.FindByIdRequest, response *pb.CategoryResponse) error {
	category, err := c.CategoryDataService.FindCategoryByID(request.CategoryId)
	if err != nil {
		return err
	}
	return common.SwapTo(category, response)
}

func (c *ShopCategory) FindCategoryByLevel(ctx context.Context, request *pb.FindByLevelRequest, response *pb.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByLevel(request.Level)
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func (c *ShopCategory) FindCategoryByParent(ctx context.Context, request *pb.FindByParentRequest, response *pb.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByParent(request.ParentId)
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func (c *ShopCategory) FindAllCategory(ctx context.Context, request *pb.FindAllRequest, response *pb.FindAllResponse) error {
	categorySlice, err := c.CategoryDataService.FindAllCategory()
	if err != nil {
		return err
	}
	categoryToResponse(categorySlice, response)
	return nil
}

func categoryToResponse(categorySlice []model.Category, response *pb.FindAllResponse) {
	for _, cg := range categorySlice {
		cr := &pb.CategoryResponse{}
		err := common.SwapTo(cg, cr)
		if err != nil {
			log.Error(err)
			break
		}
		response.Category = append(response.Category, cr)
	}
}
