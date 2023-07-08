package storage

import (
	"encoding/json"
	aero "github.com/aerospike/aerospike-client-go/v6"
	"github.com/mikhailbolshakov/cryptocare/src/kit/auth"
)

func (s *sessionStorageImpl) toSessionCacheDomain(rec *aero.Record) *auth.Session {
	if rec == nil {
		return nil
	}
	body := rec.Bins["session"].(string)
	sess := &auth.Session{}
	_ = json.Unmarshal([]byte(body), sess)
	return sess
}

func (s *sessionStorageImpl) toSessionCache(sess *auth.Session) aero.BinMap {
	sessBytes, _ := json.Marshal(sess)
	return aero.BinMap{
		"session": string(sessBytes),
		"user_id": sess.UserId,
	}
}

func (s *sessionStorageImpl) toSessionDomain(dto *session) *auth.Session {
	if dto == nil {
		return nil
	}
	det := &sessionDetails{}
	_ = json.Unmarshal([]byte(dto.Details), det)
	return &auth.Session{
		Id:             dto.Id,
		UserId:         dto.UserId,
		Username:       dto.Username,
		LoginAt:        dto.LoginAt,
		LogoutAt:       dto.LogoutAt,
		LastActivityAt: dto.LastActivityAt,
		Roles:          det.Roles,
	}
}

func (s *sessionStorageImpl) toSessionsDomain(dtos []*session) []*auth.Session {
	var res []*auth.Session
	for _, d := range dtos {
		res = append(res, s.toSessionDomain(d))
	}
	return res
}

func (s *sessionStorageImpl) toSessionDto(sess *auth.Session) *session {
	if sess == nil {
		return nil
	}
	dto := &session{
		Id:             sess.Id,
		UserId:         sess.UserId,
		Username:       sess.Username,
		LoginAt:        sess.LoginAt,
		LastActivityAt: sess.LastActivityAt,
		LogoutAt:       sess.LogoutAt,
	}
	det := &sessionDetails{
		Roles: sess.Roles,
	}
	var detailsBytes []byte
	detailsBytes, _ = json.Marshal(det)
	dto.Details = string(detailsBytes)
	return dto
}
