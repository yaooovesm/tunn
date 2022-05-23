package config

//
// Admin
// @Description:
//
type Admin struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}
