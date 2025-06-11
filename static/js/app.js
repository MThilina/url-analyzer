document.addEventListener("DOMContentLoaded", () => {
  const analyzeBtn = document.getElementById("analyzeBtn");
  const urlInput   = document.getElementById("urlInput");
  const resultCard = document.getElementById("resultCard");
  const output     = document.getElementById("output");
  const loader     = document.getElementById("loader");

  analyzeBtn.addEventListener("click", async () => {
    const url = urlInput.value.trim();
    if (!url) return;

    // Reset UI
    output.innerHTML = "";
    resultCard.classList.add("hidden");

    // Show loader
    loader.classList.remove("hidden");

    try {
      const res = await fetch("/analyze", {
        method:  "POST",
        headers: { "Content-Type": "application/json" },
        body:    JSON.stringify({ url })
      });

      const body = await res.json().catch(() => ({}));

      if (!res.ok) {
        const msg = body.message || res.statusText || "Unknown error";
        alert(`Error Message : ${msg}`);
      } else {
        // Build fields aligned with your Go model
        const fields = [
          ["URL Analyzed",       url],
          ["HTML Version",       body.htmlVersion],
          ["Page Title",         body.title],
          ["Has Login Form",     body.hasLoginForm ? "Yes" : "No"],
          ["Internal Links",     body.links.internal],
          ["External Links",     body.links.external],
          ["Inaccessible Links", body.links.inaccessible],
        ];
        for (const [tag, cnt] of Object.entries(body.headings || {})) {
          fields.push([`Count of ${tag}`, cnt]);
        }

        // Render results
        fields.forEach(([label, val]) => {
          const lab = document.createElement("div");
          lab.className = "label";
          lab.textContent = label;

          const value = document.createElement("div");
          value.className = "value";
          value.textContent = val;

          output.append(lab, value);
        });
        resultCard.classList.remove("hidden");
      }
    } catch (err) {
      alert(`Error Message : ${err.message}`);
    } finally {
      // Always hide loader
      loader.classList.add("hidden");
    }
  });
});
