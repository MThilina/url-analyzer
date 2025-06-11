document.addEventListener("DOMContentLoaded", () => {
  const analyzeBtn = document.getElementById("analyzeBtn");
  const urlInput   = document.getElementById("urlInput");
  const resultCard = document.getElementById("resultCard");
  const output     = document.getElementById("output");

  analyzeBtn.addEventListener("click", async () => {
    const url = urlInput.value.trim();
    if (!url) return;

    // Reset previous output
    output.innerHTML = "";
    resultCard.classList.add("hidden");

    try {
      const res = await fetch("/analyze", {
        method:  "POST",
        headers: { "Content-Type": "application/json" },
        body:    JSON.stringify({ url })
      });

      // Attempt to parse JSON (even on error)
      const body = await res.json().catch(() => ({}));

      if (!res.ok) {
        // Prefer server-sent message field, else HTTP status text
        const msg = body.message || res.statusText || "Unknown error";
        alert(`Error Message : ${msg}`);
        return;
      }

      // Build an array of [label, value] per your model
      const fields = [
        ["URL Analyzed",       url],
        ["HTML Version",       body.htmlVersion],
        ["Page Title",         body.title],
        ["Has Login Form",     body.hasLoginForm ? "Yes" : "No"],
        ["Internal Links",     body.links.internal],
        ["External Links",     body.links.external],
        ["Inaccessible Links", body.links.inaccessible],
      ];

      // Append heading counts
      const headings = body.headings || {};
      for (const [tag, cnt] of Object.entries(headings)) {
        fields.push([`Count of ${tag}`, cnt]);
      }

      // Render into the output-grid
      fields.forEach(([label, val]) => {
        const lab = document.createElement("div");
        lab.className = "label";
        lab.textContent = label;

        const value = document.createElement("div");
        value.className = "value";
        value.textContent = val;

        output.append(lab, value);
      });

      // Show the results card
      resultCard.classList.remove("hidden");
    }
    catch (err) {
      alert(`Error Message : ${err.message}`);
    }
  });
});
