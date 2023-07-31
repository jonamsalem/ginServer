package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// func getBooks(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, books)
// }




// func createBook(c *gin.Context) {
// 	var newBook book
// 	if err := c.BindJSON(&newBook); err != nil {
// 		return
// 	}
// 	books = append(books, newBook)
// 	c.IndentedJSON(http.StatusCreated, newBook)
// }

// func getBookbyID(id string)(*book, error){
// 	for i,b := range books {
// 		if b.ID == id {
// 			return &books[i], nil
// 		}
// 	}
// 	return nil, errors.New("Book not found")

// }


// func bookByID (c *gin.Context) {
// 	id := c.Param("id")
// 	book, err := getBookbyID(id)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// 	}
// 	c.IndentedJSON(http.StatusOK, book)
// }


// func checkoutBook(c *gin.Context) {
// 	id , ok:= c.GetQuery("id")
// 	if !ok {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing id query parameter"}) 
// 		return
// 	}
// 	book, err := getBookbyID(id)
// 	if err!= nil {
//         c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found"})
//         return
//     }
// 	if book.Quantity <= 0 {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Book not available"})
//         return
// 	}
// 	book.Quantity-=1
// 	c.IndentedJSON(http.StatusOK, book)
// }


// func returnBook(c *gin.Context) {
// 	id, ok := c.GetQuery("id")
// 	if !ok {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing id query parameter"}) 
//         return
// 	}
// 	book, err := getBookbyID(id)
// 	fmt.Println(book, err)
// 	if err!= nil {
//         c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found"})
//         return
//     }
// 	book.Quantity+=1
// 	c.IndentedJSON(http.StatusOK, book)
// }

// func deleteBook(c *gin.Context){
// 	id, ok := c.GetQuery("id")
// 	if !ok {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing id query parameter"}) 
//         return
// 	}
// 	book, err := getBookbyID(id)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found"})
//         return
// 	}
// 	index := -1
// 	for i,b := range books{
// 		if b.ID == book.ID {
// 			index = i
// 		}
// 	}
// 	books = append(books[:index], books[index+1:]...)
// 	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted"})
			
// }


type Book struct {
	ID primitive.ObjectID 		`bson:"_id,omitempty"`
	Title string          		`bson:"title, omitempty `
	Description string 		 `bson:"description, omitempty `
	Author string			`bson:"author, omitempty ` 
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
        port = "8080"
    }	



	router := gin.Default()
	router.GET("/books", getAllBooks)
	router.GET("/oneBook", getOneBook)
    router.POST("/createBook", createBook)
	router.DELETE("/deleteBook", deleteBook)
	router.PUT("/updateBook", updateBook)
	// router.PATCH("checkout", checkoutBook)
	// router.PATCH("/return", returnBook)
	// router.POST("/books", createBook)
	// router.DELETE("/books", deleteBook)

	router.Run("localhost:" + port)
}

func dbConnection() (context.Context, *mongo.Collection, *mongo.Client) {
	MONGODB_URI := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGODB_URI))
	if err!= nil  {
        log.Fatal(err)
    }
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err!= nil {
        log.Fatal(err)
    }

	database := client.Database("books")
	booksCollection := database.Collection("books")
	return ctx, booksCollection, client


}


func getAllBooks(c *gin.Context) {
	ctx, booksCollection, client := dbConnection()
	
	resultCoded, err := booksCollection.Find(ctx, bson.M{})
	if err!= nil {
        log.Fatal(err)
    }
	var books []bson.M
	if err:= resultCoded.All(ctx, &books); err!= nil {
		log.Fatal(err)
    }
	c.IndentedJSON(http.StatusOK, books)
	defer client.Disconnect(ctx)
	

}

func getOneBook(c *gin.Context) {
	ctx, booksCollection, client := dbConnection()
	ID := c.Query("ID")
	fmt.Println(ID)
	if ID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing ID query parameter",
		})}

	objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID format",
			})
			return
		}	
	var book bson.M
   if err:= booksCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&book); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Book not found"})
		return
	}

c.IndentedJSON(http.StatusOK, book)
defer client.Disconnect(ctx)
}

func createBook(c *gin.Context) {
	ctx, booksCollection, client := dbConnection()
	var book Book 
	if err := c.BindJSON(&book); err!= nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book"})
			return
		}
	var err1 error
	booksCollection.InsertOne(ctx, book)
	if err1 != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": err1.Error(),
        })
        return
    }
	c.IndentedJSON(http.StatusCreated, "book created")
	defer client.Disconnect(ctx)
}

func deleteBook(c *gin.Context) {
	ctx, booksCollection, client := dbConnection()
	ID := c.Query("ID")
	if ID== "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": "Missing ID query parameter",
        })
        return
	}
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err!= nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": "Invalid ID format",
        })
        return
    }
	var err2 error
	booksCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err2!= nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": err2.Error(),
        })
        return
    }
    c.IndentedJSON(http.StatusOK, "book deleted")

	defer client.Disconnect(ctx)

      
    }
 
	

func updateBook(c *gin.Context) {
	ctx, booksCollection, client := dbConnection()
	ID := c.Query("ID")
	if ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": "Missing ID query parameter",
        })
        return}
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err!= nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": "Invalid ID format",
        })
        return
    }
    var updatedBook Book
	if err:= c.BindJSON(&updatedBook); err!= nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": "Invalid book",
        })
        return
    }
	update := bson.M{}
	if updatedBook.Description != "" {
        update["$set"] = bson.M{"description": updatedBook.Description}
    }
    if updatedBook.Author != "" {
        update["$set"] = bson.M{"author": updatedBook.Author}
    }
    if updatedBook.Title != "" {
        update["$set"] = bson.M{"title": updatedBook.Title}
    }

	result, err := booksCollection.UpdateOne(ctx, bson.M{"_id": objectID},update)
	
	if err!= nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

	fmt.Print(result)
    c.IndentedJSON(http.StatusAccepted, gin.H{
		"message": "Book updated"})

		defer client.Disconnect(ctx)
	
}