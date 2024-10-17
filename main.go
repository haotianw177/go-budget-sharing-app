package main

import (
    "encoding/json"            //   working with JSON data
    "fmt"                      //     printing messages to the console
    "html/template"            //   for generating HTML pages
    "log"                      //    for logging errors and information
    "net/http"                 //   for creating web servers
    "sync"                     //   for synchronizing data access
    "github.com/gorilla/websocket"    // for real-time communication
)






//   data structure "Budget" represents the shared budget among users
type Budget struct {
    Name        string        // the name of the budget
    TotalAmount float64       // the total amount available
    Expenses    []Expense     // a list of expenses
    Mutex       sync.Mutex    // a lock to prevent data conflicts
}

// data structure "expsnse" represents a single expense entry
type Expense struct {
    Description string  // what the expense was for
    Amount      float64 // how much was spent
    User        string  // who made the expense
    Category    string  // the category of the expense 
}





// global variables to manage templates, data, and connections
var (
    templates    = template.Must(template.ParseFiles("./static/index.html")) // loads the HTML template
    budget       = Budget{Name: "Monthly Shared Budget", TotalAmount: 1000.0}       // initializes the budget
    upgrader     = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // allows any connection
    clients      = make(map[*websocket.Conn]bool) // keeps track of connected users
    clientsMutex sync.Mutex                       // ensures safe access to clients
    broadcast    = make(chan interface{})         // channel for sending messages to clients
    threshold    = 0.8                            //    threshold for budget warnings
)

// seems like don't need it for vercel?>>>???

func main() {
    // this is set up routes for different pages and actions
    http.HandleFunc("/", homeHandler)            // for home page
    http.HandleFunc("/ws", wsHandler)            // for WebSocket connection
    http.HandleFunc("/addExpense", addExpenseHandler) // for adding a new expense
  

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(".")))) // for serving static files


    fmt.Println("Server started on :8080")       
    go handleMessages()                          // start handling real-time messages
    log.Fatal(http.ListenAndServe(":8080", nil)) // start server on port 8080
}





// homeHandler serves the main page where users can see and add expenses
func homeHandler(w http.ResponseWriter, r *http.Request) {
    budget.Mutex.Lock()                          // lock the budget to prevent conflicts
    defer budget.Mutex.Unlock()                  // unlock after the function finishes

    // prepare data to display on the page
    data := struct {
        Name          string
        TotalAmount   float64
        TotalExpenses float64
        Expenses      []Expense
    }{
        Name:          budget.Name,               
        TotalAmount:   budget.TotalAmount,       
        TotalExpenses: calculateTotalExpenses(),  
        Expenses:      budget.Expenses,           // list of expenses
    }

    // generate HTML page using template and data
    err := templates.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError) 
    }
}

// wsHandler establishes WebSocket connection for real-time updates
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)     // upgrade the HTTP connection to WebSocket
    if err != nil {
        log.Println("WebSocket Upgrade Error:", err)
        return
    }

    clientsMutex.Lock()                          // lock the clients map
    clients[conn] = true                         // add the new client
    clientsMutex.Unlock()                        // unlock the clients map
}



// handleMessages listens for messages to send to clients
func handleMessages() {
    for {
        msg := <-broadcast                       // wait for a message to broadcast
        clientsMutex.Lock()                      // lock the clients map
        for client := range clients {            // send the message to each client
            err := client.WriteJSON(msg)         // send the message as JSON
            if err != nil {
                log.Printf("WebSocket Error: %v", err)
                client.Close()                   // close the connection if there's an error
                delete(clients, client)          // remove the client from the list
            }
        }
        clientsMutex.Unlock()                    // now unlock the clients map
    }
}

// addExpenseHandler processes new expenses submitted by users
func addExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {             // this ensure the request is a POST
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var expense Expense
    err := json.NewDecoder(r.Body).Decode(&expense) // read the expense from the request
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    budget.Mutex.Lock()                          // lock the budget data
    budget.Expenses = append(budget.Expenses, expense) // add the new expense
    totalExpenses := calculateTotalExpenses()    // update total expenses
    categoryTotals := calculateCategoryTotals()  // update totals by category
    budget.Mutex.Unlock()                        // unlock the budget data

    // prepare the message to send to clients
    msg := struct {
        Type           string
        Data           Expense
        TotalExpenses  float64
        CategoryTotals map[string]float64
    }{
        Type:           "Expense",               // message type
        Data:           expense,                 // the new expense
        TotalExpenses:  totalExpenses,           // updated total expenses
        CategoryTotals: categoryTotals,          // updated category totals
    }

    broadcast <- msg                             // send the message to all clients

    // check if expenses exceed the threshold
    if totalExpenses > threshold*budget.TotalAmount {
        notification := struct {
            Type    string
            Message string
        }{
            Type:    "Notification",
            Message: fmt.Sprintf("Budget threshold of %.0f%% exceeded!", threshold*100),
        }
        broadcast <- notification                // send a warning to all clients
    }

    // respond to the user's request indicating success
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}





// funciton that adds up all the expenses
func calculateTotalExpenses() float64 {
    total := 0.0
    for _, expense := range budget.Expenses {
        total += expense.Amount               
    }
    return total
}

// function that adds up expenses by category
func calculateCategoryTotals() map[string]float64 {
    categoryTotals := make(map[string]float64)   // map to hold totals per category
    for _, expense := range budget.Expenses {
        categoryTotals[expense.Category] += expense.Amount // now add amount to the specific category
    }
    return categoryTotals
}
