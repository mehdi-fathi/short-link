package Serialization

import (
	"fmt"
	"short-link/internal/Config"
	"short-link/internal/Core/Domin"
	"time"
)

// Character is one character from the database.
type LinkSerialized struct {
	*Domin.Link
	UrlShort string
	UpdatedAt string
}

func DeserializeLink(link *Domin.Link) *LinkSerialized {

	var dataLinkSerialized LinkSerialized
	dataLinkSerialized.Link = link

	// Example transformation: parsing a timestamp string to time.Time
	// Assuming link has a LastLogin field which is a string timestamp in the database
	// Check if UpdatedAt is valid before parsing and formatting
	if link.UpdatedAt.Valid {
		parsedTime := link.UpdatedAt.Time // Directly use the Time from sql.NullTime

		// Format the time to your desired format
		dataLinkSerialized.UpdatedAt = parsedTime.Format(`2006-01-02 15:04:05`) // Correct date format (YYYY-MM-DD HH:MM:SS)

		// Optionally assign the Link struct to the serialized data
		dataLinkSerialized.Link = link
	} else {
		// Handle the case when UpdatedAt is null
		fmt.Println("UpdatedAt is null")
		dataLinkSerialized.UpdatedAt = "-" // Or any default value you prefer
	}


	// Example transformation: parsing a timestamp string to time.Time
	// Assuming link has a LastLogin field which is a string timestamp in the database
	if link.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, link.CreatedAt)
		if err == nil {
			dataLinkSerialized.CreatedAt = parsedTime.Format(`2006-02-01 15:04:05`) // Assuming ParsedLastLogin is a time.Time field
		}
	}

	if link.ShortKey != ""{
		dataLinkSerialized.UrlShort = Config.GetBaseUrl() + "/short/" + link.ShortKey
	}

	return &dataLinkSerialized
	// Add other transformations as needed
}

func DeserializeAllLink(link map[int]*Domin.Link) map[int]*LinkSerialized {
	// Example transformation: parsing a timestamp string to time.Time
	// Assuming link has a LastLogin field which is a string timestamp in the database
	var dataLinkSerialized = make(map[int]*LinkSerialized)

	for key, data := range link {

		dataLinkSerialized[key] = DeserializeLink(data)
	}

	return dataLinkSerialized
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
