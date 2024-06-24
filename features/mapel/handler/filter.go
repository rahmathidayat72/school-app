package handler

import "apk-sekolah/features/mapel"

func (handler *MapelHandler) getMapelByFilter(guruIDParam, searchParam string) ([]mapel.MapelCore, error) {
	// Jika parameter guruID tidak kosong, maka lakukan filter berdasarkan GuruID
	if guruIDParam != "" {
		return handler.getMapelByGuruID(guruIDParam)
	}

	// Jika terdapat parameter pencarian, lakukan pencarian berdasarkan mapel
	if searchParam != "" {
		return handler.searchMapel(searchParam)
	}

	// Jika tidak ada filter atau pencarian, tampilkan semua mapel
	return handler.getAllMapel()
}

func (handler *MapelHandler) getMapelByGuruID(guruIDParam string) ([]mapel.MapelCore, error) {
	var mapelList []mapel.MapelCore
	result, err := handler.mapelService.GetByGuruID(&mapelList, guruIDParam)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (handler *MapelHandler) searchMapel(searchParam string) ([]mapel.MapelCore, error) {
	var mapelList []mapel.MapelCore

	result, err := handler.mapelService.SearchMapel(&mapelList, searchParam)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (handler *MapelHandler) getAllMapel() ([]mapel.MapelCore, error) {
	return handler.mapelService.GetAll()
}
