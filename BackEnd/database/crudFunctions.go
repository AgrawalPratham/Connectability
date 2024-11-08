package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUserDatabase func to store the user signup details into database
func RegisterUserDatabase(conn *sql.DB, user *UserDetails) error {
	//converting the password to encrypted form
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing the passowrd : ", err)
		return err
	}

	if user.Image == "" {
		user.Image = ""
	}
	
	//Inserting user data into database
	query := `Insert into user_details(name, email, password, phone, city, state, country, about, github_id, linkedin_id, resume, experience, image) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = conn.Exec(query, user.Name, user.Email, passwordHash, user.Phone, user.City, user.State, user.Country, user.About, user.Github_id, user.Linkedin_id, user.Resume, user.Experience, user.Image)
	if err != nil {
		fmt.Println("Error inserting the user registration details into database : ", err)
		return err
	}

	//Inserting user skills into database
	for _, skill := range user.Skills {
		query := "INSERT INTO user_skills(email, skill) VALUES ($1, $2)"
		_, err := conn.Exec(query, user.Email, skill)
		if err != nil {
			fmt.Printf("Error inserting skill %s: %v", skill, err)
			return err
		}
	}
	return nil
}

// LoginUserDatabase func authenticates the user drom the database
func LoginUserDatabase(conn *sql.DB, user *UserDetails) error {
	var passwordHash []byte

	query := `Select password from user_details where email = $1`
	row := conn.QueryRow(query, user.Email)
	err := row.Scan(&passwordHash)
	if err != nil {
		fmt.Println("Error extracting the password for user authentication : ", err)
		return err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(user.Password))
	if err != nil {
		fmt.Println("Login password does not matched : ", err)
		return err
	}

	fmt.Println("User successfully authenticated")
	return nil
}

// UserProfileDatabase func to extract the details of user from the database
func UserProfileDatabase(conn *sql.DB, userEmail string) (*UserDetails, error) {
	var userDetail UserDetails

	//Extracting user personal details from user_details table
	query1 := `Select * from user_details where email = $1`
	row := conn.QueryRow(query1, userEmail)
	err := row.Scan(
		&userDetail.Name,
		&userDetail.Email,
		&userDetail.Password,
		&userDetail.Phone,
		&userDetail.City,
		&userDetail.State,
		&userDetail.Country,
		&userDetail.About,
		&userDetail.Github_id,
		&userDetail.Linkedin_id,
		&userDetail.Resume,
		&userDetail.Experience,
		&userDetail.Image,
	)
	if err != nil {
		fmt.Println("Error extracting user_details : ", err)
		return nil, err
	}
	userDetail.Password = "" //to fide the password

	//Extracting user skills from user_skills table
	query2 := `Select skill from user_skills where email = $1`
	rows, err := conn.Query(query2, userEmail)
	if err != nil {
		fmt.Println("Error extracting user skills : ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var skill string
		err := rows.Scan(&skill)
		if err != nil {
			fmt.Println("Error parsing the rows of user skills : ", err)
			return nil, err
		}
		userDetail.Skills = append(userDetail.Skills, skill)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("error scanning the rows of user skills : ", err)
		return nil, err
	}

	return &userDetail, nil
}

// CreateProjectDatabase func to insert the project creation details in project_detail table and project_skills, project_images table in data present
func CreateProjectDatabase(conn *sql.DB, projectDetail *ProjectDetails) error {

	//Converting start date from string to date format
	startDate, err := time.Parse("2006-01-02", projectDetail.Start_date)
	if err != nil {
		fmt.Println("Error parsing start date from string to date format")
		return err
	}

	//Converting end date from string to date format
	endDate, err := time.Parse("2006-01-02", projectDetail.End_date)
	if err != nil {
		fmt.Println("Error parsing end date from string tom date format")
		return err
	}

	//Deriving the value of project id for the new project
	projectDetail.Project_id = 0
	row := conn.QueryRow("Select max(project_id) from project_details")
	err = row.Scan(&projectDetail.Project_id)
	if err != nil {
	}
	projectDetail.Project_id = projectDetail.Project_id + 1

	//Inserting the project details into project_details table
	query := `Insert into project_details(project_id,project_name, project_description, github_link, start_date, end_date, budget, chat_link, manager_email) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = conn.Exec(
		query,
		projectDetail.Project_id,
		projectDetail.Project_name,
		projectDetail.Project_description,
		projectDetail.Github_link,
		startDate,
		endDate,
		projectDetail.Budget,
		projectDetail.Chat_link,
		projectDetail.Manager_email)
	if err != nil {
		fmt.Println("Error inserting project details into database table project_details : ", err)
		return err
	}

	//Inserting project images into project_images table
	for _, image := range projectDetail.Images {
		query := "INSERT INTO project_images(project_id, image) VALUES ($1, $2)"
		_, err := conn.Exec(query, projectDetail.Project_id, image)
		if err != nil {
			fmt.Printf("Error inserting project image %s: %v", image, err)
			return err
		}
	}

	// //Inserting project skills into project_skills table
	for _, skill := range projectDetail.Skills {
		query := "INSERT INTO project_skills(project_id, skill) VALUES ($1, $2)"
		_, err := conn.Exec(query, projectDetail.Project_id, skill)
		if err != nil {
			fmt.Printf("Error inserting project skill %s: %v", skill, err)
			return err
		}
	}

	return nil
}


