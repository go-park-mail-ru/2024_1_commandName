package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"ProjectMessenger/domain"
	"ProjectMessenger/internal/misc"
	"github.com/prometheus/client_golang/prometheus"
)

type Sessions struct {
	db                *sql.DB
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
	sessionErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_errors",
			Help: "Number of errors some type.",
		}, []string{"error_type"},
	)

	sessionMethods := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_called_methods_count",
			Help: "Number of called methods.",
		}, []string{"method"},
	)

	sessionMethodDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "session_method_duration_seconds",
			Help:    "Histogram of session methods durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(sessionErrors, sessionMethods, sessionMethodDuration)

	return &PrometheusMetrics{
		Errors:          sessionErrors,
		Methods:         sessionMethods,
		requestDuration: sessionMethodDuration,
	}
}

func (s *Sessions) GetUserIDbySessionID(ctx context.Context, sessionID string) (userID uint, sessionExists bool) {
	s.prometheusMetrics.Methods.WithLabelValues("GetUserIDbySessionID").Inc()
	start := time.Now()
	logger := slog.With("requestID", ctx.Value("traceID"))
	var userIDInt int
	var sid string
	err := s.db.QueryRowContext(ctx, "SELECT userid, sessionid FROM auth.session WHERE sessionid = $1", sessionID).Scan(&userIDInt, &sid)
	userID = uint(userIDInt)
	logger.Debug("GetUserIDbySessionID", "userID", userID, "sessionID", sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("didn't found user by session", "userID", userID, "sessionID", sessionID)
			return 0, false
		}
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method GetUserIDbySessionID, sessions.go",
		}
		fmt.Println(customErr.Error())
		logger.Error(customErr.Error())
		s.prometheusMetrics.Errors.WithLabelValues("database").Inc()
		return 0, false
	}
	fmt.Println("found user by session userID", userID, "sessionID", sessionID)
	logger.Debug("found user by session", "userID", userID, "sessionID", sessionID)
	duration := time.Since(start)
	s.prometheusMetrics.requestDuration.WithLabelValues("GetUserIDbySessionID").Observe(duration.Seconds())
	return userID, true
}

func (s *Sessions) CreateSession(ctx context.Context, userID uint) (sessionID string) {
	s.prometheusMetrics.Methods.WithLabelValues("CreateSession").Inc()
	start := time.Now()
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("CreateSession", "userID", userID)
	sessionID = misc.RandStringRunes(32)
	_, err := s.db.ExecContext(ctx, "INSERT INTO auth.session (sessionid, userid) VALUES ($1, $2)", sessionID, userID)
	if err != nil {
		s.prometheusMetrics.Errors.WithLabelValues("database").Inc()
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method CreateSession, sessions.go",
		}
		fmt.Println(customErr.Error())
		logger.Error(customErr.Error())
		return ""
	}
	logger.Info("created session", "sessionID", sessionID, "userID", userID)
	duration := time.Since(start)
	s.prometheusMetrics.requestDuration.WithLabelValues("CreateSession").Observe(duration.Seconds())
	return sessionID
}

func (s *Sessions) DeleteSession(ctx context.Context, sessionID string) {
	s.prometheusMetrics.Methods.WithLabelValues("DeleteSession").Inc()
	start := time.Now()
	logger := slog.With("requestID", ctx.Value("traceID"))
	logger.Debug("DeleteSession", "sessionID", sessionID)
	_, err := s.db.ExecContext(ctx, "DELETE FROM auth.session WHERE sessionID = $1", sessionID)
	if err != nil {
		s.prometheusMetrics.Errors.WithLabelValues("database").Inc()
		customErr := &domain.CustomError{
			Type:    "database",
			Message: err.Error(),
			Segment: "method DeleteSession, sessions.go",
		}
		fmt.Println(customErr.Error())
		logger.Error(customErr.Error())
	}
	logger.Info("deleted session", "sessionID", sessionID)
	duration := time.Since(start)
	s.prometheusMetrics.requestDuration.WithLabelValues("DeleteSession").Observe(duration.Seconds())
}

func NewSessionStorage(db *sql.DB) *Sessions {
	slog.Info("created session storage")
	return &Sessions{
		db:                db,
		prometheusMetrics: NewPrometheusMetrics(),
	}
}
