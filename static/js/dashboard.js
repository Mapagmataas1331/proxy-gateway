document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("search");
  const table = document.getElementById("logs-table");
  const rows = table.tBodies[0].rows;

  searchInput.addEventListener("input", () => {
    const filter = searchInput.value.toLowerCase();

    for (let row of rows) {
      const cellsText = Array.from(row.cells)
        .map((cell) => cell.textContent.toLowerCase())
        .join(" ");
      row.style.display = cellsText.includes(filter) ? "" : "none";
    }
  });
});
