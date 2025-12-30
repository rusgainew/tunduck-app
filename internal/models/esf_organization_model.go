package models

type EsfOrganizationModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Token       string `json:"token"`
	DBName      string `json:"dbName"`
}
