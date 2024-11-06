package internal

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	//creating the mux router using chi
	mux := chi.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	mux.Use(corsOptions.Handler)
	mux.Use(authenticateMiddleware)

	//For LOGIN/SIGNUP Pages
	mux.Post("/registerUser", RegisterUser)
	mux.Post("/loginUser", LoginUser)
	mux.Get("/logoutUser", LogoutUser)
	mux.Get("/authenticate", Authenticate) //to check if user is already logged in previously

	//For PROFILE Page
	mux.Get("/userProfile", UserProfile) //to send user details to frontent for profile page
	// mux.Post("/updateUserProfile", UpdateUserProfile) 	//to update the user profile details

	//For CREATE PROJECT Page
	mux.Post("/createProject", CreateProject)

	//For DASHBOARD Page
	mux.Get("/allUserProjects", AllUserProjects) //to get array of objects containing data of each project of user

	mux.Post("/teamMembers", TeamMembers)
	mux.Post("/eligibleMembersForProject", EligibleMembersForProject)
	mux.Post("/inviteUserForProject", InviteUserForProject)
	mux.Get("/userInvitations", UserInvitations)

	mux.Post("/acceptInvite", AcceptInvite)
	mux.Post("/rejectInvite", RejectInvite)

	return mux
}
