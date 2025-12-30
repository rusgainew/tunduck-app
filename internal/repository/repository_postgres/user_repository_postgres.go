package repositorypostgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
)

type UserRepositoryPostgres struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserRepositoryPostgres(db *gorm.DB, log *logrus.Logger) repository.UserRepository {
	return &UserRepositoryPostgres{
		db:     db,
		logger: logger.New(log),
	}
}

func (r *UserRepositoryPostgres) Create(ctx context.Context, user *entity.User) error {
	r.logger.Debug(ctx, "Creating user in database", logrus.Fields{"username": user.Username, "email": user.Email})

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error(ctx, "Failed to create user in database", err, logrus.Fields{"username": user.Username})
		return apperror.DatabaseError("creating user", err)
	}

	r.logger.Debug(ctx, "User created successfully", logrus.Fields{"id": user.ID.String()})
	return nil
}

func (r *UserRepositoryPostgres) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	r.logger.Debug(ctx, "Fetching user by ID", logrus.Fields{"user_id": id.String()})

	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug(ctx, "User not found", logrus.Fields{"user_id": id.String()})
			return nil, nil
		}
		r.logger.Error(ctx, "Failed to fetch user by ID", err, logrus.Fields{"user_id": id.String()})
		return nil, apperror.DatabaseError("fetching user by ID", err)
	}

	r.logger.Debug(ctx, "User fetched successfully", logrus.Fields{"user_id": id.String()})
	return &user, nil
}

func (r *UserRepositoryPostgres) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.logger.Debug(ctx, "Fetching user by username", logrus.Fields{"username": username})

	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug(ctx, "User not found by username", logrus.Fields{"username": username})
			return nil, nil
		}
		r.logger.Error(ctx, "Failed to fetch user by username", err, logrus.Fields{"username": username})
		return nil, apperror.DatabaseError("fetching user by username", err)
	}

	r.logger.Debug(ctx, "User fetched successfully", logrus.Fields{"user_id": user.ID.String()})
	return &user, nil
}

func (r *UserRepositoryPostgres) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.logger.Debug(ctx, "Fetching user by email", logrus.Fields{"email": email})

	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug(ctx, "User not found by email", logrus.Fields{"email": email})
			return nil, nil
		}
		r.logger.Error(ctx, "Failed to fetch user by email", err, logrus.Fields{"email": email})
		return nil, apperror.DatabaseError("fetching user by email", err)
	}

	r.logger.Debug(ctx, "User fetched successfully", logrus.Fields{"user_id": user.ID.String()})
	return &user, nil
}

func (r *UserRepositoryPostgres) GetAll(ctx context.Context, limit int) ([]*entity.User, error) {
	r.logger.Debug(ctx, "Fetching all users", logrus.Fields{"limit": limit})

	var users []*entity.User
	query := r.db.WithContext(ctx)

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&users).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to fetch users", err)
		return nil, apperror.DatabaseError("fetching users", err)
	}

	r.logger.Debug(ctx, "Users fetched successfully", logrus.Fields{"count": len(users)})
	return users, nil
}

func (r *UserRepositoryPostgres) Update(ctx context.Context, user *entity.User) error {
	r.logger.Debug(ctx, "Updating user in database", logrus.Fields{"user_id": user.ID.String()})

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error(ctx, "Failed to update user in database", err, logrus.Fields{"user_id": user.ID.String()})
		return apperror.DatabaseError("updating user", err)
	}

	r.logger.Debug(ctx, "User updated successfully", logrus.Fields{"user_id": user.ID.String()})
	return nil
}

func (r *UserRepositoryPostgres) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Debug(ctx, "Deleting user from database", logrus.Fields{"user_id": id.String()})

	if err := r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error; err != nil {
		r.logger.Error(ctx, "Failed to delete user from database", err, logrus.Fields{"user_id": id.String()})
		return apperror.DatabaseError("deleting user", err)
	}

	r.logger.Debug(ctx, "User deleted successfully", logrus.Fields{"user_id": id.String()})
	return nil
}
