package utils

type errorMessages struct {
	SignUpError,
	LoginError,
	UserNotFound string
}

var ErrorMessages = errorMessages{
	SignUpError: "Email already exists",
	LoginError:  "Invalid email or password",

	UserNotFound: "User not found",
}

func ErrorResponse(message string) string {
	return message
}

// func ErrorResponse(message string) gin.H {
// 	return gin.H{"error": message}
// }
