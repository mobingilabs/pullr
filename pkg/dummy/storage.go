package dummy

import (
	"github.com/mobingilabs/pullr/pkg/domain"
)

type credential struct {
	username string
	password string
	email    string
}

type oauthsecret struct {
	secret   string
	username string
}

type tId = string
type tImageId = string
type tUsername = string
type tProvider = string

type storage struct {
	users  map[tUsername]domain.User
	images map[tUsername]map[tId]domain.Image
	builds map[tUsername]map[tImageId][]domain.Build

	authtokens      map[tId]string
	authcredentials map[tUsername]credential
	oauthsecrets    map[string]oauthsecret
	oauthtokens     map[tUsername]map[tProvider]string
}

func NewStorageDriver(opts map[string]interface{}) domain.StorageDriver {
	return &storage{
		users:           make(map[string]domain.User),
		images:          make(map[string]map[string]domain.Image),
		builds:          make(map[string]map[string][]domain.Build),
		authtokens:      make(map[string]string),
		authcredentials: make(map[string]credential),
		oauthsecrets:    make(map[string]oauthsecret),
		oauthtokens:     make(map[string]map[string]string),
	}
}

// Close closes nothing
func (*storage) Close() error {
	return nil
}

// AuthStorage creates an AuthStorage instance
func (s *storage) AuthStorage() domain.AuthStorage {
	return &authStorage{s}
}

// OAuthStorage creates an OAuthStorage instance
func (s *storage) OAuthStorage() domain.OAuthStorage {
	return &oauthStorage{s}
}

// UserStorage creates an UserStorage instance
func (s *storage) UserStorage() domain.UserStorage {
	return &userStorage{s}
}

// ImageStorage creates an ImageStorage instance
func (s *storage) ImageStorage() domain.ImageStorage {
	return &imageStorage{s}
}

// BuildStorage creates an BuildStorage instance
func (s *storage) BuildStorage() domain.BuildStorage {
	return &buildStorage{s}
}

// AuthStorage ================================================================
type authStorage struct {
	d *storage
}

func (s *authStorage) Get(username string) (string, error) {
	c, ok := s.d.authcredentials[username]
	if !ok {
		return "", domain.ErrNotFound
	}

	return c.password, nil
}

func (s *authStorage) GetByEmail(email string) (string, error) {
	for _, c := range s.d.authcredentials {
		if c.email == email {
			return c.password, nil
		}
	}

	return "", domain.ErrNotFound
}

func (s *authStorage) GetRefreshToken(tokenID string) (string, error) {
	token, ok := s.d.authtokens[tokenID]
	if !ok {
		return "", domain.ErrNotFound
	}

	return token, nil
}

func (s *authStorage) PutRefreshToken(username string, tokenID string) error {
	s.d.authtokens[tokenID] = tokenID
	return nil
}

func (s *authStorage) DeleteRefreshToken(tokenID string) error {
	delete(s.d.authtokens, tokenID)
	return nil
}

func (s *authStorage) PutCredentials(username string, email string, password string) error {
	s.d.authcredentials[username] = credential{username, password, email}
	return nil
}

func (s *authStorage) Delete(username string) error {
	delete(s.d.authcredentials, username)
	return nil
}

// OAuthStorage ================================================================
type oauthStorage struct {
	d *storage
}

func (s *oauthStorage) PutSecret(username, secret string) error {
	s.d.oauthsecrets[secret] = oauthsecret{username, secret}
	return nil
}

func (s *oauthStorage) PopSecret(secret string) (string, error) {
	sec, ok := s.d.oauthsecrets[secret]
	if !ok {
		return "", domain.ErrNotFound
	}

	return sec.username, nil
}

func (s *oauthStorage) GetTokens(username string) (map[string]string, error) {
	tokens, ok := s.d.oauthtokens[username]
	if !ok {
		return make(map[string]string), nil
	}

	return tokens, nil
}

func (s *oauthStorage) PutToken(username string, provider string, token string) error {
	tokens, ok := s.d.oauthtokens[username]
	if !ok {
		s.d.oauthtokens[username] = make(map[string]string)
		tokens = s.d.oauthtokens[username]
	}

	tokens[provider] = token
	return nil
}

func (s *oauthStorage) RemoveToken(username string, provider string) error {
	tokens, ok := s.d.oauthtokens[username]
	if !ok {
		return nil
	}

	delete(tokens, provider)
	return nil
}

// UserStorage ================================================================
type userStorage struct {
	d *storage
}

func (s *userStorage) Get(username string) (domain.User, error) {
	usr, ok := s.d.users[username]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}

	return usr, nil
}

func (s *userStorage) GetByEmail(email string) (domain.User, error) {
	for _, usr := range s.d.users {
		if usr.Email == email {
			return usr, nil
		}
	}

	return domain.User{}, domain.ErrNotFound
}

func (s *userStorage) Put(user domain.User) error {
	s.d.users[user.Username] = user
	return nil
}

func (s *userStorage) List(opts domain.ListOptions) ([]domain.User, domain.Pagination, error) {
	users := sortUsers(s.d.users)
	pagination := domain.Pagination{
		Last:    maxInt((len(users)/opts.PerPage)-1, 0),
		Current: opts.Page,
	}

	skip := opts.PerPage * opts.Page
	if skip > len(users) {
		return users[pagination.Last*opts.PerPage:], pagination, nil
	}
	if skip+opts.PerPage > len(users) {
		return users[skip:], pagination, nil
	}

	return users[skip : skip+opts.PerPage], pagination, nil
}

