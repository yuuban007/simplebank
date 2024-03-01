package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yuuban007/simplebank/db/mock"
	db "github.com/yuuban007/simplebank/db/sqlc"
	"github.com/yuuban007/simplebank/util"
	"go.uber.org/mock/gomock"
)

// struct to implement gomock.Matcher interface
type eqCreateUserParamsMatcher struct {
	password string
	arg      db.CreateUserParams
}

// Matches returns whether x is equal to the argument of this matcher.
func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword // update the hashed password
	return reflect.DeepEqual(e.arg, arg)
}

// String returns a description of the matcher. It is used in error messages.
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg: %v and password %v", e.arg, e.password)
}

// eqCreateUserParams is a custom gomock matcher for db.CreateUserParams
func eqCreateUserParams(password string, arg db.CreateUserParams) gomock.Matcher {
	return eqCreateUserParamsMatcher{password: password, arg: arg}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					FullName:       user.FullName,
					Email:          user.Email,
					HashedPassword: hashedPassword,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParams(password, arg)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)
			// start a server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/users"

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}
}

func randomUser() (user db.User, password string) {
	user = db.User{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}
	password = util.RandomString(8)
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}
