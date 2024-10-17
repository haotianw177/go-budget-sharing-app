// haha

// establish the real-time connection to the server
let conn = new WebSocket("wss://" + window.location.host + "/ws");

let expenseChart; // this will hold the chart showing expenses by category
let categoryTotals = {}; // this keeps track of totals in each category




// when a message is received from the server
conn.onmessage = function (evt) {
  const data = JSON.parse(evt.data); // first, convert the message from JSON to an object

  if (data.Type === "Expense") {
    // if the message is about a new expense
    addExpenseToList(data.Data); // add it to the list on the page
    updateTotalExpenses(data.TotalExpenses); // update the total expenses displayed
    categoryTotals = data.CategoryTotals; // update the category totals
    updateChart(); // now refresh the chart to reflect new data
  } else if (data.Type === "Notification") {
    // if the message is a notification
    showNotification(data.Message); // display the notification to the user
  }
};




// when the user submits a new expense
document.getElementById("expenseForm").addEventListener("submit", function (e) {
  e.preventDefault(); //prevent the page from reloading

  // gather the data from the form
  const expense = {
    User: document.getElementById("user").value,
    Description: document.getElementById("description").value,
    Amount: parseFloat(document.getElementById("amount").value),
    Category: document.getElementById("category").value,
  };

  // send the new expense to the server
  fetch("/addExpense", {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // tell the server we're sending JSON
    body: JSON.stringify(expense), // convert the expense to a JSON string
  }).then((response) => {
    if (response.ok) {
      document.getElementById("expenseForm").reset(); // clear the form
    } else {
      alert("Failed to add expense");
    }
  });
});


// adds a new expense to the list on the page
function addExpenseToList(expense) {
  const expensesList = document.getElementById("expensesList"); // the list element
  const newExpense = document.createElement("li"); // create a new list item
  // set the content of the list item
  newExpense.innerHTML = `<strong>${expense.User}:</strong> ${
    expense.Description
  } - $${expense.Amount.toFixed(2)} (${expense.Category})`;
  expensesList.appendChild(newExpense); // add the new expense to the list
}

// updates the total expenses displayed on the page
function updateTotalExpenses(total) {
  document.getElementById("totalExpenses").textContent = total.toFixed(2); // update the number
}

// updates the chart showing expenses by category
function updateChart() {
  const ctx = document.getElementById("expenseChart").getContext("2d"); // get canvas drawing area
  const categories = Object.keys(categoryTotals); // get   names of the categories
  const amounts = Object.values(categoryTotals); // Get the amounts spent in each category
  const colors = ["#4dc9f6", "#f67019", "#f53794", "#537bc4", "#acc236", "#0000ff", "#00ff00"];

  if (expenseChart) {
    // if  chart already exists, update it
    expenseChart.data.labels = categories;
    expenseChart.data.datasets[0].data = amounts;
    expenseChart.update(); // redraw the chart
  } else {
    // if the chart doesn't exist, create it
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





// displays a notification message on the page
function showNotification(message) {
  const notificationsDiv = document.getElementById("notifications"); // where we'll put the message
  const notification = document.createElement("div"); // create new div for the message
  notification.className = "notification"; 
  notification.textContent = message; // set the text to the message
  notificationsDiv.appendChild(notification); // add the message to the page
}
