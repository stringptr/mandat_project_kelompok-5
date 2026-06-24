package notification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	appnats "github.com/stringptr/SiGizi/backend/internal/infrastructure/nats"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
	"github.com/nats-io/nats.go"
)

type SSEHandler struct {
	ps   *appnats.PubSub
	jwt  *jwtutils.JWT
	repo notificationDomain.Repo
}

func NewSSEHandler(ps *appnats.PubSub, jwt *jwtutils.JWT, repo notificationDomain.Repo) *SSEHandler {
	return &SSEHandler{ps: ps, jwt: jwt, repo: repo}
}

func (h *SSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, err := h.jwt.Decode(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	var mu sync.Mutex
	var subs []*nats.Subscription
	defer func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}()

	subjects := []string{fmt.Sprintf("notif.user.%d", claims.IDUser)}
	for _, role := range claims.Roles {
		subjects = append(subjects, fmt.Sprintf("notif.role.%s", role))
	}

	for _, subject := range subjects {
		sub, err := h.ps.Subscribe(subject, func(msg []byte) {
			mu.Lock()
			defer mu.Unlock()

			fmt.Fprintf(w, "event: notification\ndata: %s\n\n", msg)

			count, err := h.repo.CountUnreadByUserID(r.Context(), claims.IDUser)
			if err == nil {
				countJSON, _ := json.Marshal(map[string]int32{"unread_count": count})
				fmt.Fprintf(w, "event: unread_count\ndata: %s\n\n", countJSON)
			}
			flusher.Flush()
		})
		if err != nil {
			log.Printf("SSE subscribe %s: %v", subject, err)
			continue
		}
		subs = append(subs, sub)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	count, _ := h.repo.CountUnreadByUserID(r.Context(), claims.IDUser)
	countJSON, _ := json.Marshal(map[string]int32{"unread_count": count})
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"ok\"}\n\nevent: unread_count\ndata: %s\n\n", countJSON)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			mu.Lock()
			count, err := h.repo.CountUnreadByUserID(r.Context(), claims.IDUser)
			if err == nil {
				countJSON, _ := json.Marshal(map[string]int32{"unread_count": count})
				fmt.Fprintf(w, "event: unread_count\ndata: %s\n\n", countJSON)
				flusher.Flush()
			}
			mu.Unlock()
		}
	}
}
