package main

import (
	"fmt"
	"go_project/database"
	"go_project/models"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	database.ConnectDB()
}

type postBooksBody struct {
	Books_data []models.BooksInfo `json:"books_data" binding:"required"`
}

type borrowBook struct {
	BookId int    `json:"book_id" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}

func HashPassword(password string) (string, error) {
	//Use 14 as cost instead of 10 in case of more unique result as mentioned in the blog.
	//in below function.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("Error occurred while comparing hashed password", err)
	return err == nil
}

func main() {
	//For deleting the table use database.DB.Migrator().DropTable(&models.<TableName>{})

	server := gin.Default()
	server.Use(cors.Default())
	// server.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"POST", "PATCH", "PUT", "GET"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Access-Control-Allow-Origin"},
	// 	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	// 	AllowCredentials: true,
	// 	AllowAllOrigins:  false,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://github.com"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))

	router := server.Group("/api")

	//Migration of new column
	// if !database.DB.Migrator().HasColumn(&models.UserInfo{}, "Password") {
	// 	database.DB.Migrator().AddColumn(&models.UserInfo{}, "Password")
	// }
	database.DB.AutoMigrate(&models.UserInfo{})
	fmt.Println("Migration is succesfully completed.")

	//UserInfo API's
	//Post request to create new users API based on given payload user details
	router.POST("/create-new-users", func(ctx *gin.Context) {
		var userPayloadData struct {
			UserData []struct {
				UserId    string `json:"UserId" binding:"required"`
				Name      string `json:"Name" binding:"required"`
				Email     string `json:"Email" binding:"required"`
				Password  string `json:"Password" binding:"required"`
				ContactNo int64  `json:"ContactNo" binding:"required"`
				Category  string `json:"Category" binding:"required"`
			} `json:"user_data" binding:"required,dive"`
		}

		err := ctx.ShouldBindJSON(&userPayloadData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_data field or one of user input fields are not present in payload details."})
			return
		}
		// fmt.Println("Printing user_data: ", userPayloadData.UserData)
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(userPayloadData.UserData); i++ {
				var currRec models.UserInfo
				currRecRes := tx.Where("email = ?", userPayloadData.UserData[i].Email).Find(&currRec)
				if currRecRes.RowsAffected != 0 {
					return io.EOF
				}
				var newUser models.UserInfo
				// newUser := userPayloadData.UserData[i]
				newUser.Name = userPayloadData.UserData[i].Name
				newUser.Email = userPayloadData.UserData[i].Email
				newUser.Category = userPayloadData.UserData[i].Category
				newUser.UserId = userPayloadData.UserData[i].UserId
				newUser.ContactNo = userPayloadData.UserData[i].ContactNo
				var hashError error
				newUser.Password, hashError = HashPassword(userPayloadData.UserData[i].Password)
				if hashError != nil {
					fmt.Printf("Error occurred for curr record: %v while hashing the password, Err: %v\n", newUser, hashError)
					return hashError
				}
				res := tx.Create(&newUser)
				if res.Error != nil {
					return res.Error
				}
			}
			return nil
		})
		if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "One of user data already exists in database."})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Transaction failure while adding new user data."})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": "success", "msg": "Created given new user data."})
		}
	})

	//Update API for userinfo details based on given userID.
	router.PUT("/update-user", func(ctx *gin.Context) {
		var userPayload struct {
			UserId    string `json:"user_id" binding:"required"`
			Name      string `json:"name" binding:"required"`
			Email     string `json:"email" binding:"required"`
			Password  string `json:"password" binding:"required"`
			ContactNo int64  `json:"contact_no" binding:"required"`
			Category  string `json:"category" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&userPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Expected payload fields are not present!"})
			return
		}
		var currRec models.UserInfo
		currRecRes := database.DB.Where("user_id = ?", userPayload.UserId).Find(&currRec)
		if currRecRes.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "Given UserID doesn't exist!"})
			return
		}

		currRec.Name = userPayload.Name
		currRec.ContactNo = userPayload.ContactNo
		currRec.Password, _ = HashPassword(userPayload.Password)
		currRec.Category = userPayload.Category
		var emailRec models.UserInfo
		emailCheckRes := database.DB.Where("email = ?", userPayload.Email).Find(&emailRec)
		if emailCheckRes.RowsAffected > 0 && emailRec.UserId != currRec.UserId {
			ctx.JSON(http.StatusNotModified, gin.H{"status": "failure", "msg": "New email for given user already exists for different user!"})
			return
		}
		updateRes := database.DB.Save(&currRec)
		if updateRes.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "failure", "msg": "Update failed for given userID!"})
		} else {
			ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "msg": "Upated the user details successfully!"})
		}
	})

	//Check if the user exists in the database based on given userId
	router.POST("/check-user-existence/user-id", func(ctx *gin.Context) {
		var userPayload struct {
			UserId   string `json:"user_id" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&userPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either user_id/password fields are not present in payload."})
			return
		}
		var currUserRec models.UserInfo
		queryRes := database.DB.Where("user_id = ?", userPayload.UserId).Find(&currUserRec)
		if queryRes.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "UserId doesn't exist!"})
		} else {
			if CheckPasswordHash(userPayload.Password, currUserRec.Password) {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Given user data exists in database."})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Password doesn't match!"})
			}
		}
	})

	//Check if the user exists based on given email.
	router.POST("/check-user-existence/email", func(ctx *gin.Context) {
		var userPayload struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&userPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either email/password fields are not present in payload."})
			return
		}
		var currUserRec models.UserInfo
		queryRes := database.DB.Where("email = ?", userPayload.Email).Find(&currUserRec)
		if queryRes.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "Email doesn't exist!"})
		} else {
			if CheckPasswordHash(userPayload.Password, currUserRec.Password) {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Given user data exists in database."})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Password doesn't match!"})
			}
		}
	})

	//All users API
	router.GET("/all-userinfo", func(ctx *gin.Context) {
		var user_list []models.UserInfo
		result := database.DB.Find(&user_list)
		fmt.Printf("Printing Result row count: %v", result.RowsAffected)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": user_list})
	})

	//Fetch userinfo using userId API
	router.GET("/userinfo/from-id/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var userIdDataList []models.UserInfo
		result := database.DB.Where("User_Id = ?", id).Find(&userIdDataList)
		fmt.Println("Printing row for given user Id: ", result.RowsAffected)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userIdDataList})
	})

	//Fetch userinfo using email API
	router.GET("/userinfo/from-email", func(ctx *gin.Context) {
		email := ctx.Query("email")
		var userEmailDataList []models.UserInfo
		result := database.DB.Where("email = ?", email).Find(&userEmailDataList)
		fmt.Println("Printing row for given user Email: ", result.RowsAffected)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userEmailDataList})
	})

	//Book's API's
	database.DB.AutoMigrate(&models.BooksInfo{})
	//Populate the data in book_info table
	//POST API for populating the new data in books_info table
	router.POST("/new-books-data", func(ctx *gin.Context) {
		var booksList postBooksBody
		err := ctx.ShouldBindJSON(&booksList)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "failure", "errors": "books_data field is require with value = array of objects"})
			return
		}
		fmt.Println("Printing books list: ", booksList.Books_data)
		for i := 0; i < len(booksList.Books_data); i++ {
			var existingRec models.BooksInfo
			errorRes := database.DB.Where("title = ? AND author = ? AND publications = ?",
				booksList.Books_data[i].Title, booksList.Books_data[i].Author, booksList.Books_data[i].Publications).First(&existingRec).Error
			fmt.Println("Printing Existing Rec : ", errorRes, existingRec)
			if errorRes == nil {
				// if errorRes != gorm.ErrRecordNotFound { ---> Another way of writing
				booksList.Books_data[i].BookId = existingRec.BookId
				database.DB.Save(&booksList.Books_data[i])
			} else {
				database.DB.Save(&booksList.Books_data[i])
			}
		}
		ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
	})

	//all-booksinfo API
	router.GET("/all-booksinfo", func(ctx *gin.Context) {
		var getAllBooks []models.BooksInfo
		result := database.DB.Find(&getAllBooks)
		fmt.Println("Printing fetched booksinfo: ", result.RowsAffected)
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": getAllBooks})
	})

	//fetch bookinfo from book_id API
	router.GET("/booksinfo/from-book-id/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var bookInfoList []models.BooksInfo
		result := database.DB.Where("book_id = ?", id).Find(&bookInfoList)
		fmt.Println("For given Id, Fetched records count: ", result.RowsAffected)
		if len(bookInfoList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": bookInfoList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Invalid Book Id"})
		}
	})

	//fetch bookinfo from Title, Author, Publications API
	router.GET("/booksinfo/from-book-data", func(ctx *gin.Context) {
		title := ctx.Query("title")
		author := ctx.Query("author")
		publications := ctx.Query("publications")

		var booksInfoList []models.BooksInfo
		result := database.DB.Where("title = ? AND author = ? AND publications = ?", title, author, publications).Find(&booksInfoList)
		fmt.Println("Printing fetched rows count: ", result.RowsAffected)
		if len(booksInfoList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksInfoList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book data found for given input."})
		}
	})

	//fetch bookinfo from Title and Author API
	router.GET("/booksinfo/from-title-author", func(ctx *gin.Context) {
		title := ctx.Query("title")
		author := ctx.Query("author")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("title = ? AND author = ?", title, author).Find(&booksDataList)
		fmt.Println("Printing fetched rows count: ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "success", "msg": "No book data found for given title and author."})
		}
	})

	//fetch bookinfo from Title and Publications API
	router.GET("/booksinfo/from-title-publication", func(ctx *gin.Context) {
		title := ctx.Query("title")
		publications := ctx.Query("publications")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("title = ? AND publications = ?", title, publications).Find(&booksDataList)
		fmt.Println("Printing fetched rows count: ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "success", "msg": "No book data found for given title and publication."})
		}
	})

	//fetch booksinfo from Author and Publications API
	router.GET("/booksinfo/from-author-publication", func(ctx *gin.Context) {
		author := ctx.Query("author")
		publications := ctx.Query("publications")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("author = ? AND publications = ?", author, publications).Find(&booksDataList)
		fmt.Println("Printing fetched rows count : ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book data found for given author and publications."})
		}
	})

	//fetch all the books based on given author name API
	router.GET("/booksinfo/from-author", func(ctx *gin.Context) {
		author := ctx.Query("author")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("author = ?", author).Find(&booksDataList)

		fmt.Println("Printing fetched rows count: ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book data found for given author."})
		}
	})

	//fetch all the books based on given publication name API
	router.GET("/booksinfo/from-publication", func(ctx *gin.Context) {
		publications := ctx.Query("publications")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("publications = ?", publications).Find(&booksDataList)

		fmt.Println("Printing fetched rows count : ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book data found for given publications."})
		}
	})

	//fetch complete book Info based on given book Title API
	router.GET("/booksinfo/from-title", func(ctx *gin.Context) {
		title := ctx.Query("title")

		var booksDataList []models.BooksInfo
		result := database.DB.Where("title = ?", title).Find(&booksDataList)

		fmt.Println("Printing fetched rows count: ", result.RowsAffected)
		if len(booksDataList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDataList})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book data found for given title."})
		}
	})

	//fetch n books based on given n value
	router.GET("/booksinfo/top-n-books/:n", func(ctx *gin.Context) {
		n := ctx.Param("n")
		var booksInfoList []models.BooksInfo
		nInteger, err := strconv.Atoi(n)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Error occurred during integer conversion for n!"})
		} else {
			result := database.DB.Limit(nInteger).Find(&booksInfoList)
			fmt.Println("For given n value, Fetched records count: ", result.RowsAffected)
			if len(booksInfoList) > 0 {
				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksInfoList})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Invalid n value to fetch n books!"})
			}
		}
	})

	//Update book data using bookId and based on given new data API
	router.PUT("/booksinfo/update-book-data", func(ctx *gin.Context) {

		var booksData models.BooksInfo
		var newBooksData models.BooksInfo
		err := ctx.ShouldBindJSON(&newBooksData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "No book id field in the given paylod."})
			return
		}
		isBookIdPresent := database.DB.Where("book_id = ?", newBooksData.BookId).First(&booksData)
		if isBookIdPresent.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "Book ID doesn't exit."})
		} else {
			var title, publications, author, genre string
			if newBooksData.Author != "" {
				author = newBooksData.Author
			}
			if newBooksData.Title != "" {
				title = newBooksData.Title
			}
			if newBooksData.Publications != "" {
				publications = newBooksData.Publications
			}
			if newBooksData.Genre != "" {
				genre = newBooksData.Genre
			}
			result := database.DB.Where("book_id = ?", booksData.BookId).Updates(models.BooksInfo{Title: title, Publications: publications, Author: author, Genre: genre})
			fmt.Println("Printing updated rows count: ", result.RowsAffected)
			ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "msg": "Updated the given data successfully."})
		}

		//Another way of reading PUT payload data from x-www-form-urlencoded
		// var booksData1 models.BooksInfo
		// title := ctx.PostForm("title")
		// author := ctx.PostForm("author")
		// publications := ctx.PostForm("publications")
		// genre := ctx.PostForm("genre")
		// book_id := ctx.PostForm("book_id")
		// isBookIdPresent1 := database.DB.Where("book_id = ?", book_id).First(&booksData1)
		// if isBookIdPresent1.Error != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "data": "Invalid Book Id in the payload"})
		// 	return
		// }
		// result1 := database.DB.Where("book_id = ?", book_id).Updates(models.BooksInfo{
		// 	Title: title, Author: author, Publications: publications, Genre: genre,
		// })
		// fmt.Println("Printing affected row count: ", result1.RowsAffected)
		// ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "data": "Updted the given data successfully!"})
	})

	//Get all the booksInfo for the given genre API
	//The input of genre: "Romance","Fantasy","Horror" or Romance,Fantasy,Horror
	router.GET("/booksinfo/from-genre", func(ctx *gin.Context) {
		genre := ctx.Query("genre")
		if genre == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Genre type is not found in the Params"})
			return
		}
		genre = strings.Trim(genre, "\"")
		var booksList []models.BooksInfo
		fmt.Println("Printing input genre:", genre[0])
		res1 := strings.Split(genre, ",")
		fmt.Println("Printing strings after splitting: ", res1)
		for i := 0; i < len(res1); i++ {
			var currBooksList []models.BooksInfo
			result := database.DB.Where("genre = ?", strings.Trim(res1[i], "\"")).Find(&currBooksList)
			fmt.Printf("Printing rows affected count: %v for genre: %s\n", result.RowsAffected, res1[i])
			booksList = append(booksList, currBooksList...)
		}
		if len(booksList) > 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksList})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No data found for given genres."})
		}
	})
	//Another way of implementing the above API in POST request by taking array of strings.
	//The input of genre: ["Romance","Horror","Fantasy"]
	router.POST("/booksinfo/from-genre", func(ctx *gin.Context) {
		var postBody struct {
			Genre []string `json:"genre" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&postBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Genre type is not mentioned in payload."})
		} else {
			var booksList []models.BooksInfo
			fmt.Println("Printing input genre data: ", postBody.Genre, postBody)
			for i := 0; i < len(postBody.Genre); i++ {
				var currBooksList []models.BooksInfo
				res := database.DB.Where("genre = ?", postBody.Genre[i]).Find(&currBooksList)
				fmt.Printf("Printing fetched rows count: %v for genre: %v\n", res.RowsAffected, postBody.Genre[i])
				booksList = append(booksList, currBooksList...)
			}
			if len(booksList) > 0 {
				ctx.JSON(http.StatusFound, gin.H{"status": "success", "data": booksList})
			} else {
				ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No data found for given genre type."})
			}
		}
	})

	//Lend & Return Books API's
	database.DB.AutoMigrate(&models.LendBooksInfo{})
	// Add to Cart API(Step before lend request)
	router.POST("/add-to-cart-books", func(ctx *gin.Context) {
		var borrowPayload struct {
			LendData []borrowBook `json:"add_to_cart_data" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&borrowPayload)
		fmt.Println("Printing err: ", err)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either addToCart_data/user_id/book_id field is not present"})
			return
		}
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(borrowPayload.LendData); i++ {
				if borrowPayload.LendData[i].BookId == 0 || borrowPayload.LendData[i].UserId == "" {
					fmt.Println("One of the add to cart fields missing: ", borrowPayload.LendData[i])
					return io.EOF
				}
				var currRec models.LendBooksInfo
				isRecPresent := tx.Where("user_id = ? AND book_id = ? AND status = ?", borrowPayload.LendData[i].UserId, borrowPayload.LendData[i].BookId, "pending").Find(&currRec)
				if isRecPresent.RowsAffected != 0 {
					return gorm.ErrInvalidData
				}
				currRes := tx.Save(models.LendBooksInfo{UserId: borrowPayload.LendData[i].UserId,
					BookId:            borrowPayload.LendData[i].BookId,
					Status:            "pending",
					BorrowedTimestamp: time.Now(),
				})
				if currRes.Error != nil {
					return currRes.Error
				}
			}
			return nil
		})
		if transactionRes == gorm.ErrInvalidData {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Add to Cart data is already existing for same input!"})
		} else if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either book_id/user_id is missing in the borrow Payload."})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Add To Cart Transaction failure occurred."})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": "success", "msg": "Add To Cart Request Created successfully!"})
		}
	})

	//Lend API to lend a book
	router.PUT("/borrow-books", func(ctx *gin.Context) {
		var borrowPayload struct {
			LendData []borrowBook `json:"lend_data" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&borrowPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either lend_data/user_id/book_id field is not present"})
			return
		}
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(borrowPayload.LendData); i++ {
				if borrowPayload.LendData[i].BookId == 0 || borrowPayload.LendData[i].UserId == "" {
					fmt.Println("One of the borrowBook fields missing: ", borrowPayload.LendData[i])
					return io.EOF
				}
				var currRec models.LendBooksInfo
				isRecPresent := tx.Where("user_id = ? AND book_id = ? AND status = ? AND returned_timestamp IS NULL", borrowPayload.LendData[i].UserId, borrowPayload.LendData[i].BookId, "borrowed").Find(&currRec)
				if isRecPresent.RowsAffected != 0 {
					return gorm.ErrInvalidData
				}
				currRes := tx.Save(models.LendBooksInfo{UserId: borrowPayload.LendData[i].UserId,
					BookId:            borrowPayload.LendData[i].BookId,
					Status:            "borrowed",
					BorrowedTimestamp: time.Now(),
				})
				if currRes.Error != nil {
					return currRes.Error
				}
				var deleteRec models.LendBooksInfo
				deleteRes := tx.Where("book_id = ? AND status = ?", borrowPayload.LendData[i].BookId, "pending").Delete(&deleteRec)
				if deleteRes.Error != nil {
					return deleteRes.Error
				}
			}
			return nil
		})
		if transactionRes == gorm.ErrInvalidData {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Borrow Data is already existing for same input!"})
		} else if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either book_id/user_id is missing in the borrow Payload."})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Lend Request transaction failure occurred."})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": "success", "msg": "Lend Request Created successfully!"})
		}
	})

	//DELETE API to delete lend cart books whose status = pending for given userId
	router.DELETE("/delete-add-to-cart-books", func(ctx *gin.Context) {
		var deletePayload struct {
			DeleteData []borrowBook `json:"delete_data" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&deletePayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "delete_data field is missing in the payload."})
			return
		}
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(deletePayload.DeleteData); i++ {
				if deletePayload.DeleteData[i].BookId == 0 || deletePayload.DeleteData[i].UserId == "" {
					fmt.Println("One of the fields in delete_book are missing: ", deletePayload.DeleteData[i])
					return io.EOF
				}
				var currRec models.LendBooksInfo
				// isCurrRecPresent := tx.Where("user_id = ? AND book_id = ? AND status = ?", deletePayload.DeleteData[i].UserId, deletePayload.DeleteData[i].BookId, "pending").Find(&currRec)
				// if isCurrRecPresent.RowsAffected == 0{
				// 	return gorm.ErrRecordNotFound
				// }
				currRes := tx.Where("user_id = ? AND book_id = ? AND status = ?", deletePayload.DeleteData[i].UserId, deletePayload.DeleteData[i].BookId, "pending").Delete(&currRec)
				if currRes.Error != nil {
					return currRes.Error
				}

			}
			return nil
		})
		if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either user_id/book_id is not present in payload data."})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Transaction failure while deleting the lend request book data"})
		} else {
			ctx.Status(http.StatusNoContent)
		}
	})

	//Return a book for given userId and bookId API
	router.POST("/return-books", func(ctx *gin.Context) {
		var returnPayload struct {
			ReturnData []borrowBook `json:"return_data" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&returnPayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "return_data field is missing in the Payload"})
			return
		}
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(returnPayload.ReturnData); i++ {
				if returnPayload.ReturnData[i].BookId == 0 || returnPayload.ReturnData[i].UserId == "" {
					fmt.Println("One of fields in return_book are missing: ", returnPayload.ReturnData[i])
					return io.EOF
				}
				var currRec models.LendBooksInfo
				isCurrRecPresent := tx.Where("user_id = ? AND book_id = ? AND status = ?", returnPayload.ReturnData[i].UserId, returnPayload.ReturnData[i].BookId, "borrowed").Find(&currRec)
				if isCurrRecPresent.RowsAffected == 0 {
					return gorm.ErrRecordNotFound
				}
				currRes := tx.Model(&currRec).Updates(models.LendBooksInfo{Status: "returned", ReturnedTimestamp: time.Now()})
				if currRes.Error != nil {
					return currRes.Error
				}
			}
			return nil
		})
		if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either user_id/book_id is not present in the Payload data."})
		} else if transactionRes == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "Borrow data not found given book_id and user_id!"})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Transaction failure when updating the return status."})
		} else {
			ctx.JSON(http.StatusAccepted, gin.H{"status": "success", "msg": "Returned book status updated successfully!"})
		}
	})

	//GET API for getting all the current addToCart BookIds for given userId
	router.GET("/add-to-cart-books", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "pending").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No add to cart books data found for given userId."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		}
	})

	//GET API for getting all the current addToCart BooksInfo instead of BookIds for given userId
	router.GET("/add-to-cart-books/booksinfo", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "pending").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		} else {
			newRes := []models.BooksInfo{}
			for i := 0; i < len(allRecs); i++ {
				var currBookInfo []models.BooksInfo
				currRecRes := database.DB.Where("book_id = ?", allRecs[i].BookId).Find(&currBookInfo)
				if currRecRes.RowsAffected != 0 {
					newRes = append(newRes, currBookInfo...)
				}
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": newRes})
		}
	})

	//GET API for getting all the current borrowed BooksId for given userId
	router.GET("/curr-borrowed-books", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "borrowed").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No borrowed books data found for given userId."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		}
	})

	//GET API for getting all the current borrowed BooksInfo for given userId
	router.GET("/curr-borrowed-books/booksinfo/", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "borrowed").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		} else {
			newRes := []models.BooksInfo{}
			for i := 0; i < len(allRecs); i++ {
				var currBookInfo []models.BooksInfo
				currRecRes := database.DB.Where("book_id = ?", allRecs[i].BookId).Find(&currBookInfo)
				if currRecRes.RowsAffected != 0 {
					newRes = append(newRes, currBookInfo...)
				}
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": newRes})
		}
	})

	//GET API for getting all the past borrowed booksId for given UserId
	router.GET("/past-borrowed-books", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "returned").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No returned books data found for given userId."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		}
	})

	//DELETE API for deleting past borrowed books for given UserId
	router.DELETE("/delete-past-borrowed-books", func(ctx *gin.Context) {
		var deletePayload struct {
			DeleteData []borrowBook `json:"delete_data" binding:"required"`
		}
		err := ctx.ShouldBindJSON(&deletePayload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "delete_data field is missing in the payload."})
			return
		}
		transactionRes := database.DB.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(deletePayload.DeleteData); i++ {
				if deletePayload.DeleteData[i].BookId == 0 || deletePayload.DeleteData[i].UserId == "" {
					fmt.Println("One of the fields in delete_book are missing: ", deletePayload.DeleteData[i])
					return io.EOF
				}
				var currRec models.LendBooksInfo
				currRes := tx.Where("user_id = ? AND book_id = ? AND status = ?", deletePayload.DeleteData[i].UserId, deletePayload.DeleteData[i].BookId, "returned").Delete(&currRec)
				if currRes.Error != nil {
					return currRes.Error
				}
			}
			return nil
		})
		if transactionRes == io.EOF {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Either user_id/book_id is not present in delete payload data."})
		} else if transactionRes != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "Transaction failure while deleting the past borrowed book data."})
		} else {
			ctx.Status(http.StatusNoContent)
		}
	})

	//GET API for getting all the past borrowed booksInfo for given UserId
	router.GET("/past-borrowed-books/booksinfo/", func(ctx *gin.Context) {
		userId := ctx.Query("user_id")
		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "user_id field should be passed in query to fetch data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("user_id = ? AND status = ?", strings.Trim(userId, "\""), "returned").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		} else {
			newRes := []models.BooksInfo{}
			for i := 0; i < len(allRecs); i++ {
				var currBookInfo []models.BooksInfo
				currRecRes := database.DB.Where("book_id = ?", allRecs[i].BookId).Find(&currBookInfo)
				if currRecRes.RowsAffected != 0 {
					newRes = append(newRes, currBookInfo...)
				}
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": newRes})
		}
	})

	//GET API to return all the curr borrowed users data for given bookID
	router.GET("/curr-borrowed-users", func(ctx *gin.Context) {
		bookId, err := strconv.Atoi(ctx.Query("book_id"))
		fmt.Println("Printing bookId: ", bookId, err, reflect.TypeOf(bookId))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "book_id should be passed as integer in query to fetch the data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("book_id = ? AND status = ?", bookId, "borrowed").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No borrowed users found for given book_id."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": "failure", "data": allRecs})
		}
	})

	//GET API to return all the past borrowed users data for given bookID
	router.GET("/past-borrowed-users", func(ctx *gin.Context) {
		bookId, err := strconv.Atoi(ctx.Query("book_id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "failure", "msg": "book_id should be passed as integer in query to fetch the data."})
			return
		}
		var allRecs []models.LendBooksInfo
		res := database.DB.Where("book_id = ? AND status = ?", bookId, "returned").Find(&allRecs)
		if res.RowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failure", "msg": "No returned users found for given book_id."})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": allRecs})
		}
	})
	server.Run(":5000")
}
