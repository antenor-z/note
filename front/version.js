function getVersion() {
    fetch(`${window.API_URL}/version`,
        {
            method: "GET",
        })
        .then(response => response.json())
        .then(data => {
            document.getElementById("version").innerText = "v" + data.version
        })
}
getVersion()