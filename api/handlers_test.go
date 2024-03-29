package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/starter-api/api"
	"gotest.tools/assert"
)

var seedData = api.SeedData()

func TestRootHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)
	assert.Equal(t, response.Code, http.StatusOK)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, string(body), "API is running")
}

func findBookById(books []*api.Book, id int) *api.Book {
	for _, b := range books {
		if b.ID == id {
			return b
		}
	}
	return nil
}

func mapToPointers(books []api.Book) []*api.Book {
	bookPointers := make([]*api.Book, len(books))
	for i := range bookPointers {
		bookPointers[i] = &books[i]
	}
	return bookPointers
}

func TestGetBookHandler(t *testing.T) {
	bookId := 2
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%d", bookId), nil)
	response := executeRequest(req)

	assert.Equal(t, response.Code, http.StatusOK)

	book := api.Book{}
	err := json.Unmarshal(response.Body.Bytes(), &book)
	if err != nil {
		log.Fatal(err)
	}
	assert.DeepEqual(t, &book, findBookById(seedData, bookId))
}

func TestGetBookHandlerNotFound(t *testing.T) {
	bookId := 10300
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%d", bookId), nil)
	response := executeRequest(req)

	assert.Equal(t, response.Code, http.StatusNotFound)
	error := api.ErrorResponse{}

	err := json.Unmarshal(response.Body.Bytes(), &error)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, error.Error, fmt.Sprintf("Could not find Book with id %d", bookId))
}

func TestListBooksHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/books", nil)
	response := executeRequest(req)

	assert.Equal(t, response.Code, http.StatusOK)
	books := make([]api.Book, 0)
	err := json.Unmarshal(response.Body.Bytes(), &books)
	if err != nil {
		log.Fatal(err)
	}

	// Need an array of pointers to the original Book structs to pass to the findBookById func
	bookPointers := mapToPointers(books)

	for _, seedBook := range seedData {
		matchingBook := findBookById(bookPointers, seedBook.ID)
		if matchingBook != nil {
			assert.DeepEqual(t, seedBook, matchingBook)
		} else {
			t.Fatal(fmt.Printf("Did not find seed book %d in results", seedBook.ID))
		}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rr, req)

	return rr
}
