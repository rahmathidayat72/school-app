package handler

func GenerateDummyDataList() []ResponseUser {
	return []ResponseUser{
		{
			ID:      "1",
			Nama:    "Udin",
			Email:   "udin@mail.com",
			Telepon: "082556546928",
			Alamat:  "pantai gading",
		},
		{
			ID:      "2",
			Nama:    "Cimin",
			Email:   "cimin@mail.com",
			Telepon: "082556546928",
			Alamat:  "Guatemala",
		},
		{
			ID:      "3",
			Nama:    "Susi",
			Email:   "susi@mail.com",
			Telepon: "082556546928",
			Alamat:  "Honduras",
		},
		{
			ID:      "4",
			Nama:    "Hamdan",
			Email:   "hamdan@mail.com",
			Telepon: "082556546928",
			Alamat:  "Yogyakarta",
		},
	}
}

func GenerateDummyDataAdd() []ResponseUser {
	return []ResponseUser{
		{
			ID:      "6",
			Nama:    "Mimin",
			Email:   "mimin@mail.com",
			Telepon: "082556546928",
			Alamat:  "jakarta",
		},
	}
}

func GenerateDummyGetUsersBy() []ResponseUser {
	return []ResponseUser{
		{
			ID:      "2",
			Nama:    "Cimin",
			Email:   "cimin@mail.com",
			Telepon: "082556546928",
			Alamat:  "Guatemala",
		},
	}
}
