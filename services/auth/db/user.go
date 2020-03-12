package db

const (

	existUserByUsernameAndPasswordStatement = "SELECT COUNT(1) AS count FROM tbl_user WHERE user_name = ? AND user_pwd = ?"
)

func ExistUserByUsernameAndPassword(username, password string) bool {
	return exist(existUserByUsernameAndPasswordStatement, username, password)
}
