package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/casbin/k8s-gatekeeper/internal/handler"
	"github.com/casbin/k8s-gatekeeper/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type ServerForTest struct {
	sync.Mutex
	running bool
	srv     *http.Server
}

func (s *ServerForTest) IsRunning() bool {
	s.Lock()
	defer s.Unlock()
	return s.running
}

func (s *ServerForTest) Stop() {
	s.srv.Shutdown(context.TODO())
}

func (s *ServerForTest) StartTestServer() {
	s.running = true
	model.IsExternalClient = true
	model.Init()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Any("/", handler.Handler)
	r.Use(func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			return
		}
		c.Next()
	})

	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := s.srv.ListenAndServeTLS("../config/certificate/server.crt", "../config/certificate/server.key")
		fmt.Println(err)
		s.Lock()
		s.running = false
		s.Unlock()
	}()
}

func RunExampleTest(workSpacePath, testCasePath string) (int, int) {

	pass := 0
	fail := 0
	//preparation stage
	exec.Command("kubectl", "apply", "-f", fmt.Sprintf("%s/model.yaml", testCasePath)).CombinedOutput()
	exec.Command("kubectl", "apply", "-f", fmt.Sprintf("%s/policy.yaml", testCasePath)).CombinedOutput()
	exec.Command("kubectl", "apply", "-f", fmt.Sprintf("%s/config/webhook_external.yaml", workSpacePath)).Run()

	server := ServerForTest{}
	server.StartTestServer()
	time.Sleep(200 * time.Millisecond)

	//run testcases
	testCaseList, _ := filepath.Glob(fmt.Sprintf("%s/testcase/*", testCasePath))
	for _, testCase := range testCaseList {
		time.Sleep(100 * time.Millisecond)
		baseName := filepath.Base(testCase)
		shouldSuccess := strings.HasPrefix(baseName, "approve")
		if server.IsRunning() {
			data, err := exec.Command("kubectl", "apply", "-f", testCase, "--dry-run=server").CombinedOutput()
			if err != nil {
				fmt.Println(err.Error())
			}
			if len(data) != 0 {
				fmt.Println(string(data))
			}

			if !server.IsRunning() {
				fail++
				fmt.Printf("[E2E Test]:FAILED Test suit %s, SERVER HAS CRASHED\n", testCase)
			} else if shouldSuccess && err != nil || !shouldSuccess && err == nil {
				fail++
				fmt.Printf("[E2E Test]:FAILED Test suit %s\n", testCase)
			} else {
				pass++
				fmt.Printf("[E2E Test]:PASSED Test suit %s\n", testCase)
			}
		} else {
			fail++
			fmt.Printf("[E2E Test]:NOTRUN Test suit %s, SERVER HAS CRASHED\n", testCase)
		}

	}

	//tear down environment
	server.Stop()
	exec.Command("kubectl", "delete", "-f", fmt.Sprintf("%s/config/webhook_external.yaml", workSpacePath)).Run()
	exec.Command("kubectl", "delete", "-f", fmt.Sprintf("%s/model.yaml", testCasePath)).Run()
	exec.Command("kubectl", "delete", "-f", fmt.Sprintf("%s/policy.yaml", testCasePath)).Run()

	return pass, fail

}
