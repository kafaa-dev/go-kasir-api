package database

import "kasir-api/models"

var Products = []models.Product{
	{ID: 1, Name: "Mie rebus", Price: 3500, Stock: 10},
	{ID: 2, Name: "Air minum 800ml", Price: 3000, Stock: 40},
	{ID: 3, Name: "Kecap", Price: 12000, Stock: 20},
}

var Categories = []models.Category{
	{ID: 1, Name: "Makanan", Description: "Apapun yang bisa dan aman untuk dimakan."},
	{ID: 2, Name: "Minuman", Description: "Apapun yang bisa dan aman untuk diminum."},
	{ID: 3, Name: "Bahan penyedap", Description: "Sesuatu yang dicampur ke makanan untuk memberikan rasa sedap."},
}
