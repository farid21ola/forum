package graph

import (
	"context"
	"github.com/farid21ola/forum/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

const (
	userloaderKey = "userloader"
)

func DataloaderMiddleware(db *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := UserLoader{
			maxBatch: 100,
			wait:     1 * time.Millisecond,
			fetch: func(ids []string) ([]*model.User, []error) {
				var users []*model.User
				query := `SELECT * FROM users WHERE id = ANY($1)`
				rows, err := db.Query(context.Background(), query, ids)
				if err != nil {
					return nil, []error{err}
				}
				defer rows.Close()

				for rows.Next() {
					var user model.User
					err = rows.Scan(&user.ID, &user.Username)
					if err != nil {
						return nil, []error{err}
					}
					users = append(users, &user)
				}

				if rows.Err() != nil {
					return nil, []error{err}
				}

				u := make(map[string]*model.User, len(users))
				for _, user := range users {
					u[user.ID] = user
				}

				result := make([]*model.User, len(ids))
				for i, id := range ids {
					result[i] = u[id]
				}
				return result, nil
			},
		}

		ctx := context.WithValue(r.Context(), userloaderKey, &userLoader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value(userloaderKey).(*UserLoader)
}
