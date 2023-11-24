package serialization

import (
	"short-link/internal/Db/Model"
	"time"
)

func DeserializeLink(user *Model.Link) {
	// Example transformation: parsing a timestamp string to time.Time
	// Assuming user has a LastLogin field which is a string timestamp in the database
	if user.UpdatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, user.UpdatedAt)
		if err == nil {
			user.UpdatedAt = parsedTime.Format(`2006-02-01 15:04:05`) // Assuming ParsedLastLogin is a time.Time field
		}
	}

	// Add other transformations as needed
}

//
//func SerializeUser(user *models.User) {
//	// Example transformation: hashing a password (if password field exists)
//	if user.Password != "" {
//		hashedPassword, err := HashPassword(user.Password)
//		if err == nil {
//			user.Password = hashedPassword
//		}
//	}
//
//	// Add other transformations as needed
//}
//
//// HashPassword is an example function for hashing passwords
//func HashPassword(password string) (string, error) {
//	// Use a proper hashing algorithm (e.g., bcrypt)
//	// return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//}
