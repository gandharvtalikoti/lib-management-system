package utils

import (
	"log"
	"time"
)

// CalculateFine calculates the fine based on the issued date and returned date.
func CalculateFine(issuedDateStr, returnedDateStr string) (float64, error) {
    issuedDate, err := time.Parse("2006-01-02", issuedDateStr)
    if err != nil {
        return 0, err
    }
    returnedDate, err := time.Parse("2006-01-02", returnedDateStr)
    if err != nil {
        return 0, err
    }

    dueDate := issuedDate.AddDate(0, 0, 30)
    var fine float64
    if returnedDate.After(dueDate) {
        daysLate := returnedDate.Sub(dueDate).Hours() / 24
        fine = daysLate * 0.25
    } else {
        fine = 0
    }
	log.Println(fine)
    return fine, nil
}
