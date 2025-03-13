document.getElementById("shortenForm").addEventListener("submit", async function(event) {
    event.preventDefault();
  
    const urlInput = document.getElementById("url").value;
    const resultDiv = document.getElementById("result");
  
    try {
      const response = await fetch("/links", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url: urlInput })
      });
  
      if (!response.ok) {
        throw new Error("Failed to shorten URL");
      }
  
      const data = await response.json();
      if (!data.shortened) {
        throw new Error("Invalid response from server");
      }
  
      const shortenedUrl = `${window.location.origin}/${data.shortened}`;
  
      resultDiv.innerHTML = `
        <p><strong>Original:</strong> <a href="${data.original}" target="_blank">${data.original}</a></p>
        <p><strong>Shortened:</strong></p>
        <div class="shortened-container">
          <a href="${shortenedUrl}" target="_blank">${shortenedUrl}</a>
          <button class="copy-btn" onclick="copyToClipboard('${shortenedUrl}')">Copy</button>
        </div>
      `;
    } catch (error) {
      resultDiv.innerHTML = `<p style="color: red;">Error: ${error.message}</p>`;
    }
  });
  
  function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
      alert("Shortened URL copied to clipboard!");
    }).catch(err => {
      console.error("Failed to copy:", err);
    });
  }
  