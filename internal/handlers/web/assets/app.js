document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("clearButton").onclick = clearSelectedNumbers;
  const selectedNumbers = [];
  initializeNumberGrid(selectedNumbers);
  setupGenerateButton(selectedNumbers);

  function getBallClass(number) {
    if (number === 40) {
      return "ball-purple";
    } else if (number >= 30) {
      return "ball-red";
    } else if (number >= 20) {
      return "ball-green";
    } else if (number >= 10) {
      return "ball-orange";
    }
    return "ball-blue";
  }

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
        numberElement.classList.add("number-circle", "lotto-ball", getBallClass(parseInt(number)));
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

  function initializeNumberGrid(selectedNumbers) {
    const numberGrid = document.getElementById("numberGrid");
    for (let i = 1; i <= 40; i++) {
      const numberElement = document.createElement("div");
      numberElement.textContent = i;
      numberElement.dataset.number = i;
      numberElement.classList.add(
        "w-10",
        "h-10",
        "mx-auto",
        "lotto-ball",
        getBallClass(i)
      );
      numberElement.onclick = handleNumberClick;
      numberGrid.appendChild(numberElement);
    }
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

  function setupGenerateButton(selectedNumbers) {
    document.getElementById("generateButton").onclick = function () {
      const numLinesValue = document.getElementById("numLines").value;
      const numPerLineValue = document.getElementById("numPerLine").value;
      const numbersJoined = selectedNumbers.join(",");
      const apiUrl = `./numbers?lines=${numLinesValue}&numPerLine=${numPerLineValue}&numbersList=${numbersJoined}`;

      fetchNumbers(apiUrl);
    };
    fetchAppVersion();
  }

  // Function to fetch and display the application version
  function fetchAppVersion() {
    fetch("/version")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch app version");
        }
        return response.text();
      })
      .then((version) => {
        document.getElementById("appVersion").textContent = version;
      })
      .catch((error) => console.error("Error fetching app version:", error));
  }

  function fetchNumbers(apiUrl) {
    fetch(apiUrl)
      .then((response) => {
        if (!response.ok) {
          if (response.headers.get("Content-Type").includes("text/plain")) {
            return response.text().then((text) => {
              throw new Error(text);
            });
          } else {
            return response.json().then((data) => {
              throw new Error(
                data.message || `HTTP error! status: ${response.status}`
              );
            });
          }
        }
        return response.json();
      })
      .then((data) => {
        displayNumbers(data);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
        displayError(error.toString().replace("Error: ", ""));
        document.getElementById("appVersion").textContent = "Unavailable";
      });
  }
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

      const lineContainer = document.createElement("div");
      lineContainer.classList.add(
        "line-container",
        "flex",
        "items-center",
        "mb-2"
      );

      // Create and append line number element
      const lineNumberElem = document.createElement("div");
      lineNumberElem.textContent = `#${index + 1}`;
      lineNumberElem.classList.add("line-number");
      lineContainer.appendChild(lineNumberElem);

      // Create and append line numbers
      const lineElem = document.createElement("div");
      lineElem.classList.add("flex", "flex-wrap", "justify-center");
      line.forEach((number) => {
        const numberElement = document.createElement("div");
        numberElement.textContent = number;
        numberElement.classList.add(
          "generated-number-circle",
          "lotto-ball",
          getBallClass(number)
        );
        lineElem.appendChild(numberElement);
      });
      lineContainer.appendChild(lineElem);

      container.appendChild(lineContainer);
    });
  }

  function displayError(errorMessage) {
    const container = document.getElementById("numbersContainer");
    container.innerHTML = `<div class="error">Error: ${errorMessage}</div>`;
  }
});
