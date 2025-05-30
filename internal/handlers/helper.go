package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/SanExpett/diploma/internal/domain"
	myerrors "github.com/SanExpett/diploma/internal/errors"
	"github.com/SanExpett/diploma/internal/metrics"
	session "github.com/SanExpett/diploma/internal/session/proto"
)

type SuccessResponse struct {
	Status int `json:"status"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

func WriteSuccess(w http.ResponseWriter, req *http.Request, metrics *metrics.HttpMetrics) error {
	response := SuccessResponse{
		Status: http.StatusOK,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return err
	}

	curRoute := mux.CurrentRoute(req)
	if curRoute == nil {
		return fmt.Errorf("unable to get current route")
	}

	pathTemplate, err := curRoute.GetPathTemplate()
	if err != nil {
		return fmt.Errorf("unable to get path template: %w", err)
	}

	metrics.IncRequestsTotal(pathTemplate, req.Method, 200)

	return nil
}

func WriteError(w http.ResponseWriter, req *http.Request, metrics *metrics.HttpMetrics, err error) error {
	statusCode, err := myerrors.ParseError(err)

	response := ErrorResponse{
		Status: statusCode,
		Err:    err.Error(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		return err
	}

	curRoute := mux.CurrentRoute(req)
	if curRoute == nil {
		return fmt.Errorf("unable to get current route")
	}

	pathTemplate, err := curRoute.GetPathTemplate()
	if err != nil {
		return fmt.Errorf("unable to get path template: %w", err)
	}

	metrics.IncRequestsTotal(pathTemplate, req.Method, statusCode)

	return nil
}

func WriteResponse(w http.ResponseWriter, r *http.Request, metrics *metrics.HttpMetrics, jsonResponse any,
	requestID any) error {
	curRoute := mux.CurrentRoute(r)
	if curRoute == nil {
		return fmt.Errorf("unable to get current route")
	}

	pathTemplate, err := curRoute.GetPathTemplate()
	if err != nil {
		return fmt.Errorf("unable to get path template: %w", err)
	}

	metrics.IncRequestsTotal(pathTemplate, r.Method, 200)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse.([]byte))
	if err != nil {
		return fmt.Errorf("[reqid=%s] failed to write response: %v\n", requestID, err)
	}

	return nil
}

func escapeUserData(user *domain.User) {
	user.Name = html.EscapeString(user.Name)
	user.Email = html.EscapeString(user.Email)
	user.Password = html.EscapeString(user.Password)
	user.Avatar = html.EscapeString(user.Avatar)
}

func escapeUserPreviewData(userPreview *domain.UserPreview) {
	userPreview.Name = html.EscapeString(userPreview.Name)
	userPreview.Avatar = html.EscapeString(userPreview.Avatar)
}

func escapeActorData(actor *domain.ActorData) {
	actor.Name = html.EscapeString(actor.Name)
	actor.Avatar = html.EscapeString(actor.Avatar)
	actor.Spouse = html.EscapeString(actor.Spouse)
	actor.BirthPlace = html.EscapeString(actor.BirthPlace)
	actor.Career = html.EscapeString(actor.Career)
}

func escapeFilmData(filmData *domain.FilmData) {
	filmData.Title = html.EscapeString(filmData.Title)
	filmData.Data = html.EscapeString(filmData.Data)
	filmData.Director = html.EscapeString(filmData.Director)
	filmData.Preview = html.EscapeString(filmData.Preview)
	var genres []domain.Genre
	for _, genre := range filmData.Genres {
		genre.Name = html.EscapeString(genre.Name)
		genres = append(genres, genre)
	}
	filmData.Genres = genres
}

func escapeSerialData(filmData *domain.SerialData) {
	filmData.Title = html.EscapeString(filmData.Title)
	filmData.Data = html.EscapeString(filmData.Data)
	filmData.Director = html.EscapeString(filmData.Director)
	filmData.Preview = html.EscapeString(filmData.Preview)
}

func escapeActorPreview(actor *domain.ActorPreview) {
	actor.Name = html.EscapeString(actor.Name)
	actor.Avatar = html.EscapeString(actor.Avatar)
}

func escapeFilmPreview(film *domain.FilmPreview) {
	film.Preview = html.EscapeString(film.Preview)
	film.Director = html.EscapeString(film.Director)
	film.Title = html.EscapeString(film.Title)
}

func escapeTopFilm(film *domain.TopFilm) {
	film.Preview = html.EscapeString(film.Preview)
	film.Title = html.EscapeString(film.Title)
	film.Data = html.EscapeString(film.Data)
}

func escapeComment(comment *domain.Comment) {
	comment.Text = html.EscapeString(comment.Text)
	comment.Author = html.EscapeString(comment.Author)
}

func escapeGenreFilms(genreFilms *domain.GenreFilms) {
	genreFilms.Name = html.EscapeString(genreFilms.Name)
}

func convertFilmPreviewToRegular(film *session.FilmPreview) domain.FilmPreview {
	filmNew := domain.FilmPreview{
		Uuid:         film.Uuid,
		Title:        film.Title,
		Preview:      film.Preview,
		Director:     film.Director,
		AverageScore: film.AvgScore,
		ScoresCount:  film.ScoresCount,
		AgeLimit:     film.AgeLimit,
		IsSerial:     film.IsSerial,
		Duration:     film.Duration,
	}
	return filmNew
}

func convertLongFilmPreviewToRegular(film *session.FindFilmLong) domain.FilmData {
	var genres []domain.Genre
	for _, genre := range film.Genres {
		genres = append(genres, domain.Genre{Name: genre.Name, Uuid: genre.Uuid})
	}

	return domain.FilmData{
		Uuid:         film.Uuid,
		Title:        film.Title,
		Preview:      film.Preview,
		Director:     film.Director,
		Date:         convertProtoToTime(film.Date),
		AgeLimit:     film.AgeLimit,
		AverageScore: film.AvgScore,
		ScoresCount:  film.ScoresCount,
		IsSerial:     film.IsSerial,
		Duration:     film.Duration,
		Genres:       genres,
	}
}

func convertFilmDataToRegular(film *session.FilmData) domain.FilmData {
	var genres []domain.Genre
	for _, genre := range film.Genres {
		genres = append(genres, domain.Genre{Name: genre.Name, Uuid: genre.Uuid})
	}

	return domain.FilmData{
		Uuid:         film.Uuid,
		Title:        film.Title,
		Preview:      film.Preview,
		Director:     film.Director,
		IsSerial:     film.IsSerial,
		Link:         film.Link,
		Data:         film.Data,
		Date:         convertProtoToTime(film.Date),
		AgeLimit:     film.AgeLimit,
		AverageScore: film.AvgScore,
		ScoresCount:  film.ScoresCount,
		Duration:     film.Duration,
		Genres:       genres,
		WithSub:      film.WithSub,
	}
}

func convertSerialDataToRegular(film *session.FilmData) domain.SerialData {
	seasons := make([]domain.Season, 0, len(film.Seasons))
	for _, season := range film.Seasons {
		episodes := make([]domain.Episode, 0, len(season.Episodes))
		for _, episode := range season.Episodes {
			episodes = append(episodes, domain.Episode{
				Link:  episode.Link,
				Title: episode.Title,
			})
		}
		seasons = append(seasons, domain.Season{
			Series: episodes,
		})
	}

	var genres []domain.Genre
	for _, genre := range film.Genres {
		genres = append(genres, domain.Genre{Name: genre.Name, Uuid: genre.Uuid})
	}

	return domain.SerialData{
		Uuid:         film.Uuid,
		Title:        film.Title,
		Preview:      film.Preview,
		Director:     film.Director,
		IsSerial:     film.IsSerial,
		Seasons:      seasons,
		Data:         film.Data,
		Date:         convertProtoToTime(film.Date),
		AgeLimit:     film.AgeLimit,
		AverageScore: film.AvgScore,
		ScoresCount:  film.ScoresCount,
		Genres:       genres,
		WithSub:      film.WithSub,
	}
}

func convertCommentToRegular(comment *session.Comment) domain.Comment {
	return domain.Comment{
		Uuid:       comment.Uuid,
		FilmUuid:   comment.FilmUuid,
		Text:       comment.Text,
		Author:     comment.Author,
		AuthorUuid: comment.AuthorUuid,
		Score:      comment.Score,
		AddedAt:    convertProtoToTime(comment.AddedAt),
	}
}

func convertActorPreviewToRegular(actor *session.ActorPreview) domain.ActorPreview {
	return domain.ActorPreview{
		Uuid:   actor.Uuid,
		Name:   actor.Name,
		Avatar: actor.Avatar,
	}
}

func convertActorPreviewLongToRegular(actor *session.ActorPreviewLong) domain.ActorData {
	return domain.ActorData{
		Uuid:       actor.Uuid,
		Name:       actor.Name,
		Avatar:     actor.Avatar,
		Birthday:   convertProtoToTime(actor.Birthday),
		Career:     actor.Career,
		BirthPlace: actor.BirthPlace,
	}
}

func convertActorDataToRegular(actor *session.ActorData) domain.ActorData {
	var filmsPreview []domain.FilmPreview
	for _, film := range actor.FilmsPreviews {
		filmRegular := convertFilmPreviewToRegular(film)
		escapeFilmPreview(&filmRegular)
		filmsPreview = append(filmsPreview, filmRegular)
	}
	return domain.ActorData{
		Uuid:       actor.Uuid,
		Name:       actor.Name,
		Avatar:     actor.Avatar,
		Birthday:   convertProtoToTime(actor.Birthday),
		BirthPlace: actor.Birthplace,
		Career:     actor.Career,
		Spouse:     actor.Spouse,
		Films:      filmsPreview,
		Height:     actor.Height,
	}
}

func convertProtoToTime(protoTime *timestamppb.Timestamp) time.Time {
	return protoTime.AsTime()
}

func convertUserToRegular(user *session.User) domain.User {
	return domain.User{
		Uuid:            user.Uuid,
		Email:           user.Email,
		Password:        user.Password,
		Name:            user.Username,
		Version:         user.Version,
		IsAdmin:         user.IsAdmin,
		Avatar:          user.Avatar,
		Birthday:        convertProtoToTime(user.Birthday),
		RegisteredAt:    convertProtoToTime(user.RegisteredAt),
		HasSubscription: user.HasSubscription,
	}
}

func convertUserSignUpDataToRegular(userData domain.UserSignUp) *session.UserSignUp {
	return &session.UserSignUp{
		Email:    userData.Email,
		Password: userData.Password,
		Username: userData.Name,
	}
}

func convertUserPreviewToRegular(user *session.UserPreview) domain.UserPreview {
	return domain.UserPreview{
		Uuid:   user.Uuid,
		Name:   user.Username,
		Avatar: user.Avatar,
	}
}

func convertGenreFilmsToRegular(genreFilms *session.GenreFilms) domain.GenreFilms {
	var filmsConverted []domain.FilmPreview
	for _, film := range genreFilms.Films {
		filmsConverted = append(filmsConverted, convertFilmPreviewToRegular(film))
	}
	return domain.GenreFilms{
		Name:  genreFilms.Genre,
		Uuid:  genreFilms.GenreUuid,
		Films: filmsConverted,
	}
}

func convertTimeToProto(time time.Time) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}

func convertActorToAddToRegular(actor domain.ActorToAdd) *session.ActorDataToAdd {
	return &session.ActorDataToAdd{
		Name:       actor.Name,
		Avatar:     actor.Avatar,
		BirthPlace: actor.BirthPlace,
		BirthdayAt: convertTimeToProto(actor.Birthday),
		Career:     actor.Career,
		Spouse:     actor.Spouse,
		Height:     actor.Height,
	}
}

func convertFilmToAdd(filmToAdd domain.FilmToAdd) *session.FilmToAdd {
	seasons := make([]*session.Season, 0, len(filmToAdd.FilmData.Seasons))
	for _, season := range filmToAdd.FilmData.Seasons {
		episodes := make([]*session.Episode, 0, len(season.Series))
		for _, episode := range season.Series {
			episodes = append(episodes, &session.Episode{
				Link:  episode.Link,
				Title: episode.Title,
			})
		}

		seasons = append(seasons, &session.Season{
			Episodes: episodes,
		})
	}

	filmData := session.FilmDataToAdd{
		Title:            filmToAdd.FilmData.Title,
		IsSerial:         filmToAdd.FilmData.IsSerial,
		Preview:          filmToAdd.FilmData.Preview,
		Director:         filmToAdd.FilmData.Director,
		Data:             filmToAdd.FilmData.Data,
		AgeLimit:         filmToAdd.FilmData.AgeLimit,
		PublishedAt:      convertTimeToProto(filmToAdd.FilmData.PublishedAt),
		Genres:           filmToAdd.FilmData.Genres,
		Duration:         filmToAdd.FilmData.Duration,
		Link:             filmToAdd.FilmData.Link,
		Seasons:          seasons,
		WithSubscription: filmToAdd.FilmData.WithSubscription,
	}

	var actors []*session.ActorDataToAdd
	for _, act := range filmToAdd.Actors {
		actors = append(actors, convertActorToAddToRegular(act))
	}

	directorData := session.DirectorDataToAdd{
		Name:     filmToAdd.DirectorToAdd.Name,
		Birthday: convertTimeToProto(filmToAdd.DirectorToAdd.Birthday),
		Avatar:   filmToAdd.DirectorToAdd.Avatar,
	}

	return &session.FilmToAdd{
		FilmData: &filmData,
		Actors:   actors,
		Director: &directorData,
	}
}

func convertTopFilmToRegular(film *session.TopFilm) domain.TopFilm {
	return domain.TopFilm{
		Uuid:     film.Uuid,
		Preview:  film.Preview,
		Title:    film.Title,
		IsSerial: film.IsSerial,
		Data:     film.Data,
	}
}

func convertCommentToAddToProto(commentToAdd domain.CommentToAdd) *session.CommentToAdd {
	return &session.CommentToAdd{
		FilmUuid:   commentToAdd.FilmUuid,
		AuthorUuid: commentToAdd.AuthorUuid,
		Text:       commentToAdd.Text,
		Score:      commentToAdd.Score,
	}
}

func convertCommentToRemoveToProto(comment domain.CommentToRemove) *session.CommentToRemove {
	return &session.CommentToRemove{
		FilmUuid:   comment.FilmUuid,
		AuthorUuid: comment.AuthorUuid,
	}
}

func convertSubsToRegular(subs []*session.Subscription) []domain.Subscription {
	subscriptions := make([]domain.Subscription, 0, len(subs))
	for _, sub := range subs {
		subscriptions = append(subscriptions, domain.Subscription{
			Uuid:        sub.Uuid,
			Title:       sub.Title,
			Amount:      sub.Price,
			Description: sub.Description,
			Duration:    sub.Duration,
		})
	}
	return subscriptions
}

func IsTokenValid(token *http.Cookie, secretKey string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %w", myerrors.ErrNotAuthorised)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: %w", myerrors.ErrNotAuthorised)
	}

	_, ok = claims["Login"]
	if !ok {
		return nil, fmt.Errorf("invalid token: %w", myerrors.ErrNotAuthorised)
	}
	_, ok = claims["IsAdmin"]
	if !ok {
		return nil, fmt.Errorf("invalid token: %w", myerrors.ErrNotAuthorised)
	}

	_, ok = claims["Version"]
	if !ok {
		return nil, fmt.Errorf("invalid token: %w", myerrors.ErrNotAuthorised)
	}

	return claims, nil
}

func ValidateLogin(e string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if emailRegex.MatchString(e) {
		return nil
	}
	return myerrors.ErrLoginIsNotValid
}

func ValidateUsername(username string) error {
	if len(username) >= 4 {
		return nil
	}
	return myerrors.ErrUsernameIsToShort
}

func ValidatePassword(password string) error {
	if len(password) >= 6 {
		return nil
	}
	return myerrors.ErrPasswordIsToShort
}

type customClaims struct {
	jwt.StandardClaims
	Login   string
	IsAdmin bool
	Version uint32
}

func GenerateTokens(login string, isAdmin bool, version uint32) (tokenSigned string,
	err error) {
	tokenCustomClaims := customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    "nimbus",
		},
		Login:   login,
		IsAdmin: isAdmin,
		Version: version,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenCustomClaims)

	tokenSigned, err = token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return tokenSigned, nil
}

func Encode(file *multipart.FileHeader) (string, error) {
	allowedExtensions := map[string]bool{
		".jpg": true,
		".png": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("wrong extensions, allowed extensions are .jpg and .png")
	}

	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		return "", fmt.Errorf("file is too big, max size is 10 MB")
	}

	avatar, err := file.Open()
	if err != nil {
		return "", err
	}
	defer avatar.Close()

	fileBytes, err := io.ReadAll(avatar)
	if err != nil {
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(fileBytes)
	return base64String, nil
}
