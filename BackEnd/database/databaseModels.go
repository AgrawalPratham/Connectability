package database

type UserDetails struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Phone       string   `json:"phone"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	Country     string   `json:"country"`
	About       string   `json:"about"`
	Github_id   string   `json:"github_id"`
	Linkedin_id string   `json:"linkedin_id"`
	Resume      string   `json:"resume"`
	Experience  string   `json:"experience"`
	Skills      []string `json:"skills"`
	Image       string   `json:"image"`
}

type ProjectDetails struct {
	Project_id          int      `json:"project_id"`
	Project_name        string   `json:"project_name"`
	Project_description string   `json:"project_description"`
	Github_link         string   `json:"github_link"`
	Start_date          string   `json:"start_date"`
	End_date            string   `json:"end_date"`
	Budget              int      `json:"budget"`
	Chat_link           string   `json:"chat_link"`
	Manager_email       string   `json:"manager_email"`
	Completed           bool     `json:"completed"`
	Images              []string `json:"images"`
	Skills              []string `json:"skills"`
}

type ProjectTeam struct {
	Project_id   int      `json:"project_id"`
	Project_head_name string `json:"project_head_name"`
	Project_head_image string `json:"project_head_image"`
	Members_data []MemberNameImage `json:"members_data"`
}

type MemberNameImage struct {
	Member_name string `json:"member_name"`
	Member_image string `json:"member_image"`
}

//used when accepting skills of a project in json to extract and return the eligible users
type RequestBody2 struct {
	Project_id int      `json:"project_id"`
	Skills     []string `json:"skills"`
}

//struct used when project manager invites a user
type RequestBody3 struct {
	Project_id     int    `json:"project_id"`
	Receiver_email string `json:"receiver_email"`
}

//struct to store the details of invitations received by the user
type UserInvitationsDetail struct {
	User_name        string           `json:"user_name"`
	User_email       string           `json:"user_email"`
	User_image       string           `json:"user_image"`
	Invitation_count int              `json:"invitation_count"`
	Project_details  []ProjectDetails `json:"project_details"`
}
