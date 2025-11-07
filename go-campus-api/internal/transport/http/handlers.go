package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-campus-api/internal/domain"
	"go-campus-api/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) http.Handler {
	h := &Handler{service: service}

	mux := http.NewServeMux()

	mux.Handle("/get_peer_status/", withMiddleware(http.HandlerFunc(h.GetPeerStatus), methodValidator("GET")))
	mux.Handle("/update_peers/", withMiddleware(http.HandlerFunc(h.UpdatePeers), methodValidator("POST"), jsonValidator()))
	mux.Handle("/get_friends_status/", withMiddleware(http.HandlerFunc(h.GetFriendsStatus), methodValidator("GET")))
	mux.Handle("/add_friend/", withMiddleware(http.HandlerFunc(h.AddFriend), methodValidator("POST"), jsonValidator()))
	mux.Handle("/delete_friend/", withMiddleware(http.HandlerFunc(h.DeleteFriend), methodValidator("POST"), jsonValidator()))

	return mux
}

func (h *Handler) GetPeerStatus(w http.ResponseWriter, r *http.Request) {
	peerName := r.URL.Query().Get("peer_name")
	if peerName == "" {
		respondError(w, http.StatusBadRequest, "peer_name is required")
		return
	}

	peer, err := h.service.GetPeerStatus(r.Context(), peerName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if peer == nil {
		respondError(w, http.StatusNotFound, "peer not found")
		return
	}

	respondJSON(w, http.StatusOK, peer)
}

func (h *Handler) UpdatePeers(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdatePeersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdatePeers(r.Context(), req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "peers updated successfully"})
}

func (h *Handler) GetFriendsStatus(w http.ResponseWriter, r *http.Request) {
	tgIDStr := r.URL.Query().Get("tg_id")
	if tgIDStr == "" {
		respondError(w, http.StatusBadRequest, "tg_id is required")
		return
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid tg_id")
		return
	}

	friends, err := h.service.GetFriendsStatus(r.Context(), tgID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, friends)
}

func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) {
	var req domain.AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.AddFriend(r.Context(), req.TgID, req.PeerName); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "friend added successfully"})
}

func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	var req domain.DeleteFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.DeleteFriend(r.Context(), req.TgID, req.PeerName); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "friend deleted successfully"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
