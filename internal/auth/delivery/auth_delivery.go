package delivery

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	firebase "firebase.google.com/go"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/swaggo/http-swagger"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/auth/repository/db"
	"ProjectMessenger/internal/auth/usecase"
	"ProjectMessenger/internal/misc"
	profileUsecase "ProjectMessenger/internal/profile/usecase"
	contacts "ProjectMessenger/microservices/contacts_service/proto"
	session "ProjectMessenger/microservices/sessions_service/proto"
)

type AuthHandler struct {
	Sessions          session.AuthCheckerClient
	Users             usecase.UserStore
	ContactsGRPC      contacts.ContactsClient
	Firebase          *firebase.App
	prometheusMetrics *PrometheusMetrics
}

type PrometheusMetrics struct {
	ActiveSessionsCount prometheus.Gauge
	Hits                *prometheus.CounterVec
	Errors              *prometheus.CounterVec
	Methods             *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	activeSessionsCount := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_sessions_total",
			Help: "Total number of active sessions.",
		},
	)

	hits := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_hits",
			Help: "Total number of hits.",
		}, []string{"status", "path"},
	)

	errorsInProject := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_errors",
			Help: "Number of errors some type.",
		}, []string{"error_type"},
	)

	methods := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_methods",
			Help: "called methods.",
		}, []string{"method"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_http_request_duration_seconds",
			Help:    "Histogram of request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	prometheus.MustRegister(activeSessionsCount, hits, errorsInProject, methods, requestDuration)

	return &PrometheusMetrics{
		ActiveSessionsCount: activeSessionsCount,
		Hits:                hits,
		Errors:              errorsInProject,
		Methods:             methods,
		requestDuration:     requestDuration,
	}
}

func NewAuthHandler(dataBase *sql.DB, sessions session.AuthCheckerClient, avatarPath string, ContactsGRPC contacts.ContactsClient, app *firebase.App) *AuthHandler {
	handler := AuthHandler{
		Sessions:          sessions,
		Users:             db.NewUserStorage(dataBase, avatarPath),
		ContactsGRPC:      ContactsGRPC,
		prometheusMetrics: NewPrometheusMetrics(),
		Firebase:          app,
	}
	return &handler
}

func NewRawAuthHandler(dataBase *sql.DB, avatarPath string) *AuthHandler {
	handler := AuthHandler{
		//Sessions: repository.NewSessionStorage(dataBase),
		Users:             db.NewRawUserStorage(dataBase, avatarPath),
		prometheusMetrics: NewPrometheusMetrics(),
	}
	return &handler
}
func (authHandler *AuthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

// Login logs user in
//
// @Summary logs user in
// @ID login
// @Accept application/json
// @Produce application/json
// @Param user body  domain.Person true "Person"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "wrong json structure | user not found | wrong password"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /login [post]
func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	sessionHttp, err := r.Cookie("session_id")
	if !errors.Is(err, http.ErrNoCookie) {
		sessionExists, _ := usecase.CheckAuthorized(ctx, sessionHttp.Value, authHandler.Sessions)
		authHandler.prometheusMetrics.Methods.WithLabelValues("CheckAuthorized").Inc()
		if sessionExists {
			authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
			authHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
			misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "session already exists"})
			return
		}
	}
	if r.Method != http.MethodPost {
		authHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			authHandler.prometheusMetrics.Hits.WithLabelValues("500", r.URL.String()).Inc()
			return
		}
	}
	decoder := json.NewDecoder(r.Body)
	var jsonUser domain.Person
	err = decoder.Decode(&jsonUser)
	if err != nil {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
	}

	authHandler.prometheusMetrics.Methods.WithLabelValues("LoginUser").Inc()
	sessionID, err := usecase.LoginUser(ctx, jsonUser, authHandler.Users, authHandler.Sessions)

	if err != nil {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		authHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	authHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
	authHandler.prometheusMetrics.ActiveSessionsCount.Inc()
	duration := time.Since(start)
	authHandler.prometheusMetrics.requestDuration.WithLabelValues("/login").Observe(duration.Seconds())
	misc.WriteStatusJson(ctx, w, 200, nil)
}

