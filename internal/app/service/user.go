package service

import (
	"errors"
	"kazokku/internal/app/delivery/dto"
	"kazokku/internal/app/repository"
	"kazokku/internal/domain"
	"kazokku/internal/helpers"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserService interface {
	Create(ctx *fiber.Ctx, data dto.UserRegisterRequest) (uint, error)
	GetAll(ctx *fiber.Ctx, query dto.UserQuery) ([]dto.UserResponse, error)
	GetByID(ctx *fiber.Ctx, userID uint) (dto.UserResponse, error)
	UpdateByID(ctx *fiber.Ctx, data dto.UserUpdateRequest) error
}

type userService struct {
	db        *pgx.Conn
	userRepo  repository.UserRepository
	ccRepo    repository.CreditCardRepository
	photoRepo repository.PhotoRepository
}

func NewUserService(db *pgx.Conn, userRepo repository.UserRepository, ccRepo repository.CreditCardRepository, photoRepo repository.PhotoRepository) UserService {
	return userService{
		db:        db,
		userRepo:  userRepo,
		ccRepo:    ccRepo,
		photoRepo: photoRepo,
	}
}

func (s userService) Create(ctx *fiber.Ctx, data dto.UserRegisterRequest) (uint, error) {
	if err := data.Validate(); err != nil {
		return 0, helpers.NewResponseError(err, fiber.StatusBadRequest)
	}

	files, err := ctx.MultipartForm()
	if err != nil {
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	if len(files.File["photos"]) < 1 {
		return 0, helpers.NewResponseError(errors.New("Please provide photos fields."), fiber.StatusBadRequest)
	}

	tx, err := s.db.Begin(ctx.Context())
	if err != nil {
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	data.Password, err = helpers.HashPassword(data.Password)
	if err != nil {
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// create user record
	id, err := s.userRepo.Insert(ctx.Context(), tx, helpers.UserRegisterDTOtoUserDomain(data))
	if err != nil {
		tx.Rollback(ctx.Context())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, helpers.NewResponseError(errors.New("Email already registered."), fiber.StatusConflict)
		}
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// create credit card record
	err = s.ccRepo.Insert(ctx.Context(), tx, helpers.UserRegisterDTOtoCCDomain(data, id))
	if err != nil {
		tx.Rollback(ctx.Context())
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// save photos
	var photos []domain.Photo

	for _, file := range files.File["photos"] {
		if !helpers.IsImage(file.Header.Get("Content-Type")) {
			continue
		}

		savedFile, err := helpers.SaveFile(id, file)
		if err != nil {
			tx.Rollback(ctx.Context())
			return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
		}
		photos = append(photos, domain.Photo{
			UserID:   id,
			Filepath: savedFile,
		})
	}

	// create photo records
	err = s.photoRepo.InsertBatch(ctx.Context(), tx, photos)
	if err != nil {
		tx.Rollback(ctx.Context())
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// commit transaction
	err = tx.Commit(ctx.Context())
	if err != nil {
		return 0, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	return id, nil
}

func (s userService) GetAll(ctx *fiber.Ctx, query dto.UserQuery) ([]dto.UserResponse, error) {
	var users []dto.UserResponse
	if err := query.Validate(); err != nil {
		return users, helpers.NewResponseError(err, fiber.StatusBadRequest)
	}

	if query.OrderBy == "" {
		query.OrderBy = "name"
	}

	if query.SortBy == "" {
		query.SortBy = "asc"
	}

	if query.Offset < 0 {
		query.Offset = 0
	}

	if query.Limit <= 0 {
		query.Limit = 30
	}

	data, err := s.userRepo.GetAll(ctx.Context(), query)
	if err != nil {
		return users, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	for _, user := range data {
		var photos []string
		for _, photo := range user.Photos {
			photos = append(photos, filepath.Join("/photos", photo.Filepath))
		}

		users = append(users, dto.UserResponse{
			ID:      user.ID,
			Name:    user.Name.String,
			Email:   user.Email.String,
			Address: user.Address.String,
			Photos:  photos,
			CreditCard: dto.CreditCardResponse{
				Type:    user.CreditCard.Type.String,
				Number:  helpers.GetLast4Digits(user.CreditCard.Number.String),
				Name:    user.CreditCard.Name.String,
				Expired: user.CreditCard.Expired.String,
			},
		})
	}

	return users, nil
}

func (s userService) GetByID(ctx *fiber.Ctx, userID uint) (dto.UserResponse, error) {
	var user dto.UserResponse

	data, err := s.userRepo.GetByID(ctx.Context(), userID)
	if err != nil {
		return user, helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	if data.IsEmpty() {
		return user, helpers.NewResponseError(errors.New("User not found."), fiber.StatusNotFound)
	}

	user.ID = data.ID
	user.Name = data.Name.String
	user.Email = data.Email.String
	user.Address = data.Address.String
	user.Photos = make([]string, len(data.Photos))
	user.CreditCard = dto.CreditCardResponse{
		Type:    data.CreditCard.Type.String,
		Number:  helpers.GetLast4Digits(data.CreditCard.Number.String),
		Name:    data.CreditCard.Name.String,
		Expired: data.CreditCard.Expired.String,
	}

	for i, photo := range data.Photos {
		user.Photos[i] = filepath.Join("/photos", photo.Filepath)
	}

	return user, nil
}

func (s userService) UpdateByID(ctx *fiber.Ctx, data dto.UserUpdateRequest) error {
	if err := data.Validate(); err != nil {
		return helpers.NewResponseError(err, fiber.StatusBadRequest)
	}

	files, err := ctx.MultipartForm()
	if err != nil {
		return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	if data.Password != "" {
		data.Password, err = helpers.HashPassword(data.Password)
		if err != nil {
			return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
		}
	}

	tx, err := s.db.Begin(ctx.Context())
	if err != nil {
		return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// update user record
	err = s.userRepo.Update(ctx.Context(), tx, data.UserID, helpers.UserUpdateDTOtoUserDomain(data))
	if err != nil {
		tx.Rollback(ctx.Context())
		return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// update credit card record
	err = s.ccRepo.Update(ctx.Context(), tx, helpers.UserUpdateDTOtoCCDomain(data, data.UserID))
	if err != nil {
		tx.Rollback(ctx.Context())
		return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	// save photos
	if len(files.File["photos"]) > 0 {
		var photos []domain.Photo
		for _, file := range files.File["photos"] {
			if !helpers.IsImage(file.Header.Get("Content-Type")) {
				continue
			}

			savedFile, err := helpers.SaveFile(data.UserID, file)
			if err != nil {
				tx.Rollback(ctx.Context())
				return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
			}
			photos = append(photos, domain.Photo{
				UserID:   data.UserID,
				Filepath: savedFile,
			})
		}

		// create photo records
		err = s.photoRepo.InsertBatch(ctx.Context(), tx, photos)
		if err != nil {
			tx.Rollback(ctx.Context())
			return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
		}
	}

	// commit transaction
	err = tx.Commit(ctx.Context())
	if err != nil {
		return helpers.NewResponseError(helpers.ErrInternal, fiber.StatusInternalServerError)
	}

	return nil
}
