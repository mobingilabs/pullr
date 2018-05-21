package domain_test

import (
	"os"
	"testing"

	. "github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/dummy"
)

var (
	expiredAuthToken    = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjEwNTM0NDYsInN1YiI6InRlc3QifQ.KQajE5tMGO0aJDScicFjUHGq7ly8PXekEEy_N3c_HdoNAMQdhVleSY_bPYFf1MN88k5R3zP8wc9uJDaivzykJBdCZsTNjU0rORmVXipzDFot814ebPBWnDlkYbr8fa6du3oQmJosfzTv2EAGPvN-Ra500o8ErkJetHFU7DAHuO283wf6CNzQtCc5sarkEX7MAIJnfn_VI_51gExP0a6iMQlAKbCPVgdatyJeX708Nl3ZnPRm0Xr9CJRkOxq3zIyKVzxAeZw8Vptvik0TtDi-avMOfL8694uPwTAuzVxbIthaAcdekug_uNhKPrR7R3GWtfQt1ojo649nEcCfGNUQjg"
	expiredRefreshToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjA2MjI0NTgsImp0aSI6InZBWTJJM0s0cjB6VENjQmlkSFNsb25UdlJOcFpObWItV3luVzIyejkwbUU9Iiwic3ViIjoidGVzdCJ9.vtnxFzzt_AE3MheTwa5xsd2uboDRJfQVmbZaLFb97MsVOcp6s_i2KTLpCwRgJIfEOPx47NtgRICBILdiHjjqJQexLsK9EWaEY7wV5Q9gpEAlmWA8MoBIL0EijgdjtCTrBtCAU8j7oyU6B2fcPXVCvVlWXaRjjtpLoBPFyW7ZcwMc5OdYt3cjpJ73NPirpzm27KWZvb_Qv5W-Km8I49vc6iASFQcSMqj3TDOnHPIZfVH8cm0yOjc0XFnigJ8DmROcBv_IRiFHfB_Dpeay5YsMOzD1k6BTivZ2cY9FaBlDW1xPOsrFBCQiOSHJ0mnJPjW5bJUO_QrPpKa7DtG7tW021Q"
)

func newAuthService(t *testing.T, storage StorageDriver) *DefaultAuthService {
	authsvc, err := NewAuthService(storage.AuthStorage(), storage.UserStorage(), &TestLogger{}, AuthConfig{
		Key: "../../certs/auth.key",
		Crt: "../../certs/auth.crt",
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return authsvc
}

func TestNewAuthService(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)

	configTests := []struct {
		name string
		key  string
		crt  string
	}{
		{"wrong crt path", "wrong", "../../certs/auth.crt"},
		{"wrong key path", "../../certs/auth.key", "wrong"},
	}

	for _, tt := range configTests {
		_, err := NewAuthService(storage.AuthStorage(), storage.UserStorage(), &TestLogger{}, AuthConfig{
			Key: tt.key,
			Crt: tt.crt,
		})

		if _, ok := err.(*os.PathError); !ok {
			t.Errorf("for %s expected to have an open error, got: %v", tt.name, err)
		}
	}
}

func TestAuthService_Register(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	authsvc := newAuthService(t, storage)

	err := authsvc.Register("test", "test@test.com", "12345")
	if err != nil {
		t.Fatal(err)
	}

	cred, err := storage.AuthStorage().Get("test")
	if err != nil {
		t.Fatal(err)
	}

	if cred == "" {
		t.Fatal("hashed password shouldn't be empty")
	}
}

func TestAuthService_Register_SameIdentifier(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	authsvc := newAuthService(t, storage)

	err := authsvc.Register("test", "test@test.com", "12345")
	if err != nil {
		t.SkipNow()
	}

	err = authsvc.Register("test", "another@email.com", "abcdef")
	if err != ErrUserUsernameExist {
		t.Error("register didn't return error for existent username")
	}

	err = authsvc.Register("john", "test@test.com", "12345")
	if err != ErrUserEmailExist {
		t.Error("register didn't return error for existent email")
	}
}

func TestAuthService_Login(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	storage.AuthStorage().PutCredentials("test", "test@test.com", "$2a$10$q/3Fqkzv0kR/cXdzdKgkl.P3LSmKZ24s84SV.hvDkCDTDCK1i4GwG")
	authsvc := newAuthService(t, storage)

	secrets, err := authsvc.Login("test", "12345")
	if err != nil {
		t.Fatal(err)
	}

	if secrets.Username == "" || secrets.AuthToken == "" || secrets.RefreshToken == "" {
		t.Errorf("auth secrets is empty")
	}

	secrets, err = authsvc.Login("test", "abcde")
	if err != ErrAuthBadCredentials {
		t.Error("wrong credentials didn't failed")
	}
}

func TestAuthService_Grant(t *testing.T) {
	storage := dummy.NewStorageDriver(nil)
	storage.AuthStorage().PutCredentials("test", "test@test.com", "$2a$10$q/3Fqkzv0kR/cXdzdKgkl.P3LSmKZ24s84SV.hvDkCDTDCK1i4GwG")
	authsvc := newAuthService(t, storage)

	loginSecrets, err := authsvc.Login("test", "12345")
	if err != nil {
		t.Fatal(err)
	}

	tokenTests := []struct {
		name    string
		auth    string
		refresh string
		err     error
	}{
		{"both tokens given", loginSecrets.AuthToken, loginSecrets.RefreshToken, nil},
		{"auth token empty", "", loginSecrets.RefreshToken, ErrAuthUnauthorized},
		{"refresh token empty", loginSecrets.AuthToken, "", ErrAuthUnauthorized},
		{"both tokens empty", "", "", ErrAuthUnauthorized},
		{"invalid auth token", "xxx", loginSecrets.RefreshToken, ErrAuthUnauthorized},
		{"invalid refresh token", loginSecrets.AuthToken, "xxx", ErrAuthUnauthorized},
		{"expired auth token", expiredAuthToken, loginSecrets.RefreshToken, nil},
		{"expired auth & refresh token", expiredAuthToken, expiredRefreshToken, ErrAuthUnauthorized},
	}

	for _, tt := range tokenTests {
		tt := tt
		_, err := authsvc.Grant(tt.refresh, tt.auth)
		if err != tt.err {
			t.Errorf("for %s error expected to be %v but got %v", tt.name, tt.err, err)
		}
	}
}
