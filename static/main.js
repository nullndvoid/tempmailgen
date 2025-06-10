// Constants and configuration.
const CONFIG = {
  API_BASE_URL: document.baseURI + "/api",
  // Refresh every 30 seconds.
  REFRESH_INTERVAL: 30_000,
};

class TempMailApp {
  constructor() {
    console.info("App loaded!");
    this.bindEvents();
  }

  bindEvents() {
    const btn = document.querySelector("#new-mailbox-submit");
    if (btn) {
      btn.addEventListener("click", (e) => {
        // Don't submit the form immediately.
        e.preventDefault();
        this.getNewMailbox();
      });
    }
  }

  getNewMailbox() {
    console.info("Getting new mailbox!");
  }
}

// Initialize app when DOM is loaded.
document.addEventListener("DOMContentLoaded", () => {
  new TempMailApp();
});

export default TempMailApp;
