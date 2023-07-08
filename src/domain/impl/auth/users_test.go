package auth

//
//import (
//	"fmt"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/suite"
//	"gitlab.moi-service.ru/moi-service/platform/auth/domain"
//	"gitlab.moi-service.ru/moi-service/platform/auth/errors"
//	"gitlab.moi-service.ru/moi-service/platform/auth/logger"
//	"gitlab.moi-service.ru/moi-service/platform/auth/mocks"
//	authPb "gitlab.moi-service.ru/moi-service/platform/auth/proto"
//	"gitlab.moi-service.ru/moi-service/platform/chat/proto"
//	"gitlab.moi-service.ru/moi-service/platform/kit"
//	"gitlab.moi-service.ru/moi-service/platform/kit/common"
//	kitMocks "gitlab.moi-service.ru/moi-service/platform/kit/mocks"
//	kitTestSuite "gitlab.moi-service.ru/moi-service/platform/kit/test/suite"
//	"testing"
//)
//
//type usersTestSuite struct {
//	kitTestSuite.Suite
//	queue             *kitMocks.Queue
//	userStorage       *mocks.UserStorage
//	chatService       *mocks.ChatRepository
//	passwordGenerator *mocks.PasswordGenerator
//	emailAdapter      *mocks.EmailRepository
//	securityService   *mocks.SecurityService
//	userSvc           domain.UserService
//}
//
//func (s *usersTestSuite) SetupSuite() {
//	s.Suite.Init(logger.LF())
//	s.queue = &kitMocks.Queue{}
//	s.userStorage = &mocks.UserStorage{}
//	s.chatService = &mocks.ChatRepository{}
//	s.passwordGenerator = &mocks.PasswordGenerator{}
//	s.emailAdapter = &mocks.EmailRepository{}
//	s.securityService = &mocks.SecurityService{}
//
//	s.queue.On("Publish", mock.Anything, mock.Anything, authPb.QueueTopicUserDraftCreated, mock.Anything).Return(nil)
//	s.queue.On("Publish", mock.Anything, mock.Anything, authPb.QueueTopicUserLocked, mock.Anything).Return(nil)
//	s.queue.On("Publish", mock.Anything, mock.Anything, authPb.QueueTopicUserDeleted, mock.Anything).Return(nil)
//	s.chatService.On("CreateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*chat.CreateUserRequest")).Return(&chat.CreateUserResponse{
//		ChatUserId:   "1",
//		ChatUsername: "username",
//	}, nil)
//
//	s.userSvc = NewUserService(s.userStorage, s.chatService, s.securityService, s.queue, s.passwordGenerator, s.emailAdapter)
//}
//
//func (s *usersTestSuite) SetupTest() {
//	s.userStorage.ExpectedCalls = nil
//	s.securityService.ExpectedCalls = nil
//	s.emailAdapter.Calls = nil
//	s.emailAdapter.ExpectedCalls = nil
//	s.userStorage.On("CreateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User")).Return(nil)
//	s.userStorage.On("UpdateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User")).Return(nil)
//}
//
//func TestUsersSuite(t *testing.T) {
//	suite.Run(t, new(usersTestSuite))
//}
//
//func (s *usersTestSuite) clientDefault() (*domain.CreateUserRequest, *domain.User) {
//	rq := &domain.CreateUserRequest{
//		Type:       authPb.UserTypeClient,
//		FirstName:  "test",
//		MiddleName: "test",
//		LastName:   "test",
//		Email:      "user@test.com",
//		Phone:      "7903222333444",
//		IsGuest:    false,
//		Groups:     []string{authPb.UserGrpClient},
//	}
//
//	exp := &domain.User{
//		Id:       kit.NewId(),
//		Username: rq.Phone,
//		Type:     rq.Type,
//		AuthType: domain.AUTH_TYPE_SMS_CODE,
//		Status:   authPb.UserStatusDraft,
//		Details: &domain.UserDetails{
//			FirstName:    rq.FirstName,
//			MiddleName:   rq.MiddleName,
//			LastName:     rq.LastName,
//			Email:        rq.Email,
//			Phone:        rq.Phone,
//			Avatar:       rq.Avatar,
//			ChatUserId:   "1",
//			ChatUsername: "username",
//			IsGuest:      rq.IsGuest,
//			Groups:       rq.Groups,
//		},
//	}
//
//	return rq, exp
//
//}
//
//func (s *usersTestSuite) consultantDefault() (*domain.CreateUserRequest, *domain.User) {
//
//	rq := &domain.CreateUserRequest{
//		Type:       authPb.UserTypeConsultant,
//		Password:   "123",
//		FirstName:  "test",
//		MiddleName: "test",
//		LastName:   "test",
//		Email:      "user@test.com",
//		Phone:      "7903222333444",
//		IsGuest:    false,
//		Groups:     []string{authPb.UserGrpConsultantProperty},
//	}
//
//	expected := &domain.User{
//		Id:       kit.NewId(),
//		Username: rq.Email,
//		Type:     rq.Type,
//		AuthType: domain.AUTH_TYPE_PASSWORD,
//		Status:   authPb.UserStatusDraft,
//		Details: &domain.UserDetails{
//			FirstName:    rq.FirstName,
//			MiddleName:   rq.MiddleName,
//			LastName:     rq.LastName,
//			Email:        rq.Email,
//			Phone:        rq.Phone,
//			Avatar:       rq.Avatar,
//			ChatUserId:   "1",
//			ChatUsername: "username",
//			IsGuest:      rq.IsGuest,
//			Groups:       rq.Groups,
//		},
//	}
//	return rq, expected
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenRequestValid() {
//	rq, expected := s.clientDefault()
//
//	group := &domain.Group{
//		Code:    authPb.UserGrpClient,
//		Default: true,
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return([]*domain.Group{group}, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	got, err := s.userSvc.Create(s.Ctx, rq)
//	if err != nil {
//		s.Fatal(err)
//	}
//
//	expected.Id = got.Id
//	expected.CreatedAt = got.CreatedAt
//	expected.UpdatedAt = got.UpdatedAt
//
//	s.Equal(expected, got)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenInvalidType() {
//	rq, _ := s.clientDefault()
//	rq.Type = "wrong_type"
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenInvalidPhone() {
//	rq, _ := s.clientDefault()
//	rq.Phone = "abc"
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenEmptyPhone() {
//	rq, _ := s.clientDefault()
//	rq.Phone = ""
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenInvalidEmail() {
//	rq, _ := s.clientDefault()
//	rq.Email = "invalid-email"
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenWrongGroups() {
//	rq, _ := s.clientDefault()
//	rq.Groups = []string{authPb.UserGrpIntegration, authPb.UserGrpSupervisor}
//
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserTypeClient,
//			Default: true,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return(groups, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.AssertAppErr(err, errors.ErrCodeUserGroupNotAllowed)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Client_WhenUserExists() {
//	rq, _ := s.clientDefault()
//	group := &domain.Group{
//		Code:    authPb.UserGrpClient,
//		Default: true,
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return([]*domain.Group{group}, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, &domain.User{}, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Consultant_WhenRequestValid() {
//	rq, expected := s.consultantDefault()
//
//	group := &domain.Group{
//		Code:    authPb.UserGrpConsultantProperty,
//		Default: true,
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return([]*domain.Group{group}, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	got, err := s.userSvc.Create(s.Ctx, rq)
//	if err != nil {
//		s.Fatal(err)
//	}
//	s.NotEmpty(got.Password)
//	expected.Id = got.Id
//	expected.CreatedAt = got.CreatedAt
//	expected.UpdatedAt = got.UpdatedAt
//	expected.Password = got.Password
//	s.Equal(expected, got)
//	s.AssertNotCalled(&s.emailAdapter.Mock, "SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest"))
//}
//
//func (s *usersTestSuite) Test_UserCreate_WhenIdAndStatusPassedFromOutside_Ok() {
//	rq, expected := s.consultantDefault()
//	rq.Id = kit.NewId()
//	rq.Status = authPb.UserStatusActive
//	expected.Id = rq.Id
//	expected.Status = rq.Status
//
//	group := &domain.Group{
//		Code:    authPb.UserGrpConsultantProperty,
//		Default: true,
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return([]*domain.Group{group}, nil)
//
//	s.emailAdapter.On("SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest")).Return(nil)
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	got, err := s.userSvc.Create(s.Ctx, rq)
//	if err != nil {
//		s.Fatal(err)
//	}
//	s.NotEmpty(got.Password)
//	expected.Id = got.Id
//	expected.CreatedAt = got.CreatedAt
//	expected.UpdatedAt = got.UpdatedAt
//	expected.Password = got.Password
//	s.Equal(expected, got)
//	s.AssertCalled(&s.emailAdapter.Mock, "SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest"))
//}
//
//func (s *usersTestSuite) Test_UserCreate_Consultant_WhenDefaultGroups_Ok() {
//	rq, _ := s.consultantDefault()
//	rq.Groups = nil
//
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: true,
//		},
//		{
//			Code:    authPb.UserGrpSupervisor,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return(groups, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	user, err := s.userSvc.Create(s.Ctx, rq)
//	s.Nil(err)
//	s.Equal(1, len(user.Details.Groups))
//	s.Equal(authPb.UserGrpConsultantProperty, user.Details.Groups[0])
//}
//
//func (s *usersTestSuite) Test_UserCreate_Consultant_WhenEmptyGroups_Fail() {
//	rq, _ := s.consultantDefault()
//	rq.Groups = nil
//
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), rq.Type).Return(groups, nil)
//
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.AssertAppErr(err, errors.ErrCodeUserNoGroup)
//}
//
//func (s *usersTestSuite) Test_UserCreate_Consultant_WhenGuest() {
//	rq, _ := s.consultantDefault()
//	rq.IsGuest = true
//	s.userStorage.On("GetByUsername", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.Create(s.Ctx, rq)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_SetStatus_WhenActiveFromDraft() {
//	userId := "1"
//	status := authPb.UserStatusActive
//	_, expected := s.clientDefault()
//	expected.Id = userId
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, expected, nil)
//	s.userStorage.On("UpdateStatus", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User")).Return(nil)
//	got, err := s.userSvc.SetStatus(s.Ctx, userId, status)
//	if err != nil {
//		s.Fatal(err)
//	}
//	s.Equal(status, got.Status)
//	s.Empty(got.DeletedAt)
//}
//
//func (s *usersTestSuite) Test_SetStatus_WhenSameStatus() {
//	userId := "1"
//	status := authPb.UserStatusDraft
//	_, expected := s.clientDefault()
//	expected.Id = userId
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, expected, nil)
//
//	userStorageMock := s.userStorage.On("UpdateStatus", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User"))
//	userStorageMock.RunFn = func(args mock.Arguments) {
//		userStorageMock.ReturnArguments = mock.Arguments{args.Get(1).(*domain.User), nil}
//	}
//
//	got, err := s.userSvc.SetStatus(s.Ctx, userId, status)
//	if err != nil {
//		s.Fatal(err)
//	}
//
//	s.Equal(status, got.Status)
//	s.Empty(got.DeletedAt)
//}
//
//func (s *usersTestSuite) Test_Delete_Ok() {
//	userId := "1"
//	_, expected := s.clientDefault()
//	expected.Id = userId
//	var gotUserId string
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, expected, nil)
//	s.userStorage.On("Delete", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User")).
//		Run(func(args mock.Arguments) {
//			u := args.Get(1).(*domain.User)
//			gotUserId = u.Id
//			now := kit.Now()
//			expected.DeletedAt = &now
//		}).
//		Return(nil)
//	err := s.userSvc.Delete(s.Ctx, userId)
//	if err != nil {
//		s.Fatal(err)
//	}
//	s.Equal(expected.Id, gotUserId)
//	s.NotEmpty(expected.DeletedAt)
//}
//
//func (s *usersTestSuite) Test_SetStatus_WhenActivationWithEmptyChatUser() {
//	userId := "1"
//	status := authPb.UserStatusActive
//	_, expected := s.clientDefault()
//	expected.Id = userId
//	expected.Details.ChatUserId = ""
//
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, expected, nil)
//
//	userStorageMock := s.userStorage.On("UpdateStatus", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.User"))
//	userStorageMock.RunFn = func(args mock.Arguments) {
//		userStorageMock.ReturnArguments = mock.Arguments{args.Get(1).(*domain.User), nil}
//	}
//
//	_, err := s.userSvc.SetStatus(s.Ctx, userId, status)
//	fmt.Println(err)
//	s.Error(err)
//}
//
//func (s *usersTestSuite) Test_Search_WithoutOnlineStatusesCriteria() {
//	searchResp := &domain.UserSearchResponse{
//		PagingResponse: &common.PagingResponse{
//			Total: 2,
//			Index: 0,
//		},
//		Users: []*domain.User{
//			{
//				Id:      "1",
//				Details: &domain.UserDetails{ChatUserId: "1"},
//			},
//			{
//				Id:      "2",
//				Details: &domain.UserDetails{ChatUserId: "2"},
//			},
//		},
//	}
//	s.userStorage.On("Search", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*domain.UserSearchCriteria")).
//		Return(searchResp, nil)
//
//	rs, err := s.userSvc.Search(s.Ctx, &domain.UserSearchCriteria{})
//	if err != nil {
//		s.Fatal(err)
//	}
//
//	s.NotEmpty(rs.Users)
//	s.Equal(2, len(rs.Users))
//	s.Equal(2, rs.Total)
//}
//
//func (s *usersTestSuite) Test_UpdateDetails_WhenUserNotFound() {
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(false, nil, nil)
//	_, err := s.userSvc.UpdateUserDetails(s.Ctx, "123", &domain.UserDetails{})
//	s.AssertAppErr(err, errors.ErrCodeUserNotFound)
//}
//
//func (s *usersTestSuite) Test_UpdateDetails_WhenUserPhoneLoginModified() {
//	_, user := s.clientDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	_, err := s.userSvc.UpdateUserDetails(s.Ctx, user.Id, &domain.UserDetails{
//		Phone: "123",
//	})
//	s.AssertAppErr(err, errors.ErrCodeUserPhoneLoginModification)
//}
//
//func (s *usersTestSuite) Test_UpdateDetails_WhenUserEmailLoginModified() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	_, err := s.userSvc.UpdateUserDetails(s.Ctx, user.Id, &domain.UserDetails{
//		Email: "123",
//	})
//	s.AssertAppErr(err, errors.ErrCodeUserEmailLoginModification)
//}
//
//func (s *usersTestSuite) Test_UpdateDetails_WhenModifiedLastName_Ok() {
//	_, user := s.clientDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	user, err := s.userSvc.UpdateUserDetails(s.Ctx, "123", &domain.UserDetails{
//		FirstName:  user.Details.FirstName,
//		MiddleName: user.Details.MiddleName,
//		LastName:   "modified",
//		Avatar:     user.Details.Avatar,
//		Email:      user.Details.Email,
//		Phone:      user.Details.Phone,
//	})
//	if err != nil {
//		s.Fatal(err)
//	}
//	s.NotEmpty(user)
//	s.Equal("modified", user.Details.LastName)
//}
//
//func (s *usersTestSuite) Test_ResetPassword_Ok() {
//	s.passwordGenerator.On("Generate", mock.AnythingOfType("*context.valueCtx")).Return("passssssword", nil)
//	s.emailAdapter.On("SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest")).Return(nil)
//	_, user := s.clientDefault()
//	user.Status = authPb.UserStatusActive
//	user.Type = authPb.UserTypeAdmin
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Id).Return(true, user, nil)
//	s.userStorage.On("UpdateUser", mock.AnythingOfType("*context.valueCtx"), user.Id).Return(nil)
//
//	err := s.userSvc.ResetPassword(s.Ctx, user.Id, "")
//	s.Nil(err)
//	s.AssertCalled(&s.emailAdapter.Mock, "SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest"))
//}
//
//func (s *usersTestSuite) Test_ResetPassword_UserTypeNotValid() {
//
//	s.passwordGenerator.On("Generate", mock.AnythingOfType("*context.valueCtx")).Return("passssssword", nil)
//	_, user := s.clientDefault()
//	user.Status = authPb.UserStatusActive
//	user.Type = authPb.UserTypeClient
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Id).Return(true, user, nil)
//	s.userStorage.On("UpdateUser", mock.AnythingOfType("*context.valueCtx"), user.Id).Return(nil)
//
//	err := s.userSvc.ResetPassword(s.Ctx, user.Id, "")
//	s.AssertAppErr(err, errors.ErrCodeUserNoValidType)
//	s.AssertNotCalled(&s.emailAdapter.Mock, "SendAsync", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*communications.EmailRequest"))
//}
//
//func (s *usersTestSuite) Test_ResetPassword_NotActiveStatusOfUser() {
//
//	_, user := s.clientDefault()
//	user.Type = authPb.UserTypeAdmin
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), user.Id).Return(true, user, nil)
//
//	err := s.userSvc.ResetPassword(s.Ctx, user.Id, "")
//	s.AssertAppErr(err, errors.ErrCodeUserNotActive)
//}
//
//func (s *usersTestSuite) Test_AddGroups_WhenNoGroupConfigured_Fail() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(nil, nil)
//	_, err := s.userSvc.AddGroups(s.Ctx, user.Id, []string{"group"})
//	s.AssertAppErr(err, errors.ErrCodeUserNoGroupSpecified)
//}
//
//func (s *usersTestSuite) Test_AddGroups_WhenGroupNotAllowed_Fail() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(groups, nil)
//	_, err := s.userSvc.AddGroups(s.Ctx, user.Id, []string{"group"})
//	s.AssertAppErr(err, errors.ErrCodeUserGroupNotAllowed)
//}
//
//func (s *usersTestSuite) Test_AddGroups_WhenDuplicatedWithExistent_Ok() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(groups, nil)
//	user, err := s.userSvc.AddGroups(s.Ctx, user.Id, []string{authPb.UserGrpConsultantProperty, authPb.UserGrpConsultantProperty})
//	s.Nil(err, err)
//	s.Equal(1, len(user.Details.Groups))
//}
//
//func (s *usersTestSuite) Test_AddGroups_WhenAdded_Ok() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//		{
//			Code:    authPb.UserGrpSupervisor,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(groups, nil)
//	user, err := s.userSvc.AddGroups(s.Ctx, user.Id, []string{authPb.UserGrpSupervisor})
//	s.Nil(err, err)
//	s.Equal(2, len(user.Details.Groups))
//}
//
//func (s *usersTestSuite) Test_DeleteGroups_WhenExist_Ok() {
//	_, user := s.consultantDefault()
//	user.Details.Groups = []string{authPb.UserGrpConsultantProperty, authPb.UserGrpSupervisor}
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//		{
//			Code:    authPb.UserGrpSupervisor,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(groups, nil)
//	user, err := s.userSvc.DeleteGroups(s.Ctx, user.Id, []string{authPb.UserGrpConsultantProperty})
//	s.Nil(err, err)
//	s.Equal(1, len(user.Details.Groups))
//}
//
//func (s *usersTestSuite) Test_DeleteGroups_WhenNoGroups_Fail() {
//	_, user := s.consultantDefault()
//	user.Details.Groups = []string{authPb.UserGrpConsultantProperty, authPb.UserGrpSupervisor}
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	groups := []*domain.Group{
//		{
//			Code:    authPb.UserGrpConsultantProperty,
//			Default: false,
//		},
//		{
//			Code:    authPb.UserGrpSupervisor,
//			Default: false,
//		},
//	}
//	s.securityService.On("GetGroupsByUserType", mock.AnythingOfType("*context.valueCtx"), user.Type).Return(groups, nil)
//	_, err := s.userSvc.DeleteGroups(s.Ctx, user.Id, []string{authPb.UserGrpConsultantProperty, authPb.UserGrpSupervisor})
//	s.AssertAppErr(err, errors.ErrCodeUserNoGroup)
//}
//
//func (s *usersTestSuite) Test_GrantRoles_WhenInvalidRole_Fail() {
//	_, user := s.consultantDefault()
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	roles := []*domain.Role{
//		{
//			Code: "role1",
//		},
//	}
//	s.securityService.On("GetAllRoles", mock.AnythingOfType("*context.valueCtx")).Return(roles, nil)
//	_, err := s.userSvc.GrantRoles(s.Ctx, user.Id, []string{"role2"})
//	s.AssertAppErr(err, errors.ErrCodeUserInvalidRole)
//}
//
//func (s *usersTestSuite) Test_GrantRoles_WhenDuplicated_Fail() {
//	_, user := s.consultantDefault()
//	user.Details.Roles = []string{"role1", "role2"}
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	roles := []*domain.Role{
//		{
//			Code: "role1",
//		},
//		{
//			Code: "role2",
//		},
//		{
//			Code: "role3",
//		},
//	}
//	s.securityService.On("GetAllRoles", mock.AnythingOfType("*context.valueCtx")).Return(roles, nil)
//	user, err := s.userSvc.GrantRoles(s.Ctx, user.Id, []string{"role2", "role2", "role3", "role3"})
//	s.Nil(err, err)
//	s.Equal(3, len(user.Details.Roles))
//}
//
//func (s *usersTestSuite) Test_RevokeRoles_WhenNotExists_Ok() {
//	_, user := s.consultantDefault()
//	user.Details.Roles = []string{"role1", "role2"}
//	s.userStorage.On("Get", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(true, user, nil)
//	user, err := s.userSvc.RevokeRoles(s.Ctx, user.Id, []string{"role2", "role2", "role3", "role3"})
//	s.Nil(err, err)
//	s.Equal(1, len(user.Details.Roles))
//	s.Equal("role1", user.Details.Roles[0])
//}
