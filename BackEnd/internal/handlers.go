package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/agrawalpratham/Connectability/BackEnd/config"
	"github.com/agrawalpratham/Connectability/BackEnd/database"
)

// RegisterUser function stores the details of new user in database
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//variable to store extracted registration form data
	var user = &database.UserDetails{}

	//decoding registration form data from json format in struct of user details
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Error parsing the registration form : %v", err)
		return
	}

	//Calling function to insert user registration details into database
	err = database.RegisterUserDatabase(config.App.DBConnection.SQL, user)
	if err != nil {
		fmt.Println("Error in calling UserRegistrationDatabase func to insert user registration details in database : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not successfully registered"))
		return
	}

	fmt.Println("User successfully registered")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully registered"))
}

// LoginUser checks the authentication of user during login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	//variable to store extracted login form data
	var user = &database.UserDetails{}

	//decoding login form data from json format in struct of user details
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintf(w, "Error parsing the registration form : %v", err)
		return
	}

	//calling LoginUserDatabase func to authenticate user
	err = database.LoginUserDatabase(config.App.DBConnection.SQL, user)
	if err != nil {
		fmt.Println("Error while authentication the user : ", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authentication failed: Unauthorized access"))
		return
	}

	session, _ := config.App.Session.Get(r, "Connectability")
	session.Values["userEmail"] = user.Email
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}else {
    log.Println("Session saved successfully")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User successfully logged in"))
}

// LogoutUser log the user out
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := config.App.Session.Get(r, "Connectability")
	delete(session.Values, "userEmail")
	err := session.Save(r, w)
	// session.Options.MaxAge = -1

	if err != nil {
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	// config.App.UserEmail = ""
	w.WriteHeader(http.StatusOK)
}

// Authenticate to check is user is logged in previously
func Authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := config.App.Session.Get(r, "Connectability")
	userEmail := session.Values["userEmail"].(string)
	if userEmail != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// UserProfile to return the user data when user comes to their profile page
func UserProfile(w http.ResponseWriter, r *http.Request) {
	session, err := config.App.Session.Get(r, "Connectability")
	if err != nil {
		log.Println(err.Error())
	} else if session.IsNew {
		log.Println("session is new")
	}
	log.Println(session.Values)
	userEmail, ok := session.Values["userEmail"].(string)
	if( !ok){
		log.Println("user id is missing")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("user id is missing"))
		return
	}
	userDetails, err := database.UserProfileDatabase(config.App.DBConnection.SQL, userEmail)
	if err != nil {
		fmt.Println("Error while calling database func to extract user details : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while extracting user details"))
		return
	}

	jsonUserDetails, err := json.Marshal(userDetails)
	if err != nil {
		fmt.Println("Error marshalling the user details into json : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling user details into json format"))
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserDetails)
}

// CreateProject to store the details of new project in database
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var projectDetail = &database.ProjectDetails{}

	//Unmarshalling the project details data
	err := json.NewDecoder(r.Body).Decode(projectDetail)
	if err != nil {
		fmt.Println("Error decoding project details from json format : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding project details from json format"))
		return
	}

	//Setting the project manager email to the current user email
	session, _ := config.App.Session.Get(r, "Connectability")
	userEmail := session.Values["userEmail"].(string)
	projectDetail.Manager_email = userEmail

	//Calling database func to insert project creation details into database
	err = database.CreateProjectDatabase(config.App.DBConnection.SQL, projectDetail)
	if err != nil {
		fmt.Println("Error calling database  func to insert project details into database : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database  func to insert project details into database"))
		return
	}

	fmt.Println("Project creation details successfully inserted in database")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project creation details successfully inserted in database"))
}

// AllUserProjects to return and array of json objects of each project where user si working or is manager
func AllUserProjects(w http.ResponseWriter, r *http.Request) {

	session, _ := config.App.Session.Get(r, "Connectability")
	userEmail := session.Values["userEmail"].(string)

	//Calling database func to get an array of all user projects
	Projects, err := database.AllUserProjectsDatabase(config.App.DBConnection.SQL, userEmail)
	if err != nil {
		fmt.Println("Error retrieving all user projects data : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error retrieving all user projects data"))
		return
	}

	//Marshalling User projects data into json format
	jsonProjects, err := json.Marshal(Projects)
	if err != nil {
		fmt.Println("Error marshalling user projects data into json format")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling user projects data into json format"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonProjects)
}

