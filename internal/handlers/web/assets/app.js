document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("clearButton").onclick = clearSelectedNumbers;
  const numberGrid = document.getElementById("numberGrid");
  const selectedNumbers = [];

  // Function to toggle selection
  function toggleNumberSelection(number) {
    const index = selectedNumbers.indexOf(number);
    if (index > -1) {
      selectedNumbers.splice(index, 1); // Remove number if already selected
    } else {
      selectedNumbers.push(number); // Add number if not already selected
    }
  }

  // Function to update the display of selected numbers
  function updateDisplay() {
    const selectedNumbersElement = document.getElementById("selectedNumbers");
    selectedNumbersElement.innerHTML = ""; // Clear previous content
    if (selectedNumbers.length > 0) {
      selectedNumbers.forEach((number) => {
        const numberElement = document.createElement("div");
        numberElement.textContent = number;
        numberElement.classList.add("number-circle");
        selectedNumbersElement.appendChild(numberElement);
      });
    } else {
      selectedNumbersElement.textContent = "None";
    }
  }

  // Function to handle number click
  function handleNumberClick(event) {
    const number = event.target.dataset.number;
    event.target.classList.toggle("selected");
    toggleNumberSelection(number);
    updateDisplay();
  }

  // Initialize the number grid
  for (let i = 1; i <= 40; i++) {
    const numberElement = document.createElement("div");
    numberElement.textContent = i;
    numberElement.dataset.number = i;
    numberElement.classList.add(
      "w-10",
      "h-10",
      "bg-blue-500",
      "text-white",
      "flex",
      "items-center",
      "justify-center",
      "rounded-full",
      "mx-auto"
    );
    numberElement.onclick = handleNumberClick;
    numberGrid.appendChild(numberElement);
  }

  // Function to handle generate button click
  // Function to clear selected numbers and generated lines
  function clearSelectedNumbers() {
    selectedNumbers.length = 0; // Clear the array
    updateDisplay();
    // Clear the generated lines
    const numbersContainer = document.getElementById("numbersContainer");
    numbersContainer.innerHTML = "None";
    // Use a different approach to remove the 'selected' class from all elements
    const numberElements = numberGrid.querySelectorAll(".selected");
    numberElements.forEach(function (element) {
      element.classList.remove("selected");
    });
  }

  document.getElementById("generateButton").onclick = function () {
    // Construct the API URL
    let numLinesValue = document.getElementById("numLines").value;
    let numPerLineValue = document.getElementById("numPerLine").value;
    let numbersJoined = selectedNumbers.join(",");
    let apiUrl = `./numbers?lines=${numLinesValue}&numPerLine=${numPerLineValue}&numbersList=${numbersJoined}`;

    // Fetch data from the API
    fetch(apiUrl)
      .then(response => response.json().then(data => ({
        status: response.status,
        body: data
      })))
      .then(obj => {
        if (obj.status !== 200) {
          throw new Error(obj.body.message || `HTTP error! status: ${obj.status}`);
        }
        displayNumbers(obj.body);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
        displayError(error.message);
      });
  };
});
// Move the clear button event listener setup inside the DOMContentLoaded event where clearSelectedNumbers is defined
// This code block is removed as it is now redundant

function displayNumbers(data) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = ""; // Clear previous results

  data.lines.forEach((line, index) => {
    if (index > 0) {
      // Add a divider before each new line except the first
      const divider = document.createElement("div");
      divider.classList.add("line-divider");
      container.appendChild(divider);
    }
    const lineElem = document.createElement("div");
    lineElem.classList.add("flex", "flex-wrap", "justify-center", "mb-2");
    line.forEach((number) => {
      const numberElement = document.createElement("div");
      numberElement.textContent = number;
      numberElement.classList.add("generated-number-circle");
      lineElem.appendChild(numberElement);
    });
    container.appendChild(lineElem);
  });
}

function displayError(error) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = `<div class="error">Error: ${error.message}</div>`;
}
