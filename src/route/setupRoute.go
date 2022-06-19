package route

import (
	"medias-ms/src/config"
	"net/http"
	"strings"

	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var total404Requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total_404",
		Help: "Total number of 404 requests.",
	},
	[]string{"path"},
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

var uniqueClients = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "http_unique_clients",
	Help: "Number of unique clients.",
}, []string{"ip", "timestamp", "browser"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		ip := strings.Split(r.RemoteAddr, ":")[0]
		browser := r.UserAgent()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()
		uniqueClients.WithLabelValues(ip, time.Now().Format(time.UnixDate), browser).Inc()

		if statusCode == 404 {
			total404Requests.WithLabelValues(path).Inc()
		}

		timer.ObserveDuration()
	})
}

func SetupRoutes(container config.ControllerContainer) *mux.Router {
	route := mux.NewRouter()

	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
	prometheus.Register(total404Requests)

	routerWithApiAsPrefix := route.PathPrefix("/api").Subrouter()

	routerWithApiAsPrefix.Use(prometheusMiddleware)

	routerWithApiAsPrefix.Path("/metrics").Handler(promhttp.Handler())

	routerWithApiAsPrefix.HandleFunc("/medias", container.MediaController.Upload).Methods("POST")
	routerWithApiAsPrefix.HandleFunc("/medias/{id}", container.MediaController.Delete).Methods("DELETE")
	routerWithApiAsPrefix.HandleFunc("/medias/{id}", container.MediaController.GetById).Methods("GET")

	return routerWithApiAsPrefix
}
