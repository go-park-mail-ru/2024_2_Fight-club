package router

import (
	ads "2024_2_FIGHT-CLUB/internal/ads/controller"
	auth "2024_2_FIGHT-CLUB/internal/auth/controller"
	chat "2024_2_FIGHT-CLUB/internal/chat/controller"
	city "2024_2_FIGHT-CLUB/internal/cities/controller"
	review "2024_2_FIGHT-CLUB/internal/reviews/contoller"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRoutes(authHandler *auth.AuthHandler, adsHandler *ads.AdHandler, cityHandler *city.CityHandler, chatHandler *chat.ChatHandler, reviewHandler *review.ReviewHandler) *mux.Router {
	router := mux.NewRouter()
	api := "/api"

	// User Authentication Routes
	router.HandleFunc(api+"/auth/register", authHandler.RegisterUser).Methods("POST") // Register a new user
	router.HandleFunc(api+"/auth/login", authHandler.LoginUser).Methods("POST")       // Login user
	router.HandleFunc(api+"/auth/logout", authHandler.LogoutUser).Methods("DELETE")   // Logout user
	// User Management Routes
	router.HandleFunc(api+"/users", authHandler.PutUser).Methods("PUT")                            // Update user
	router.HandleFunc(api+"/users/{userId}", authHandler.GetUserById).Methods("GET")               // Get user by ID
	router.HandleFunc(api+"/users", authHandler.GetAllUsers).Methods("GET")                        // Get all users
	router.HandleFunc(api+"/session", authHandler.GetSessionData).Methods("GET")                   // Get session data
	router.HandleFunc(api+"/users/{userId}/housing", adsHandler.GetUserPlaces).Methods("GET")      // Get User Ads
	router.HandleFunc(api+"/users/{userId}/favorites", adsHandler.GetUserFavorites).Methods("GET") // Get User Favorites
	// Ad Management Routes
	router.HandleFunc(api+"/housing", adsHandler.GetAllPlaces).Methods("GET")                             // Get all ads
	router.HandleFunc(api+"/housing/{adId}", adsHandler.GetOnePlace).Methods("GET")                       // Get ad by ID
	router.HandleFunc(api+"/housing", adsHandler.CreatePlace).Methods("POST")                             // Create a new ad
	router.HandleFunc(api+"/housing/{adId}", adsHandler.UpdatePlace).Methods("PUT")                       // Update ad by ID
	router.HandleFunc(api+"/housing/{adId}", adsHandler.DeletePlace).Methods("DELETE")                    // Delete ad by ID
	router.HandleFunc(api+"/housing/cities/{city}", adsHandler.GetPlacesPerCity).Methods("GET")           // Get ads by city
	router.HandleFunc(api+"/housing/{adId}/images/{imageId}", adsHandler.DeleteAdImage).Methods("DELETE") // Delete image from ad
	router.HandleFunc(api+"/housing/{adId}/like", adsHandler.AddToFavorites).Methods("POST")              // Add to favorites
	router.HandleFunc(api+"/housing/{adId}/dislike", adsHandler.DeleteFromFavorites).Methods("POST")      // Delete from favorites
	// CSRF Token Route
	router.HandleFunc(api+"/csrf/refresh", authHandler.RefreshCsrfToken).Methods("GET") // Refresh CSRF token
	// City Management Routes
	router.HandleFunc(api+"/cities", cityHandler.GetCities).Methods("GET")         // Get All Cities
	router.HandleFunc(api+"/cities/{city}", cityHandler.GetOneCity).Methods("GET") // Get One City
	// Chat Management Routes
	router.HandleFunc(api+"/messages/chats", chatHandler.GetAllChats).Methods("GET") //Get All Chats
	router.HandleFunc(api+"/messages/chat/{id}", chatHandler.GetChat).Methods("GET") //Get One Chats
	router.HandleFunc(api+"/messages/setconn", chatHandler.SetConnection)            //Set connection
	// Reviews Management Routese
	router.HandleFunc(api+"/reviews", reviewHandler.CreateReview).Methods("POST")
	router.HandleFunc(api+"/reviews/{userId}", reviewHandler.GetUserReviews).Methods("GET")
	router.HandleFunc(api+"/reviews/{hostId}", reviewHandler.DeleteReview).Methods("DELETE")
	router.HandleFunc(api+"/reviews/{hostId}", reviewHandler.UpdateReview).Methods("PUT")

	router.Handle(api+"/metrics", promhttp.Handler())

	return router
}
