package converter

import (
	"database/sql"
	"time"
	"user_service/internal/model"
	desc "user_service/pkg/user_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)


func UserModelToProto(u *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if u.UpdatedAt.Valid {
		updatedAt = timestamppb.New(u.UpdatedAt.Time)
	}

	return &desc.User{
		Id: u.ID,
		Info: &desc.UserInfo{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			PhoneNumber:     u.Phone,
		},
		CreatedAt: timestamppb.New(u.CreatedAt), // безопасно, всегда есть
		UpdatedAt: updatedAt,                    // nullable
	}
}

func UserProtoToModel(u *desc.User) *model.User {
	if u == nil {
		return nil
	}

	var updatedAt sql.NullTime
	if u.UpdatedAt != nil {
		updatedAt = sql.NullTime{
			Time:  u.UpdatedAt.AsTime(),
			Valid: true,
		}
	}

	var createdAt time.Time
	if u.CreatedAt != nil {
		createdAt = u.CreatedAt.AsTime()
	}

	firstName, lastName, email, phone := "", "", "", ""
	if u.Info != nil {
		firstName = u.Info.FirstName
		lastName = u.Info.LastName
		email = u.Info.Email
		phone = u.Info.PhoneNumber
	}

	return &model.User{
		ID:        u.Id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		CreatedAt: createdAt,   // всегда есть
		UpdatedAt: updatedAt,   // nullable
	}
}