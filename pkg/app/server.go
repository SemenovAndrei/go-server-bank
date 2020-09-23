package app

import (
	"encoding/json"
	"github.com/i-hit/go-server-bank.git/pkg/app/dto"
	"github.com/i-hit/go-server-bank.git/pkg/card"
	"log"
	"net/http"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getCards", s.getCards)
	s.mux.HandleFunc("/addCard", s.addCard)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getCards(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	if userId == "" {
		dtos := dto.CardErrDTO{Err: card.ErrUserId.Error()}
		replyErr(w, dtos)
		return
	}

	cards := s.cardSvc.AllCards()
	var dtos []*dto.CardDTO

	for _, c := range cards {
		if userId == c.UserId {
			dtos = append(
				dtos,
				&dto.CardDTO{
					Id:     c.Id,
					UserId: c.UserId,
					Number: c.Number,
					Type:   c.Type,
					System: c.System,
				})
		}
	}
	reply(w, dtos)
}

func (s *Server) addCard(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	cardType := r.FormValue("type")
	cardSystem := r.FormValue("system")

	userCard, err := s.cardSvc.Add(userId, cardType, cardSystem)

	if err != nil {
		dtos := dto.CardErrDTO{Err: err.Error()}
		replyErr(w, dtos)
		return
	}

	var dtos []*dto.CardDTO
	dtos = append(dtos,
		&dto.CardDTO{
			Id:     userCard.Id,
			UserId: userCard.UserId,
			Number: userCard.Number,
			Type:   userCard.Type,
			System: userCard.System,
		})

	reply(w, dtos)
}

func reply(w http.ResponseWriter, dtos []*dto.CardDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}

func replyErr(w http.ResponseWriter, dtos dto.CardErrDTO) {
	respBody, err := json.Marshal(dtos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
	}
}
