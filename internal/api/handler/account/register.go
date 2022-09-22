// Copyright 2021 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package account

import (
	"net/http"
	"time"

	"github.com/harness/gitness/internal/api/render"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/internal/token"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/check"

	"github.com/dchest/uniuri"
	"github.com/rs/zerolog/hlog"
	"golang.org/x/crypto/bcrypt"
)

// HandleRegister returns an http.HandlerFunc that processes an http.Request
// to register the named user account with the system.
func HandleRegister(users store.UserStore, system store.SystemStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := hlog.FromRequest(r)

		username := r.FormValue("username")
		password := r.FormValue("password")

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Err(err).
				Str("email", username).
				Msg("Failed to hash password")

			render.InternalError(w)
			return
		}

		// TODO: allow to provide email and name separately ...
		user := &types.User{
			Name:     username,
			Email:    username,
			Password: string(hash),
			Salt:     uniuri.NewLen(uniuri.UUIDLen),
			Created:  time.Now().UnixMilli(),
			Updated:  time.Now().UnixMilli(),
		}

		if err = check.User(user); err != nil {
			log.Debug().Err(err).
				Str("email", username).
				Msg("invalid user input")

			render.UserfiedErrorOrInternal(w, err)
			return
		}

		if err = users.Create(ctx, user); err != nil {
			log.Err(err).
				Str("email", username).
				Msg("Failed to create user")

			render.InternalError(w)
			return
		}

		// if the registered user is the first user of the system,
		// assume they are the system administrator and grant the
		// user system admin access.
		if user.ID == 1 {
			user.Admin = true
			if err = users.Update(ctx, user); err != nil {
				log.Err(err).
					Str("email", username).
					Int64("user_id", user.ID).
					Msg("Failed to enable admin user")

				render.InternalError(w)
				return
			}
		}

		expires := time.Now().Add(system.Config(ctx).Token.Expire)
		token, err := token.GenerateExp(user, expires.Unix(), user.Salt)
		if err != nil {
			log.Err(err).
				Str("email", username).
				Int64("user_id", user.ID).
				Msg("Failed to generate token")

			render.InternalError(w)
			return
		}

		// return the token if the with_user boolean
		// query parameter is set to true.
		if r.FormValue("return_user") == "true" {
			render.JSON(w, http.StatusOK, &types.UserToken{
				User: user,
				Token: &types.Token{
					Value:   token,
					Expires: expires.UTC(),
				},
			})
		} else {
			// else return the token only.
			render.JSON(w, http.StatusOK, &types.Token{Value: token})
		}
	}
}