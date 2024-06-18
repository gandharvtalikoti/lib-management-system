package models

type IssuedBook struct {
    ID         int    `json:"id,omitempty"`
    UserID     int    `json:"user_id"`
    BookID     int    `json:"book_id"`
    IssuedDate string `json:"issued_date"`
    DueDate    string `json:"due_date"`
    ReturnedDate string `json:"returned_date,omitempty"`
}
