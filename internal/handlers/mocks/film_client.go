package mocks

import (
	"context"

	proto "github.com/SanExpett/diploma/internal/session/proto"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=film_client.go -destination=films_client_mock.go -package=mocks
type FilmsClient interface {
	GetAllFilmsPreviews(ctx context.Context, in *proto.AllFilmsPreviewsRequest, opts ...grpc.CallOption) (*proto.AllFilmsPreviewsResponse, error)
	GetFilmsPreviewsWithSub(ctx context.Context, in *proto.AllFilmsPreviewsRequest, opts ...grpc.CallOption) (*proto.AllFilmsPreviewsResponse, error)
	GetFilmDataByUuid(ctx context.Context, in *proto.FilmDataByUuidRequest, opts ...grpc.CallOption) (*proto.FilmDataByUuidResponse, error)
	GetFilmPreviewByUuid(ctx context.Context, in *proto.FilmPreviewByUuidRequest, opts ...grpc.CallOption) (*proto.FilmPreviewByUuidResponse, error)
	RemoveFilmByUuid(ctx context.Context, in *proto.RemoveFilmByUuidRequest, opts ...grpc.CallOption) (*proto.RemoveFilmByUuidResponse, error)
	GetActorDataByUuid(ctx context.Context, in *proto.ActorDataByUuidRequest, opts ...grpc.CallOption) (*proto.ActorDataByUuidResponse, error)
	GetActorsByFilm(ctx context.Context, in *proto.ActorsByFilmRequest, opts ...grpc.CallOption) (*proto.ActorsByFilmResponse, error)
	PutFavorite(ctx context.Context, in *proto.PutFavoriteRequest, opts ...grpc.CallOption) (*proto.PutFavoriteResponse, error)
	DeleteFavorite(ctx context.Context, in *proto.DeleteFavoriteRequest, opts ...grpc.CallOption) (*proto.DeleteFavoriteResponse, error)
	GetAllFavoriteFilms(ctx context.Context, in *proto.GetAllFavoriteFilmsRequest, opts ...grpc.CallOption) (*proto.GetAllFavoriteFilmsResponse, error)
	GetAllFilmsByGenre(ctx context.Context, in *proto.GetAllFilmsByGenreRequest, opts ...grpc.CallOption) (*proto.GetAllFilmsByGenreResponse, error)
	GetAllGenres(ctx context.Context, in *proto.GetAllGenresRequest, opts ...grpc.CallOption) (*proto.GetAllGenresResponse, error)
	AddFilm(ctx context.Context, in *proto.AddFilmRequest, opts ...grpc.CallOption) (*proto.AddFilmResponse, error)
	FindFilmsShort(ctx context.Context, in *proto.FindFilmsShortRequest, opts ...grpc.CallOption) (*proto.FindFilmsShortResponse, error)
	FindFilmsLong(ctx context.Context, in *proto.FindFilmsShortRequest, opts ...grpc.CallOption) (*proto.FindFilmsLongResponse, error)
	FindSerialsShort(ctx context.Context, in *proto.FindFilmsShortRequest, opts ...grpc.CallOption) (*proto.FindFilmsShortResponse, error)
	FindSerialsLong(ctx context.Context, in *proto.FindFilmsShortRequest, opts ...grpc.CallOption) (*proto.FindFilmsLongResponse, error)
	FindActorsShort(ctx context.Context, in *proto.FindActorsShortRequest, opts ...grpc.CallOption) (*proto.FindActorsShortResponse, error)
	FindActorsLong(ctx context.Context, in *proto.FindActorsShortRequest, opts ...grpc.CallOption) (*proto.FindActorsLongResponse, error)
	GetTopFilms(ctx context.Context, in *proto.GetTopFilmsRequest, opts ...grpc.CallOption) (*proto.GetTopFilmsResponse, error)
	GetAllFilmComments(ctx context.Context, in *proto.AllFilmCommentsRequest, opts ...grpc.CallOption) (*proto.AllFilmCommentsResponse, error)
	AddComment(ctx context.Context, in *proto.AddCommentRequest, opts ...grpc.CallOption) (*proto.AddCommentResponse, error)
	RemoveComment(ctx context.Context, in *proto.RemoveCommentRequest, opts ...grpc.CallOption) (*proto.RemoveCommentResponse, error)
}