// TeamMembers to return a json object containing project id and an array of names of team members
func TeamMembers(w http.ResponseWriter, r *http.Request) {

	type RequestBody struct {
		Project_id int `json:"project_id"`
	}
	var requestBody RequestBody

	//Unmarshalling the project id from json received from frontend
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println("Error unmarshalling the project id from frontend json : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Extracting the team members name and image from the database
	Team_members, err := database.TeamMembersDatabase(config.App.DBConnection.SQL, requestBody.Project_id)
	if err != nil {
		fmt.Println("Error calling database function to extract team members name and image for project id : ", requestBody.Project_id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database function to extract team members name and image "))
		return
	}

	//Extracting the team head name and image from the database
	Team_head, err := database.TeamHeadDatabase(config.App.DBConnection.SQL, requestBody.Project_id)
	if err != nil {
		fmt.Println("Error calling database function to extract team head name and image  for project id : ", requestBody.Project_id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database function to extract team haed and image  name"))
		return
	}

	//Marshalling the team members name to json
	var Project_team = &database.ProjectTeam{}
	Project_team.Project_id = requestBody.Project_id
	Project_team.Project_head_name = Team_head.Member_name
	Project_team.Project_head_image = Team_head.Member_image
	Project_team.Members_data = *Team_members

	jsonProject_team, err := json.Marshal(Project_team)
	if err != nil {
		fmt.Println("Error converting Project Team data to json  format for project id : ", requestBody.Project_id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error converting Project Team data to json  format"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonProject_team)
}

//EligibleMembersForProject to return the details of users having required skills for the project
func EligibleMembersForProject(w http.ResponseWriter, r *http.Request){
	
	var requestBody database.RequestBody2

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err!= nil{
		fmt.Println("Error unmarshalling the skills needed to filter users for project : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("rror unmarshalling the skills needed to filter users for project"))
		return
	}

	Users_detail, err := database.EligibleMembersForProjectDatabase(config.App.DBConnection.SQL, requestBody)
	if err!= nil{
		fmt.Println("Error extracting details of members with required skill for project id : ", requestBody.Project_id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error extracting details of members with required skill for project"))
		return 
	}

	jsonUsers_detail, err := json.Marshal(Users_detail)
	if err!= nil{
		fmt.Println("Error marshalling the details of Users with required skill")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling the details of Users with required skill"))
		return 
	}

	fmt.Println("DONE")
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUsers_detail)
}

//InviteUserForProject to insert the invitation request record in database
func InviteUserForProject(w http.ResponseWriter, r *http.Request){

	var requestBody3 database.RequestBody3

	//Unmarshalling the project id and receiever email of the invitation 
	err := json.NewDecoder(r.Body).Decode(&requestBody3)
	if err!= nil{
		fmt.Println("Error unmarshalling the project id and receiver email of invitation : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error unmarshalling the project id and receiver email of invitation"))
		return
	}

	//Inserting the invitation record in database
	err = database.InviteUserForProjectDatabase(config.App.DBConnection.SQL, requestBody3)
	if err!= nil{
		fmt.Println("Error calling the database function to insert the invitation request record : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling the database function to insert the invitation request record"))
		return
	}

	fmt.Println("Invitation request record successfully inserted in database")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Invitation request record successfully inserted in database"))
}

//UserInvitations to send the details of user invitation projects 
func UserInvitations(w http.ResponseWriter, r *http.Request){

	session, _ := config.App.Session.Get(r, "Connectability")
	userEmail := session.Values["userEmail"].(string)

	//Extracting details of user invitations
	userInvitationDetail, err := database.UserInvitationsDatabase(config.App.DBConnection.SQL, userEmail)
	if err!= nil{
		fmt.Println("Error calling database func to extract details of invitations received by user : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database func to extract details of invitations received by user"))
		return
	}

	//Marshalling details of user invitations 
	jsonUserInvitationDetail, err := json.Marshal(userInvitationDetail)
	if err!= nil{
		fmt.Println("Error marshalling the user invitation details into json format : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling the user invitation details into json format"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUserInvitationDetail)
}

func AcceptInvite(w http.ResponseWriter, r *http.Request){
	var requestBody3 database.RequestBody3
	err := json.NewDecoder(r.Body).Decode(&requestBody3)
	if err != nil{
		fmt.Println("Error unmarshalling the project id and receiver email of the accepted request : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error unmarshalling the project id and receiver email of the accepted request"))
		return
	}

	session, _ := config.App.Session.Get(r, "Connectability")
	userEmail := session.Values["userEmail"].(string)
	requestBody3.Receiver_email = userEmail

	err = database.AcceptInviteDatabase(config.App.DBConnection.SQL, requestBody3)
	if err!= nil{
		fmt.Println("Error calling database func to update recordes for accepted request : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database func to update recordes for accepted request"))
		return
	}

	fmt.Printf("Request successfully accepted for project id : %d by user emai : %s \n", requestBody3.Project_id, requestBody3.Receiver_email)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request successfully accepted"))
}

func RejectInvite(w http.ResponseWriter, r *http.Request){
	var requestBody3 database.RequestBody3
	err := json.NewDecoder(r.Body).Decode(&requestBody3)
	if err != nil{
		fmt.Println("Error unmarshalling the project id and receiver email of the rejected request : ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error unmarshalling the project id and receiver email of the rejected request"))
		return
	}
	session, _ := config.App.Session.Get(r, "Connectability")
	requestBody3.Receiver_email = session.Values["userEmail"].(string)

	err = database.RejectInviteDatabase(config.App.DBConnection.SQL, requestBody3)
	if err!= nil{
		fmt.Println("Error calling database func to update recordes for rejected request : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error calling database func to update recordes for rejected request"))
		return
	}

	fmt.Printf("Request successfully rejected for project id : %d by user emai : %s \n", requestBody3.Project_id, requestBody3.Receiver_email)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request successfully rejected"))
}