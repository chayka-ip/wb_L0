package apiserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"level_zero/internal/app/cache"
	"level_zero/internal/app/logger"
	models "level_zero/internal/app/model"
	"level_zero/order"
	"level_zero/store"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

var (
	routeOrderCreate   = "/api/add_order"
	routeListOrders    = "/api/orders"
	routeGetOrderByUid = "/api/order"
	routeViewOrder     = "/view_order"
)

type server struct {
	router   *mux.Router
	logger   logger.Logger
	store    store.Store
	cache    cache.Cache
	stanInfo *NatsInfo
}

func NewServer(store store.Store, cache cache.Cache, logger logger.Logger, stanInfo *NatsInfo) *server {
	s := &server{
		router:   mux.NewRouter(),
		logger:   logger,
		store:    store,
		cache:    cache,
		stanInfo: stanInfo,
	}

	s.configureRouter()
	s.configureNats()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc(routeOrderCreate, s.handleOrderCreate()).Methods("POST")
	s.router.HandleFunc(routeListOrders, s.handleListOrders()).Methods("GET")
	s.router.HandleFunc(routeGetOrderByUid, s.handleGetOrderById()).Methods("GET")
	s.router.HandleFunc(routeViewOrder, s.handleViewOrderById()).Methods("GET")
}

func (s *server) configureNats() {
	clusterId := s.stanInfo.ClusterId
	clientId := s.stanInfo.ClientId
	sstan := fmt.Sprintf("cluster_id: %s | client_id: %s", clusterId, clientId)
	failConnMsg := fmt.Sprintf("Could not connect to STAN with settings: %s\n", sstan)

	conn, err := stan.Connect(clusterId, clientId)
	if err != nil {
		s.logger.Error(failConnMsg)
		return
	}

	_, err = conn.Subscribe(clusterId, s.natsHandlerOrderReceive, stan.DurableName(clientId))
	if err != nil {
		s.logger.Error(failConnMsg)
		return
	}

	fmt.Printf("Listening STAN on %s\n", sstan)
}

func (s *server) natsHandlerOrderReceive(m *stan.Msg) {
	// s.logger.Info("Message recieved from STAN\n")
	req := &order.Order{}
	err := json.Unmarshal(m.Data, req)
	if err != nil {
		s.logger.Info("Can't parse data from STAN request")
		return
	}
	err = s.createNewOrder(req)
	if err != nil {
		s.logger.Info(err)
	}
}

func (s *server) handleOrderCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &order.Order{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.createNewOrder(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusCreated, req)
	}
}

func (s *server) createNewOrder(o *order.Order) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	ord := &models.Order{
		OrderUID: o.OrderUID,
		Data:     data,
	}

	if err := s.store.Order().Create(ord); err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("Added Order with UID: %v", o.OrderUID))
	s.cache.AddToStore(ord.OrderUID, ord.Data)

	return nil
}

// for testing purposes
func (s *server) handleListOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count := 100
		orderList, err := s.store.Order().ListOrders(count)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var out []*order.Order
		for i := range orderList {
			ord := &order.Order{}
			err := json.Unmarshal(orderList[i].Data, ord)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			out = append(out, ord)
		}

		s.respond(w, r, http.StatusOK, out)
	}
}

func (s *server) handleGetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getOrderIdFromRequest(r)
		data, err := s.getOrderDataById(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, data)
	}
}

func (s *server) handleViewOrderById() http.HandlerFunc {
	type orderInfo struct {
		IsValid  bool
		OrderUID string
		Data     order.Order
	}

	return func(w http.ResponseWriter, r *http.Request) {
		info := orderInfo{IsValid: false}
		id := getOrderIdFromRequest(r)
		data, err := s.getOrderDataById(id)
		if err == nil {
			var ord order.Order
			err = json.Unmarshal(data, &ord)
			if err == nil {
				info.IsValid = true
				info.OrderUID = id
				info.Data = ord
			}
		}
		if err != nil {
			s.logger.Info("Can't parse order for view\n", err)
		}

		t_path := getTemplatePath("view_order.html")
		t := template.Must(template.ParseFiles(t_path))

		t.Execute(w, info)
	}
}

func getOrderIdFromRequest(r *http.Request) string {
	return r.URL.Query().Get("id")
}

func (s *server) getOrderDataById(id string) (json.RawMessage, error) {
	return s.cache.Get(id)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func getTemplatePath(name string) string {
	workDir, _ := os.Getwd()
	return filepath.Join(workDir, "internal", "app", "web", "static", "templates", name)
}
