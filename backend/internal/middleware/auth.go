package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const AuthUserKey contextKey = "auth_user"

type AuthUser struct {
	UserID   int
	Username string
	Role     string
	ClubID   *int
}

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			authUser, err := parseAuthUserFromClaims(claims)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, authUser.Username)
			ctx = context.WithValue(ctx, AuthUserKey, authUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthUser(ctx context.Context) (AuthUser, bool) {
	value := ctx.Value(AuthUserKey)
	user, ok := value.(AuthUser)
	return user, ok
}

func parseAuthUserFromClaims(claims jwt.MapClaims) (AuthUser, error) {
	username, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(username) == "" {
		return AuthUser{}, fmt.Errorf("missing subject")
	}

	role, ok := claims["role"].(string)
	if !ok || strings.TrimSpace(role) == "" {
		return AuthUser{}, fmt.Errorf("missing role")
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return AuthUser{}, fmt.Errorf("missing user_id")
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok || userIDFloat <= 0 {
		return AuthUser{}, fmt.Errorf("invalid user_id")
	}
	userID := int(userIDFloat)

	var clubID *int
	if clubRaw, exists := claims["club_id"]; exists && clubRaw != nil {
		clubFloat, ok := clubRaw.(float64)
		if !ok || clubFloat <= 0 {
			return AuthUser{}, fmt.Errorf("invalid club_id")
		}
		club := int(clubFloat)
		clubID = &club
	}

	return AuthUser{
		UserID:   userID,
		Username: username,
		Role:     strings.ToUpper(strings.TrimSpace(role)),
		ClubID:   clubID,
	}, nil
}
