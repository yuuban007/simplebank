package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	mockdb "github.com/yuuban007/simplebank/db/mock"
	db "github.com/yuuban007/simplebank/db/sqlc"
	"github.com/yuuban007/simplebank/token"
	"github.com/yuuban007/simplebank/util"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	username := util.RandomOwner()
	account := randomAccount(username)

	testCases := []struct {
		name          string
		accountID     int64
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			accountID: 0,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute)
			},
			// right now this func will not be invoke
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			server := newTestServer(store, t)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			// setUpAuth
			tc.setUpAuth(t, request, server.TokenMaker)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}
}

/*
	 func TestCreateAccountAPI(t *testing.T) {
		account := randomAccount()
		// arg := db.CreateAccountParams{
		// 	Owner:    account.Owner,
		// 	Balance:  0,
		// 	Currency: account.Currency,
		// }

		testCases := []struct {
			name          string
			body          gin.H
			buildStubs    func(store *mockdb.MockStore)
			checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				body: gin.H{
					"owner": "jjj",
					"currency":"USD",
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq()).
						Times(1).
						Return(account, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},
			},
			// {
			// 	name:      "NotFound",
			// 	accountID: account.ID,
			// 	buildStubs: func(store *mockdb.MockStore) {
			// 		store.EXPECT().
			// 			GetAccount(gomock.Any(), gomock.Eq(account.ID)).
			// 			Times(1).
			// 			Return(db.Account{}, sql.ErrNoRows)
			// 	},
			// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			// 		require.Equal(t, http.StatusNotFound, recorder.Code)
			// 	},
			// },
			// {
			// 	name:      "Invalid ID",
			// 	accountID: 0,
			// 	// right now this func will not be invoke
			// 	buildStubs: func(store *mockdb.MockStore) {
			// 		store.EXPECT().
			// 			GetAccount(gomock.Any(), gomock.Any()).
			// 			Times(0)
			// 	},
			// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
			// 	},
			// },
			// {
			// 	name:      "InternalServerError",
			// 	accountID: account.ID,
			// 	buildStubs: func(store *mockdb.MockStore) {
			// 		store.EXPECT().
			// 			GetAccount(gomock.Any(), gomock.Eq(account.ID)).
			// 			Times(1).
			// 			Return(db.Account{}, sql.ErrConnDone)
			// 	},
			// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
			// 	},
			// },
			// // TODO: add more test cases
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

				url := "/accounts"
				request, err := http.NewRequest(http.MethodPost, url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, request)

				// check response
				tc.checkResponse(t, recorder)
			})

		}
	}
*/
func randomAccount(username string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
		// CreatedAt: time.Now(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
