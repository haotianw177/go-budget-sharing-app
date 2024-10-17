// Establish a real-time connection to the server
let conn = new WebSocket("ws://" + window.location.host + "/ws");

let expenseChart; // Will hold the chart showing expenses by category
let categoryTotals = {}; // Keeps track of totals in each category

// When a message is received from the server
conn.onmessage = function (evt) {
  const data = JSON.parse(evt.data); // Convert the message from JSON to an object

  if (data.Type === "Expense") {
    // If the message is about a new expense
    addExpenseToList(data.Data); // Add it to the list on the page
    updateTotalExpenses(data.TotalExpenses); // Update the total expenses displayed
    categoryTotals = data.CategoryTotals; // Update the category totals
    updateChart(); // Refresh the chart to reflect new data
  } else if (data.Type === "Notification") {
    // If the message is a notification
    showNotification(data.Message); // Display the notification to the user
  }
};

// When the user submits a new expense
document.getElementById("expenseForm").addEventListener("submit", function (e) {
  e.preventDefault(); // Prevent the page from reloading

  // Gather the data from the form
  const expense = {
    User: document.getElementById("user").value,
    Description: document.getElementById("description").value,
    Amount: parseFloat(document.getElementById("amount").value),
    Category: document.getElementById("category").value,
  };

  // Send the new expense to the server
  fetch("/addExpense", {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // Tell the server we're sending JSON
    body: JSON.stringify(expense), // Convert the expense to a JSON string
  }).then((response) => {
    if (response.ok) {
      document.getElementById("expenseForm").reset(); // Clear the form
    } else {
      alert("Failed to add expense"); // Show an error if something went wrong
    }
  });
});

// Adds a new expense to the list on the page
function addExpenseToList(expense) {
  const expensesList = document.getElementById("expensesList"); // The list element
  const newExpense = document.createElement("li"); // Create a new list item
  // Set the content of the list item
  newExpense.innerHTML = `<strong>${expense.User}:</strong> ${
    expense.Description
  } - $${expense.Amount.toFixed(2)} (${expense.Category})`;
  expensesList.appendChild(newExpense); // Add the new expense to the list
}

// Updates the total expenses displayed on the page
function updateTotalExpenses(total) {
  document.getElementById("totalExpenses").textContent = total.toFixed(2); // Update the number
}

// Updates the chart showing expenses by category
function updateChart() {
  const ctx = document.getElementById("expenseChart").getContext("2d"); // Get the drawing area
  const categories = Object.keys(categoryTotals); // Get the names of the categories
  const amounts = Object.values(categoryTotals); // Get the amounts spent in each category
  const colors = ["#4dc9f6", "#f67019", "#f53794", "#537bc4", "#acc236"]; // Colors for the chart

  if (expenseChart) {
    // If the chart already exists, update it
    expenseChart.data.labels = categories;
    expenseChart.data.datasets[0].data = amounts;
    expenseChart.update(); // Redraw the chart
  } else {
    // If the chart doesn't exist, create it
    expenseChart = new Chart(ctx, {
      type: "pie",
      data: {
        labels: categories, 
        datasets: [{ data: amounts, backgroundColor: colors }],
      },
      options: {
        responsive: true, 
      },
    });
  }
}

// Displays a notification message on the page
function showNotification(message) {
  const notificationsDiv = document.getElementById("notifications"); // Where we'll put the message
  const notification = document.createElement("div"); // Create a new div for the message
  notification.className = "notification"; // Apply some styling
  notification.textContent = message; // Set the text to the message
  notificationsDiv.appendChild(notification); // Add the message to the page
}
