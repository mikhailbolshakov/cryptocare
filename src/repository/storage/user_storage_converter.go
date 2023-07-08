package storage

import (
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
	"github.com/mikhailbolshakov/cryptocare/src/kit/storages/pg"
)

func (s *userStorageImpl) toUserDto(u *auth.User) *user {
	if u == nil {
		return nil
	}
	dto := &user{
		Id:          u.Id,
		Username:    u.Username,
		Password:    pg.StringToNull(u.Password),
		Type:        u.Type,
		ActivatedAt: u.ActivatedAt,
		LockedAt:    u.LockedAt,
	}
	det := &userDetails{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Groups:    u.Groups,
		Roles:     u.Roles,
	}
	var detailsBytes []byte
	detailsBytes, _ = json.Marshal(det)
	dto.Details = string(detailsBytes)
	return dto
}

func (s *userStorageImpl) toUserCacheDomain(rec *aero.Record) *auth.User {
	if rec == nil {
		return nil
	}
	body := rec.Bins["user"].(string)
	user := &auth.User{}
	_ = json.Unmarshal([]byte(body), user)
	return user

}

func (s *userStorageImpl) toUserCache(user *auth.User) aero.BinMap {
	usrBytes, _ := json.Marshal(user)
	return aero.BinMap{
		"username": user.Username,
		"user":     string(usrBytes),
	}
}

func (s *userStorageImpl) toUserDomain(dto *user) *auth.User {
	if dto == nil {
		return nil
	}
	det := &userDetails{}
	_ = json.Unmarshal([]byte(dto.Details), det)
	return &auth.User{
		Id:          dto.Id,
		Username:    dto.Username,
		Password:    pg.NullToString(dto.Password),
		Type:        dto.Type,
		FirstName:   det.FirstName,
		LastName:    det.LastName,
		ActivatedAt: dto.ActivatedAt,
		LockedAt:    dto.LockedAt,
		Groups:      det.Groups,
		Roles:       det.Roles,
	}
}

func (s *userStorageImpl) toUsersDomain(dtos []*user) []*auth.User {
	var res []*auth.User
	for _, d := range dtos {
		res = append(res, s.toUserDomain(d))
	}
	return res
}
