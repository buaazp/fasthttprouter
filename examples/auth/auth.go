package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// basicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
func basicAuth(ctx *fasthttp.RequestCtx) (username, password string, ok bool) {
	auth := ctx.Request.Header.Peek("Authorization")
	if auth == nil {
		return
	}
	return parseBasicAuth(string(auth))
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// BasicAuth is the basic auth handler
func BasicAuth(h fasthttp.RequestHandler, requiredUser, requiredPassword string) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := basicAuth(ctx)

		// WARNING: This is insecure and for demonstration purposes only.
		// DO NOT use plain-text passwords for real apps.
		// This comparison == is also vulnerable to a timing attack.
		// 
		// One easy way to compare passwords safely is hash passwords with
		// scrypt and use the CompareHashAndPassword function
		// which works in constant time from
		// https://github.com/elithrar/simple-scrypt/blob/master/scrypt.go:
		
		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(ctx)
			return
		}
		// Request Basic Authentication otherwise
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	})
}

// Index is the index handler
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Not protected!\n")
}

// Protected is the Protected handler
func Protected(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Protected!\n")
}

func main() {
	user := "gordon"
	pass := "secret!"

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/protected/", BasicAuth(Protected, user, pass))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
