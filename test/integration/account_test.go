package integration_test

import (
	"github.com/ecommerce-api/internal/repository"
	"github.com/ecommerce-api/internal/service"
	"github.com/ecommerce-api/pkg/dto"
	security2 "github.com/ecommerce-api/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func registerAndSetPayload(s *RepositoryTestSuite) {
	user := s.registerUserData()

	security2.PayloadData = &security2.Payload{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func (s *RepositoryTestSuite) TestAccount() {
	registerAndSetPayload(s)
	accountService := service.NewAccountService(s.db)

	s.Run("get my profile", func() {
		profile, err := accountService.Profile()

		assert.Nil(s.T(), err)
		assert.Equal(s.T(), "user@gmail.com", profile.Email)
		assert.Equal(s.T(), "user", profile.Username)

	})

}

// [note] : mailer is disabled in test
func (s *RepositoryTestSuite) TestRegisterUserShouldExpected() {
	s.registerUserData()

	//find user expectation
	expectEmail := "user@gmail.com"
	repo := repository.NewUserRepository(s.db)
	user, err := repo.ByEmail(expectEmail)

	require.NoError(s.T(), err)
	require.Equal(s.T(), expectEmail, user.Email)

}

// TestCreateAccessTokenThenAuthorizedShouldExpected : Login
func (s *RepositoryTestSuite) TestLoginThenAuthorizedShouldExpected() {
	user := s.registerUserData()

	token, exp, err := security2.CreateToken(user)

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), token)
	require.NotEmpty(s.T(), exp)

	jwt, err := security2.ParseToken(token)
	jwtRaw := *jwt

	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(jwtRaw["user_id"].(float64)), user.ID)
}

func (s *RepositoryTestSuite) TestGetLoginShouldReturnCredential() {
	user := s.registerUserData()

	authService := service.NewAccountService(s.db)

	request := dto.AuthAccessTokenRequest{
		Email:    user.Email,
		Password: "localhost",
	}

	cred, err := authService.Login(&request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), cred.Username, user.Username)
	assert.Equal(s.T(), cred.Email, user.Email)
	assert.Equal(s.T(), cred.Role, user.Role)

	assert.Nil(s.T(), err)

}
