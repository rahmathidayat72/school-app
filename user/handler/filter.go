package handler

import "apk-sekolah/user"

func (handler *UserHandler) getUsersByFilter(roleParam, searchParam string) ([]user.UserCore, error) {
	// Jika parameter role adalah "admin" atau "user", maka kita filter data sesuai peran
	if roleParam != "" {
		return handler.getUsersByRole(roleParam)
	}

	// Jika terdapat parameter pencarian, lakukan pencarian berdasarkan nama, email, dan alamat
	if searchParam != "" {
		return handler.searchUsers(searchParam)
	}

	// Jika tidak ada filter atau pencarian, tampilkan semua pengguna
	return handler.getAllUsers()
}

func (handler *UserHandler) getUsersByRole(roleParam string) ([]user.UserCore, error) {
	var userList []user.UserCore
	result, err := handler.userService.GetByRole(&userList, roleParam)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (handler *UserHandler) searchUsers(searchParam string) ([]user.UserCore, error) {
	var userList []user.UserCore
	result, err := handler.userService.SearchUsers(&userList, searchParam)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (handler *UserHandler) getAllUsers() ([]user.UserCore, error) {
	return handler.userService.GetAll()
}
