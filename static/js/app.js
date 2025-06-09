function analyzeURL() {
    const url = document.getElementById("urlInput").value;
    fetch("/analyze", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url }),
    })
      .then((response) => response.json())
      .then((data) => {
        document.getElementById("result").innerHTML =
          "<h3>Analysis Result</h3><pre>" +
          JSON.stringify(data, null, 2) +
          "</pre>";
      })
      .catch((err) => {
        document.getElementById("result").innerHTML =
          '<p style="color:red;">Error: ' + err.message + "</p>";
      });
  }
  