package service

import (
	"blog-backend/model"
	"blog-backend/repository"

	"errors"
	"gorm.io/gorm"
)

type MomentService struct {
	repo *repository.MomentRepository
}

func NewMomentService() *MomentService {
	return &MomentService{
		repo: repository.NewMomentRepository(),
	}
}

// Create 创建说说
func (s *MomentService) Create(moment *model.Moment) error {
	if moment.Content == "" {
		return errors.New("说说内容不能为空")
	}
	return s.repo.Create(moment)
}

// Update 更新说说
func (s *MomentService) Update(id uint, content, images string) error {
	moment, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if content != "" {
		moment.Content = content
	}
	moment.Images = images

	return s.repo.Update(moment)
}

// Delete 删除说说
func (s *MomentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetByID 获取说说详情
func (s *MomentService) GetByID(id uint) (*model.Moment, error) {
	return s.repo.GetByID(id)
}

// List 获取说说列表
func (s *MomentService) List(page, pageSize int, status *int, keyword string, userID *uint, ip string) ([]model.Moment, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	moments, total, err := s.repo.List(page, pageSize, status, keyword)
	if err != nil {
		return nil, 0, err
	}

	// 获取当前用户点赞过的说说ID列表
	if len(moments) > 0 {
		momentIDs := make([]uint, len(moments))
		for i, m := range moments {
			momentIDs[i] = m.ID
		}

		likedIDs, err := s.repo.GetLikedMomentIDs(momentIDs, userID, ip)
		if err == nil {
			// 标记已点赞的说说
			likedMap := make(map[uint]bool)
			for _, id := range likedIDs {
				likedMap[id] = true
			}

			for i := range moments {
				moments[i].Liked = likedMap[moments[i].ID]
			}
		}
	}

	return moments, total, nil
}

// GetRecent 获取最新说说
func (s *MomentService) GetRecent(limit int) ([]model.Moment, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return s.repo.GetRecent(limit)
}

// Like 点赞/取消点赞说说
func (s *MomentService) Like(id uint, userID *uint, ip string) (bool, error) {
	// 检查是否已点赞
	liked, err := s.repo.CheckLiked(id, userID, ip)
	if err != nil {
		return false, err
	}

	// 获取说说
	moment, err := s.repo.GetByID(id)
	if err != nil {
		return false, err
	}

	// 使用事务确保数据一致性
	var isLiked bool
	err = s.repo.Transaction(func(tx *gorm.DB) error {
		if liked {
			// 已点赞，执行取消点赞
			if err := s.repo.DeleteLikeTx(tx, id, userID, ip); err != nil {
				return err
			}
			
			// 减少点赞数
			if moment.LikeCount > 0 {
				moment.LikeCount--
			}
			if err := s.repo.UpdateTx(tx, moment); err != nil {
				return err
			}
			
			isLiked = false
		} else {
			// 未点赞，执行点赞
			like := &model.MomentLike{
				MomentID: id,
				UserID:   userID,
				IP:       ip,
			}
			if err := s.repo.CreateLikeTx(tx, like); err != nil {
				return err
			}

			// 增加点赞数
			moment.LikeCount++
			if err := s.repo.UpdateTx(tx, moment); err != nil {
				return err
			}
			
			isLiked = true
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return isLiked, nil
}
