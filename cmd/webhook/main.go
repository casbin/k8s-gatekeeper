// Copyright 2022 The Casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/casbin/k8s-gatekeeper/internal/handler"
	"github.com/casbin/k8s-gatekeeper/internal/model"
)

func tlsHandler(c *gin.Context) {
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
	})
	err := secureMiddleware.Process(c.Writer, c.Request)
	// If there was an error, do not continue.
	if err != nil {
		return
	}
	c.Next()
}

func main() {
	name := flag.Bool("externalClient", true, "is running inside the k8s cluster")
	flag.Parse()
	model.IsExternalClient = *name

	model.Init()

	r := gin.Default()
	r.Any("/", handler.Handler)
	r.Use(tlsHandler)
	r.RunTLS(":8080", "config/certificate/server.crt", "config/certificate/server.key")
}
