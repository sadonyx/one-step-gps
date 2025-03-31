package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"github.com/sadonyx/assessmentAPI/internal/session"
	"github.com/sadonyx/assessmentAPI/internal/tokens"
)

const SessionCookieName = session.SessionCookieName

// Creates a new router manager.
func NewRouter(sm *session.SessionManager, tm *tokens.TokenManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/user-preferences", Middleware(sm)(preferencesHandler(sm, tm)))
	mux.HandleFunc("/", eventsHandler(sm, tm))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(mux)
}

// Returns handler for sending live device events using Server-side Events.
func eventsHandler(sm *session.SessionManager, tm *tokens.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context()

		shortLivedToken := r.URL.Query().Get("slt")
		if len(shortLivedToken) == 0 {
			log.Println("No slt in query parameters.")
			http.Error(w, "Bad Request: missing slt", http.StatusBadRequest)
			return
		}

		var frequency float64
		token := tm.Get(shortLivedToken)
		if token == nil {
			log.Println("Invalid/expired token.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token.SessionID == "" {
			log.Println("Token has no session ID.")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userPreferences, err := sm.GetSession(ctx, token.SessionID)
		if err != nil {
			frequency = 5
		} else {
			// maintain default -- necessary if user does not initially have cookie
			if userPreferences.Preferences.PollingFrequency > 0 {
				frequency = userPreferences.Preferences.PollingFrequency
			}
		}

		pollingDuration := time.Duration(frequency) * time.Second

		rc := http.NewResponseController(w)
		ticker := time.NewTicker(pollingDuration)
		defer ticker.Stop()

		sseHelper(w, rc)

		clientGone := r.Context().Done()

		for {
			select {
			case <-clientGone:
				fmt.Println("Client disconnected")
				ticker.Stop()
				return
			case <-ticker.C:
				sseHelper(w, rc)
			}
		}
	}
}

// Returns handler for updating user preferences.
func preferencesHandler(sm *session.SessionManager, tm *tokens.TokenManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sessionID := session.GetSessionFromContext(r.Context())
		userPreferences := session.GetSessionDataFromContext(r.Context())

		switch r.Method {
		case http.MethodGet:
			tokenValue := tm.Create(sessionID, 30*time.Second)

			w.Header().Set("Authorization", tokenValue)

			jsonStr, err := json.Marshal(userPreferences)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonStr)
		case http.MethodPost:
			var updatedPreferences UpdatedPreferencesRequest
			err := json.NewDecoder(r.Body).Decode(&updatedPreferences)
			updatedPreferences.Visits = userPreferences.Visits + 1

			mapped := mapRequestToPreferences(updatedPreferences)

			updatedDocument := sm.UpdateSession(r.Context(), sessionID, mapped)
			if err != nil {
				http.Error(w, "Failed to update session", http.StatusInternalServerError)
				return
			}

			if updatedDocument != nil {

				jsonStr, err := json.Marshal(updatedDocument["preferences"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Write(jsonStr)
			}
		}
	}
}

// Helper function that requests OSGPS API and returns parsed/structured data.
func sseHelper(w http.ResponseWriter, rc *http.ResponseController) {
	oneStepGpsApiKey := os.Getenv("ONE_STEP_GPS_API_KEY")
	if oneStepGpsApiKey == "" {
		http.Error(w, "Missing API key", http.StatusInternalServerError)
		return
	}

	getUrl := fmt.Sprintf("https://track.onestepgps.com/v3/api/public/device?latest_point=true&api-key=%s", oneStepGpsApiKey)
	res, err := http.Get(getUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error making API request: %v", err), http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("API returned non-200 status code: %d", res.StatusCode), res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	log.Println("Sent OSG API data.")

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	if err := rc.Flush(); err != nil {
		http.Error(w, "Error flushing response", http.StatusInternalServerError)
	}
}

// Cookie authorization middleware.
func Middleware(sm *session.SessionManager) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			ctx := r.Context()

			// get session ID from cookie
			cookie, err := r.Cookie(SessionCookieName)
			fmt.Printf("Cookie: %v\n", cookie)

			if err == http.ErrNoCookie || cookie == nil {
				// session does not exists, create one
				println("Creating new session")
				sessionID, err := sm.CreateSession(ctx)
				if err != nil {
					println("failed create")
					http.Error(w, "Failed to create session", http.StatusInternalServerError)
					return
				}

				session, err := sm.GetSession(ctx, sessionID)
				if err != nil {
					println("failed get")
					http.Error(w, "Failed to get session details", http.StatusInternalServerError)
					return
				}

				setCookie(sessionID, w, r)

				// store session in context
				ctx = context.WithValue(ctx, SessionCookieName, sessionID)
				ctx = context.WithValue(ctx, "preferences", session.Preferences)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// session exists, validate it
			sessionID := cookie.Value
			session, err := sm.GetSession(ctx, sessionID)

			if err != nil {
				// invalid/expired session, create a new one
				sessionID, err := sm.CreateSession(ctx)
				if err != nil {
					http.Error(w, "Failed to create session", http.StatusInternalServerError)
					return
				}

				session, err := sm.GetSession(ctx, sessionID)
				if err != nil {
					http.Error(w, "Failed to get session details", http.StatusInternalServerError)
					return
				}

				setCookie(sessionID, w, r)

				ctx = context.WithValue(ctx, SessionCookieName, sessionID)
				ctx = context.WithValue(ctx, "preferences", session.Preferences)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			println("Session exists and is valid")
			// valid session, store in context
			ctx = context.WithValue(ctx, SessionCookieName, sessionID)
			ctx = context.WithValue(ctx, "preferences", session.Preferences)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper function to set cookie to response.
func setCookie(sID string, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sID,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(session.SessionExpirationTime.Seconds()),
	})
}

// Helper function to convert request object to `Preferences` struct to be updated in MongoDB.
func mapRequestToPreferences(req UpdatedPreferencesRequest) session.Preferences {
	return session.Preferences{
		SortOrder:              req.SortOrder,
		HiddenDevices:          req.HiddenDevices,
		Visits:                 req.Visits,
		ShowVisibilityControls: req.ShowVisibilityControls,
		PollingFrequency:       req.PollingFrequency,
	}
}

// -------- Originally used geocode, however it is costly! --------

// func getGeocode(lat float64, lng float64) string {
// 	s := fmt.Sprintf("<API-Key>", lat, lng)

// 	res, resErr := http.Get(s)
// 	if resErr != nil {
// 		fmt.Printf("Error making API request: %v", resErr)
// 		return ""
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		fmt.Printf("API returned non-200 status code: %d", res.StatusCode)
// 		return "No specified address"
// 	}

// 	body, errBody := io.ReadAll(res.Body)
// 	if errBody != nil {
// 		fmt.Printf("Error reading response body: %v", errBody)
// 		return ""
// 	}

// 	// fmt.Println(string(body))

// 	var response struct {
// 		Results []struct {
// 			FormattedAddress string `json:"formatted_address"`
// 		} `json:"results"`
// 		Address string `json:"formatted_address"`
// 	}
// 	err := json.Unmarshal(body, &response)
// 	if err != nil {
// 		return ""
// 	}

// 	if len(response.Results) > 0 {
// 		return response.Results[0].FormattedAddress
// 	} else {
// 		return "No specified address"
// 	}
// }