// Logout logs user out
//
// @Summary logs user out
// @ID logout
// @Produce json
// @Success 200 {object}  domain.Response[int]
// @Failure 400 {object}  domain.Response[domain.Error] "no session to logout"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /logout [get]
func (authHandler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	session, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		authHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "no session to logout"})
		return
	}

	sessionExists, _ := usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	if !sessionExists {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		authHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "no session to logout"})
		return
	}
	authHandler.prometheusMetrics.Methods.WithLabelValues("LogoutUser").Inc()
	usecase.LogoutUser(ctx, session.Value, authHandler.Sessions)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	authHandler.prometheusMetrics.ActiveSessionsCount.Dec()
	authHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
	misc.WriteStatusJson(ctx, w, 200, nil)
	duration := time.Since(start)
	authHandler.prometheusMetrics.requestDuration.WithLabelValues("/logout").Observe(duration.Seconds())
}

// Register registers user
//
// @Summary registers user
// @ID register
// @Accept json
// @Produce json
// @Param user body  domain.Person true "Person"
// @Success 200 {object}  domain.Response[int]
// @Failure 405 {object}  domain.Response[domain.Error] "use POST"
// @Failure 400 {object}  domain.Response[domain.Error] "user already exists | required field empty | wrong json structure"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /register [post]
func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	logger := slog.With("requestID", ctx.Value("traceID"))
	if r.Method != http.MethodPost {
		authHandler.prometheusMetrics.Hits.WithLabelValues("405", r.URL.String()).Inc()
		authHandler.prometheusMetrics.Errors.WithLabelValues("405").Inc()
		misc.WriteStatusJson(ctx, w, 405, domain.Error{Error: "use POST"})
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			authHandler.prometheusMetrics.Hits.WithLabelValues("401", r.URL.String()).Inc()
			return
		}
	}

	decoder := json.NewDecoder(r.Body)
	var jsonUser domain.Person
	err := decoder.Decode(&jsonUser)
	if err != nil {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: "wrong json structure"})
	}

	authHandler.prometheusMetrics.Methods.WithLabelValues("RegisterAndLoginUser").Inc()
	sessionID, userID, err := usecase.RegisterAndLoginUser(ctx, jsonUser, authHandler.Users, authHandler.Sessions)
	if err != nil {
		authHandler.prometheusMetrics.Errors.WithLabelValues("400").Inc()
		misc.WriteStatusJson(ctx, w, 400, domain.Error{Error: err.(*domain.CustomError).Message})
		authHandler.prometheusMetrics.Hits.WithLabelValues("400", r.URL.String()).Inc()
		return
	}

	authHandler.prometheusMetrics.Methods.WithLabelValues("AddToAllContacts").Inc()
	ok := profileUsecase.AddToAllContacts(ctx, userID, authHandler.Users, authHandler.ContactsGRPC)
	if !ok {
		logger.Error("Register: contacts failed", "userID", userID)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	authHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
	misc.WriteStatusJson(ctx, w, 200, nil)
	duration := time.Since(start)
	authHandler.prometheusMetrics.requestDuration.WithLabelValues("/register").Observe(duration.Seconds())
}

// CheckAuth checks that user is authenticated
//
// @Summary checks that user is authenticated
// @ID checkAuth
// @Produce json
// @Success 200 {object}  domain.Response[int]
// @Failure 401 {object}  domain.Response[domain.Error] "Person not authorized"
// @Failure 500 {object}  domain.Response[domain.Error] "Internal server error"
// @Router /checkAuth [get]
func (authHandler *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authHandler.prometheusMetrics.Methods.WithLabelValues("CheckAuthorized").Inc()
		authorized, _ = usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	}
	if authorized {
		misc.WriteStatusJson(ctx, w, 200, nil)
	} else {
		authHandler.prometheusMetrics.Errors.WithLabelValues("401").Inc()
		misc.WriteStatusJson(ctx, w, 401, domain.Error{Error: "Person not authorized"})
	}
	duration := time.Since(start)
	authHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
	authHandler.prometheusMetrics.requestDuration.WithLabelValues("/checkAuth").Observe(duration.Seconds())
}

func (authHandler *AuthHandler) CheckAuthNonAPI(w http.ResponseWriter, r *http.Request) (authorized bool, userID uint) {
	ctx := r.Context()
	session, err := r.Cookie("session_id")
	fmt.Println(err)
	if err == nil && session != nil {
		authHandler.prometheusMetrics.Methods.WithLabelValues("CheckAuthorized").Inc()
		authorized, userID = usecase.CheckAuthorized(ctx, session.Value, authHandler.Sessions)
	}
	if !authorized {
		authHandler.prometheusMetrics.Errors.WithLabelValues("401").Inc()
		misc.WriteStatusJson(ctx, w, 401, domain.Error{Error: "Person not authorized"})
	}
	authHandler.prometheusMetrics.Hits.WithLabelValues("200", r.URL.String()).Inc()
	return authorized, userID
}
