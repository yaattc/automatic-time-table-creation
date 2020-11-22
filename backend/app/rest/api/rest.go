package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/service"

	"github.com/go-pkgz/auth"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	log "github.com/go-pkgz/lgr"
	R "github.com/go-pkgz/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"
)

// Rest defines a simple web server for routing to calendar REST api methods
type Rest struct {
	Version string

	Authenticator *auth.Service

	httpServer *http.Server
	lock       sync.Mutex
	DataStore  *service.DataStore

	teacherRest teacherCtrlGroup
	uniRest     uniCtrlGroup
}

const hardBodyLimit = 1024 * 64 // limit size of body

// Run starts the web-server for listening
func (s *Rest) Run(port int) {
	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(port, s.routes())
	s.httpServer.ErrorLog = log.ToStdLogger(log.Default(), "WARN")
	s.lock.Unlock()

	log.Printf("[INFO] started web server at port %d", port)
	err := s.httpServer.ListenAndServe()
	log.Printf("[WARN] web server terminated reason: %s", err)
}

func (s *Rest) makeHTTPServer(port int, routes chi.Router) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

// notFound returns standard 404 not found message
func (s *Rest) notFound(w http.ResponseWriter, r *http.Request) {
	rest.SendErrorJSON(w, r, http.StatusNotFound, nil, "not found", rest.ErrBadRequest)
}

func (s *Rest) controllerGroups() (teacherCtrlGroup, uniCtrlGroup) {
	teacherGroup := teacherCtrlGroup{dataService: s.DataStore}
	uniGroup := uniCtrlGroup{dataService: s.DataStore}
	return teacherGroup, uniGroup
}

func (s *Rest) routes() chi.Router {
	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	r.Use(R.AppInfo("attc", "yaattc", s.Version))
	r.Use(R.Recoverer(log.Default()))
	r.Use(R.Ping, middleware.RealIP)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.NotFound(s.notFound)

	authHandler, _ := s.Authenticator.Handlers()
	s.teacherRest, s.uniRest = s.controllerGroups()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(5 * time.Second))
		r.Mount("/auth", authHandler)
	})

	m := s.Authenticator.Middleware()

	r.With(m.Auth).Route("/api/v1", func(rapi chi.Router) {
		rapi.Group(func(rt chi.Router) {
			rt.Post("/teacher", s.teacherRest.addTeacherCtrl)
			rt.Delete("/teacher", s.teacherRest.deleteTeacherCtrl)
			rt.Get("/teacher", s.teacherRest.listTeachersCtrl)
			rt.Post("/teacher/{id}/preferences", s.teacherRest.setTeacherPreferencesCtrl)
		})

		rapi.Group(func(rg chi.Router) {
			rg.Post("/group", s.uniRest.addGroup)
			rg.Get("/group", s.uniRest.listGroups)
			rg.Delete("/group", s.uniRest.deleteGroup)
		})

		rapi.Group(func(rsy chi.Router) {
			rsy.Post("/study_year", s.uniRest.addStudyYear)
			rsy.Get("/study_year", s.uniRest.listStudyYears)
			rsy.Delete("/study_year", s.uniRest.deleteStudyYear)
		})

		rapi.Group(func(rcrs chi.Router) {
			rcrs.Post("/course", s.uniRest.addCourse)
		})

		rapi.Group(func(rsched chi.Router) {
			rsched.Get("/time_slot", s.uniRest.listTimeSlots)
		})
	})

	return r
}
