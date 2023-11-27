package serialization

import (
	"short-link/internal/Core/Domin"
	"time"
)

func DeserializeLink(link *Domin.Link) {
	// Example transformation: parsing a timestamp string to time.Time
	// Assuming link has a LastLogin field which is a string timestamp in the database
	if link.UpdatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, link.UpdatedAt)
		if err == nil {
			link.UpdatedAt = parsedTime.Format(`2006-02-01 15:04:05`) // Assuming ParsedLastLogin is a time.Time field
		}
	}

	// Add other transformations as needed
}

func DeserializeAllLink(link map[int]*Domin.Link) {
	// Example transformation: parsing a timestamp string to time.Time
	// Assuming link has a LastLogin field which is a string timestamp in the database

	for _, data := range link {

		if data.UpdatedAt != "" {
			parsedTime, err := time.Parse(time.RFC3339, data.UpdatedAt)
			if err == nil {
				data.UpdatedAt = parsedTime.Format(`2006-02-01 15:04:05`) // Assuming ParsedLastLogin is a time.Time field
			}
		}
	}

	// Add other transformations as needed
}

func isMap(i interface{}) (bool, map[interface{}]interface{}) {
	// The type of val will be map[interface{}]interface{} if i is a map
	val, ok := i.(map[interface{}]interface{})
	return ok, val
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
