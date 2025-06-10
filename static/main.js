// Constants and configuration.
const CONFIG = {
  API_BASE_URL: document.baseURI + "/api",
  // Refresh every 30 seconds.
  REFRESH_INTERVAL: 30_000,
};

class TempMailApp {
  constructor() {
    console.info("App loaded!");
  }
}

// Initialize app when DOM is loaded.
document.addEventListener("DOMContentLoaded", () => {
  new TempMailApp();
});

export default TempMailApp;
