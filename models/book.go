package models

type Book struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Author    string `json:"author"`
    ISBN      string `json:"isbn"`
    Stock     int    `json:"stock"`
    Available int    `json:"available"`
}
