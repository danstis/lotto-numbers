document.addEventListener("DOMContentLoaded", function () {
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
    document.getElementById("selectedNumbers").textContent =
      selectedNumbers.join(", ");
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
    numberElement.classList.add("w-10", "h-10", "bg-blue-500", "text-white", "flex", "items-center", "justify-center", "rounded-full", "mx-auto");
    numberElement.onclick = handleNumberClick;
    numberGrid.appendChild(numberElement);
  }

  // Function to handle generate button click
  // Function to clear selected numbers
  function clearSelectedNumbers() {
    selectedNumbers.length = 0; // Clear the array
    updateDisplay();
    const numberElements = numberGrid.getElementsByClassName("selected");
    while (numberElements.length) {
      numberElements[0].classList.remove("selected");
    }
  }

  document.getElementById("generateButton").onclick = function () {
    // Construct the API URL
    let numbersJoined = selectedNumbers.join(", ");
    let apiUrl = `./numbers?lines=${numLines}&numPerLine=${numPerLine}&numbersList=${numbersJoined}`;

    // Fetch data from the API
    fetch(apiUrl)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => displayNumbers(data))
      .catch((error) => {
        console.error("Error fetching data:", error);
        displayError(error);
      });
  };
});
// Add event listener to the clear button
document.getElementById("clearButton").onclick = clearSelectedNumbers;

function displayNumbers(data) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = ""; // Clear previous results

  data.lines.forEach((line, index) => {
    const lineElem = document.createElement("div");
    lineElem.textContent = `Line ${index + 1}: ${line.join(", ")}`;
    container.appendChild(lineElem);
  });
}

function displayError(error) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = `<div class="error">Error: ${error.message}</div>`;
}