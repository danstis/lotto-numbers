document
  .getElementById("generateButton")
  .addEventListener("click", function () {
    // Get input values
    let numLines = document.getElementById("numLines").value;
    let numPerLine = document.getElementById("numPerLine").value;
    let numbers = document.getElementById("numbers").value;

    // Construct the API URL
    let apiUrl = `./numbers?lines=${numLines}&numPerLine=${numPerLine}&numbersList=${numbers}`;

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
  });

function displayNumbers(data) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = ""; // Clear previous results

  data.lines.forEach((line) => {
    const lineElem = document.createElement("div");
    lineElem.textContent = `Line: ${line.join(", ")}`;
    container.appendChild(lineElem);
  });
}

function displayError(error) {
  const container = document.getElementById("numbersContainer");
  container.innerHTML = `<div class="error">Error: ${error.message}</div>`;
}
