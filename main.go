package main

import (
    "encoding/json"            //  working with JSON data
    "fmt"                      //  printing messages to the console
    "html/template"            // for generating HTML pages
    "log"                      // for logging errors and information
    "net/http"                 // for creating web servers
    "sync"                     // for synchronizing data access
    "github.com/gorilla/websocket" // for real-time communication
)




// data structure "Budget" represents the shared budget among users
type Budget struct {
    Name        string        // the name of the budget
    TotalAmount float64       // the total amount available
    Expenses    []Expense     // a list of expenses
    Mutex       sync.Mutex    // a lock to prevent data conflicts
}

// expense represents a single expense entry
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
    upgrader     = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // Allows any connection
    clients      = make(map[*websocket.Conn]bool) // Keeps track of connected users
    clientsMutex sync.Mutex                       // Ensures safe access to clients
    broadcast    = make(chan interface{})         // Channel for sending messages to clients
    threshold    = 0.8                            // 80% threshold for budget warnings
)

func main() {
    // Set up routes for different pages and actions
    http.HandleFunc("/", homeHandler)            // Home page
    http.HandleFunc("/ws", wsHandler)            // WebSocket connection
    http.HandleFunc("/addExpense", addExpenseHandler) // Adding a new expense
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(".")))) // Serving static files

    fmt.Println("Server started on :8080")       // Inform that the server is running
    go handleMessages()                          // Start handling real-time messages
    log.Fatal(http.ListenAndServe(":8080", nil)) // Start the web server on port 8080
}

// homeHandler serves the main page where users can see and add expenses
func homeHandler(w http.ResponseWriter, r *http.Request) {
    budget.Mutex.Lock()                          // Lock the budget to prevent conflicts
    defer budget.Mutex.Unlock()                  // Unlock after the function finishes

    // Prepare the data to display on the page
    data := struct {
        Name          string
        TotalAmount   float64
        TotalExpenses float64
        Expenses      []Expense
    }{
        Name:          budget.Name,               // Budget name
        TotalAmount:   budget.TotalAmount,        // Total budget amount
        TotalExpenses: calculateTotalExpenses(),  // Calculate total expenses
        Expenses:      budget.Expenses,           // List of expenses
    }

    // Generate the HTML page using the template and data
    err := templates.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError) // Handle any errors
    }
}

// wsHandler establishes a WebSocket connection for real-time updates
func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)     // Upgrade the HTTP connection to WebSocket
    if err != nil {
        log.Println("WebSocket Upgrade Error:", err)
        return
    }

    clientsMutex.Lock()                          // Lock the clients map
    clients[conn] = true                         // Add the new client
    clientsMutex.Unlock()                        // Unlock the clients map
}

// handleMessages listens for messages to send to clients
func handleMessages() {
    for {
        msg := <-broadcast                       // Wait for a message to broadcast
        clientsMutex.Lock()                      // Lock the clients map
        for client := range clients {            // Send the message to each client
            err := client.WriteJSON(msg)         // Send the message as JSON
            if err != nil {
                log.Printf("WebSocket Error: %v", err)
                client.Close()                   // Close the connection if there's an error
                delete(clients, client)          // Remove the client from the list
            }
        }
        clientsMutex.Unlock()                    // Unlock the clients map
    }
}

// addExpenseHandler processes new expenses submitted by users
func addExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {             // Ensure the request is a POST
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var expense Expense
    err := json.NewDecoder(r.Body).Decode(&expense) // Read the expense from the request
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    budget.Mutex.Lock()                          // Lock the budget data
    budget.Expenses = append(budget.Expenses, expense) // Add the new expense
    totalExpenses := calculateTotalExpenses()    // Update total expenses
    categoryTotals := calculateCategoryTotals()  // Update totals by category
    budget.Mutex.Unlock()                        // Unlock the budget data

    // Prepare the message to send to clients
    msg := struct {
        Type           string
        Data           Expense
        TotalExpenses  float64
        CategoryTotals map[string]float64
    }{
        Type:           "Expense",               // Message type
        Data:           expense,                 // The new expense
        TotalExpenses:  totalExpenses,           // Updated total expenses
        CategoryTotals: categoryTotals,          // Updated category totals
    }

    broadcast <- msg                             // Send the message to all clients

    // Check if expenses exceed the threshold
    if totalExpenses > threshold*budget.TotalAmount {
        notification := struct {
            Type    string
            Message string
        }{
            Type:    "Notification",
            Message: fmt.Sprintf("Budget threshold of %.0f%% exceeded!", threshold*100),
        }
        broadcast <- notification                // Send a warning to all clients
    }

    // Respond to the user's request indicating success
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