func (s *userStorage) Delete(username string) error {
	delete(s.d.users, username)
	return nil
}

// ImageStorage ================================================================
type imageStorage struct {
	d *storage
}

func (s *imageStorage) Get(username string, key string) (domain.Image, error) {
	usrImages, ok := s.d.images[username]
	if !ok {
		return domain.Image{}, domain.ErrNotFound
	}

	img, ok := usrImages[key]
	if !ok {
		return domain.Image{}, domain.ErrNotFound
	}

	return img, nil
}

func (s *imageStorage) GetMany(username string, keys []string) (map[string]domain.Image, error) {
	usrImages, ok := s.d.images[username]
	if !ok {
		return nil, nil
	}

	foundImages := make(map[string]domain.Image, len(keys))
	for key, img := range usrImages {
		for i := range keys {
			if keys[i] == key {
				foundImages[key] = img
				break
			}
		}
	}

	return foundImages, nil
}

func (s *imageStorage) List(username string, opts domain.ListOptions) ([]domain.Image, domain.Pagination, error) {
	usrImages, ok := s.d.images[username]
	if !ok {
		return []domain.Image{}, domain.Pagination{}, nil
	}

	pagination := domain.Pagination{
		Last:    maxInt((len(usrImages)/opts.PerPage)-1, 0),
		Current: opts.Page,
	}

	sortedImages := sortImages(usrImages)
	skip := opts.PerPage * opts.Page
	if skip > len(sortedImages) {
		return sortedImages[pagination.Last*opts.PerPage:], pagination, nil
	}
	if skip+opts.PerPage > len(sortedImages) {
		return sortedImages[skip:], pagination, nil
	}

	return sortedImages[skip : skip+opts.PerPage], pagination, nil
}

func (s *imageStorage) Put(username string, image domain.Image) error {
	usrImages, ok := s.d.images[username]
	if !ok {
		s.d.images[username] = make(map[tId]domain.Image)
		usrImages = s.d.images[username]
	}

	usrImages[domain.ImageKey(image)] = image
	return nil
}

func (s *imageStorage) Update(username string, key string, image domain.Image) error {
	usrImages, ok := s.d.images[username]
	if !ok {
		return domain.ErrNotFound
	}

	_, ok = usrImages[key]
	if !ok {
		return domain.ErrNotFound
	}

	if key != domain.ImageKey(image) {
		delete(usrImages, key)
	}

	usrImages[key] = image
	return nil
}

func (s *imageStorage) Delete(username string, key string) error {
	usrImages, ok := s.d.images[username]
	if !ok {
		return domain.ErrNotFound
	}

	_, ok = usrImages[key]
	if !ok {
		return domain.ErrNotFound
	}

	delete(usrImages, key)
	return nil
}

// BuildStorage ================================================================

type buildStorage struct {
	d *storage
}

func (s *buildStorage) GetAll(username string, imgKey string, opts domain.ListOptions) ([]domain.Build, domain.Pagination, error) {
	builds, ok := s.d.builds[username][imgKey]
	if !ok {
		return nil, domain.Pagination{}, domain.ErrNotFound
	}

	pagination := domain.Pagination{
		Last:    maxInt((len(builds)/opts.PerPage)-1, 0),
		Current: opts.Page,
	}

	sortedBuilds := sortBuilds(builds)
	skip := opts.PerPage * opts.Page
	if skip > len(sortedBuilds) {
		return sortedBuilds[pagination.Last*opts.PerPage:], pagination, nil
	}
	if skip+opts.PerPage > len(sortedBuilds) {
		return sortedBuilds[skip:], pagination, nil
	}

	return sortedBuilds[skip : skip+opts.PerPage], pagination, nil
}

func (s *buildStorage) GetLast(username string, imgKey string) (domain.Build, error) {
	imgBuilds, ok := s.d.builds[username][imgKey]
	if !ok {
		return domain.Build{}, domain.ErrNotFound
	}

	if len(imgBuilds) == 0 {
		return domain.Build{}, domain.ErrNotFound
	}

	return imgBuilds[0], nil
}

func (s *buildStorage) List(username string, opts domain.ListOptions) ([]domain.Build, domain.Pagination, error) {
	images, ok := s.d.builds[username]
	if !ok {
		return nil, domain.Pagination{}, domain.ErrNotFound
	}

	pagination := domain.Pagination{
		Last:    maxInt((len(images)/opts.PerPage)-1, 0),
		Current: opts.Page,
	}

	sortedBuilds := sortImageBuilds(images)
	skip := opts.PerPage * opts.Page
	if skip > len(sortedBuilds) {
		return sortedBuilds[pagination.Last*opts.PerPage:], pagination, nil
	}
	if skip+opts.PerPage > len(sortedBuilds) {
		return sortedBuilds[skip:], pagination, nil
	}

	return sortedBuilds[skip : skip+opts.PerPage], pagination, nil
}

func (s *buildStorage) Update(username string, imgKey string, build domain.Build) error {
	builds, ok := s.d.builds[username][imgKey]
	if !ok || len(builds) == 0 {
		return domain.ErrNotFound
	}

	builds[0] = build
	return nil
}

func (s *buildStorage) Put(username string, imgKey string, build domain.Build) error {
	usrImgs, ok := s.d.builds[username]
	if !ok {
		s.d.builds[username] = make(map[string][]domain.Build)
		usrImgs = s.d.builds[username]
	}

	_, ok = usrImgs[imgKey]
	if !ok {
		usrImgs[imgKey] = nil
	}

	usrImgs[imgKey] = append(usrImgs[imgKey])
	return nil
}