// AllUserProjectsDatabase func to return an array of objects of project details where user is working of is manager
func AllUserProjectsDatabase(conn *sql.DB, userEmail string) (*[]ProjectDetails, error) {

	//Extracting an array of Project ids of all the projects where user is working or is manager
	Project_ids, err := AllUserProjectsIdDatabase(conn, userEmail)
	if err != nil {
		fmt.Println("Error finding the project ids of user projects : ", err)
		return nil, err
	}

	//Array of all the projects where user is working or is manager
	var Projects []ProjectDetails

	// Iterating over each project ID and retrieve the project details
	for _, Project_id := range Project_ids {
		Project, err := GetProjectDetailsByIDDatabase(conn, Project_id)
		if err != nil {
			fmt.Println("Error fetching project details for project ID:", Project_id, err)
			return nil, err
		}
		Project.Project_id = Project_id
		Projects = append(Projects, *Project)
	}

	return &Projects, nil

}

	// AllUserProjectIdDatabase return array of project id of projects where user is working of is manager
	func AllUserProjectsIdDatabase(conn *sql.DB, userEmail string) ([]int, error) {
		var Project_ids []int

		query := `Select project_id from project_details where manager_email = $1`
		rows, err := conn.Query(query, userEmail)
		if err != nil {
			fmt.Println("Error extracting project id of projects where user in project manager : ", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var Project_id int
			err := rows.Scan(&Project_id)
			if err != nil {
				fmt.Println("Error parsing the rows of project id of projects where user in project manager : ", err)
				return nil, err
			}
			Project_ids = append(Project_ids, Project_id)
		}

		query = `Select project_id from team_table where member_email = $1`
		rows, err = conn.Query(query, userEmail)
		if err != nil {
			fmt.Println("Error extracting project id of projects where user in working : ", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var Project_id int
			err := rows.Scan(&Project_id)
			if err != nil {
				fmt.Println("Error parsing the rows of project id of projects where user in working : ", err)
				return nil, err
			}
			Project_ids = append(Project_ids, Project_id)
		}

		return Project_ids, nil
	}

	// GetProjectDetailsByIDDatabase return project details of particular project id
	func GetProjectDetailsByIDDatabase(conn *sql.DB, projectID int) (*ProjectDetails, error) {
		//variable to store all the details of particular project specified by projectId
		var project = &ProjectDetails{}

		//Extracting project details from project_details table
		query := `SELECT project_name, project_description, github_link, start_date, end_date, budget, chat_link, manager_email, completed FROM project_details WHERE project_id = $1`
		err := conn.QueryRow(query, projectID).Scan(
			&project.Project_name,
			&project.Project_description,
			&project.Github_link,
			&project.Start_date,
			&project.End_date,
			&project.Budget,
			&project.Chat_link,
			&project.Manager_email,
			&project.Completed,
		)
		if err != nil {
			fmt.Println("Error extracting project details of project id : ", projectID, err)
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, project.Start_date)
		if err != nil {
			fmt.Println("Error parsing start date :", err)
			return nil, err
		}
		formattedStartDate := parsedTime.Format("2006-01-02")
		project.Start_date = formattedStartDate

		parsedTime, err = time.Parse(time.RFC3339, project.End_date)
		if err != nil {
			fmt.Println("Error parsing start date :", err)
			return nil, err
		}
		formattedEndDate := parsedTime.Format("2006-01-02")
		project.End_date = formattedEndDate

		//Extracting project skills from project_skills table
		var skills []string
		query = `Select skill from project_skills where project_id = $1`
		rows, err := conn.Query(query, projectID)
		if err != nil {
			fmt.Println("Error extracting skills of project id : ", projectID, err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var skill string
			err := rows.Scan(&skill)
			if err != nil {
				fmt.Println("Error parshing skills of project id : ", projectID, err)
				return nil, err
			}
			skills = append(skills, skill)
		}
		project.Skills = skills

		//Extracting project images from project_images table
		var images []string
		query = `Select image from project_images where project_id = $1`
		rows, err = conn.Query(query, projectID)
		if err != nil {
			fmt.Println("Error extracting images of project id : ", projectID, err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var image string
			err := rows.Scan(&image)
			if err != nil {
				fmt.Println("Error parshing images of project id : ", projectID, err)
				return nil, err
			}
			images = append(images, image)
		}
		project.Images = images

		return project, nil
	}


// TeamMembersDatabase return array of names of team members
func TeamMembersDatabase(conn *sql.DB, Project_id int) (*[]MemberNameImage, error) {
	query := `Select u.name, u.image from team_table t JOIN user_details u ON t.member_email = u.email where t.project_id = $1`
	rows, err := conn.Query(query, Project_id)
	if err != nil {
		fmt.Println("Error extracting the team members name for project id : ", Project_id, err)
		return nil, err
	}

	var MembersData []MemberNameImage
	for rows.Next() {
		var MemberData MemberNameImage
		err := rows.Scan(&MemberData.Member_name, &MemberData.Member_image)
		if err != nil {
			fmt.Println("Error parsing the names of team members for project Id : ", Project_id, err)
			return nil, err
		}
		MembersData = append(MembersData, MemberData)
	}
	return &MembersData, nil
}

func TeamHeadDatabase(conn *sql.DB, Project_id int)(*MemberNameImage, error){
	query := `Select name, image from user_details where email = ( Select manager_email from project_details where project_id = $1 )`
	row := conn.QueryRow(query, Project_id)

	var ProjectHeadData MemberNameImage
	err := row.Scan( &ProjectHeadData.Member_name, &ProjectHeadData.Member_image)
	if err!= nil{
		fmt.Println("Error extracting the project head data for project id : ", Project_id, err)
		return nil, err
	}
	return &ProjectHeadData, nil
}

//EligibleMembersForProjectDatabase return details of all the users having the required skills for the project
func EligibleMembersForProjectDatabase(conn *sql.DB, requestBody RequestBody2) ( *[]UserDetails, error){

	query := `Select Distinct email 
				From user_skills 
				Where LOWER(skill) in (`+ formatSkillsForQuery(requestBody.Skills) +`)
				And email NOT In ( Select member_email from team_table where project_id = $1)`
	rows, err := conn.Query(query, requestBody.Project_id)
	if err!= nil{
		fmt.Println("Error extracting the email of users with required skills for project id : ", requestBody.Project_id, err)
		return nil, err
	}
	var Users_email []string
	for rows.Next(){
		var User_email string
		err :=rows.Scan(&User_email)
		if err != nil{
			fmt.Println("Error parsing the email of users with required skills for project id : ", requestBody.Project_id, err)
			return nil, err
		}
		Users_email = append(Users_email, User_email)
	}

	var Users_detail []UserDetails
	for _, User_email := range Users_email{
		User_detail, err := UserProfileDatabase(conn, User_email)
		if err != nil{
			fmt.Println("Error extracting detail of user with required skill : ", err)
			return nil,err
		}
		Users_detail = append(Users_detail, *User_detail)
	}

	return &Users_detail, nil
}

	//to format the skill array in a string
	func formatSkillsForQuery(skills []string) string {
		var formattedSkills string
		for i, skill := range skills {
			formattedSkills += "LOWER('" + skill + "')"
			if i < len(skills)-1 {
				formattedSkills += ", "
			}
		}
		return formattedSkills
	}

//InviteUserForProjectDatabase to insert the request record in database
func InviteUserForProjectDatabase(conn *sql.DB, requestBody3 RequestBody3) error {

	//Checking if the user is already requested
	var requestId int
	row := conn.QueryRow("Select request_id from request_table where project_id = $1 and receiver_email = $2", requestBody3.Project_id, requestBody3.Receiver_email)
	err := row.Scan(&requestId)
	if err == nil{
		fmt.Println("User already requested. Cannot request again, Error : ", err)
		return errors.New("User already requested. Cannot request again")
	}

	//Finding the request id for the new request
	requestId = 0
	row = conn.QueryRow("Select max(request_id) from request_table")
	err = row.Scan(&requestId)
	if err!= nil{
	}
	requestId++

	//Inserting the request record in database
	query := `Insert into request_table(request_id, project_id, receiver_email, closed, accepted, rejected) values ($1,$2, $3, $4, $5, $6)`
	_, err = conn.Exec(query, requestId, requestBody3.Project_id, requestBody3.Receiver_email, false, false, false)
	if err != nil{
		fmt.Println("Error inserting the request record in request_table : ", err)
		return err
	}

	return nil
}

//UserInvitationsDatabase to return the details of user invitations from the database
func UserInvitationsDatabase(conn *sql.DB, userEmail string) (*UserInvitationsDetail, error){
	//variable to store the invitation details of the user
	var userInvitationDetail UserInvitationsDetail
	userInvitationDetail.User_email = userEmail

	row := conn.QueryRow("Select name, image from user_details where email = $1", userEmail)
	err := row.Scan(&userInvitationDetail.User_name, &userInvitationDetail.User_image)
	if err != nil{
		fmt.Println("Error extracting the invitation details for the user : ", err)
		return nil,err
	}

	//extraction project ids of projects who sended invitation to user
	var projectIds []int
	rows, err := conn.Query("Select project_id from request_table where receiver_email = $1 and closed = false and accepted = false and rejected = false", userEmail)
	if err!= nil{
		fmt.Println("Error extracting the project ids of invitation details for the user : ", err)
		return nil,err
	}
	for rows.Next(){
		var projectId int
		err := rows.Scan(&projectId)
		if err!= nil{
			fmt.Println("Error parshing the project ids of invitation details for the user : ", err)
			return nil,err
		}
		projectIds = append(projectIds, projectId)
	}

	//extracting details of each project which invited the user
	for _, projectId :=range projectIds{
		Project_detail, err := GetProjectDetailsByIDDatabase(conn, projectId)
		if( err != nil){
			fmt.Println("Error extracting project details of user invitations with project id : ", projectId, err)
			return nil, err
		}
		Project_detail.Project_id = projectId
		userInvitationDetail.Project_details = append(userInvitationDetail.Project_details, *Project_detail)
	}

	userInvitationDetail.Invitation_count = len(userInvitationDetail.Project_details)

	return &userInvitationDetail, nil
}

func AcceptInviteDatabase(conn *sql.DB, requestBody3 RequestBody3) error{
	query := `Update request_table 
				Set accepted = true
				Where project_id = $1 and receiver_email = $2`
	_, err := conn.Exec(query, requestBody3.Project_id, requestBody3.Receiver_email)
	if err!= nil{
		fmt.Println("Error updating the accepted column for the accepted request : ", err)
		return err
	}

	query = `Insert into team_table(project_id, member_email) values($1, $2)`
	_, err = conn.Exec(query, requestBody3.Project_id, requestBody3.Receiver_email)
	if err != nil{
		fmt.Println("Error inserting the new member record in team table : ", err)
		return err
	}

	return nil
}


func RejectInviteDatabase(conn *sql.DB, requestBody3 RequestBody3) error{
	query := `Update request_table 
				Set rejected = true
				Where project_id = $1 and receiver_email = $2`
	_, err := conn.Exec(query, requestBody3.Project_id, requestBody3.Receiver_email)
	if err!= nil{
		fmt.Println("Error updating the rejected column for the reject request : ", err)
		return err
	}
	
	return nil
}