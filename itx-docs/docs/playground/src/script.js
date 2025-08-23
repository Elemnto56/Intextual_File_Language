window.addEventListener("load", function () {
    console.log("Page has loaded")
})

function EnableLight() {
    let x = document.getElementById("bd")
    let y = document.getElementById("bt")

    if (x.className == "") {
        x.className = "dark-mode"
        y.innerText = "Enable Light Mode?"
    } else {
        y.innerText = "Enable Dark Mode?"
        x.className = ""
    }
}

function Run() {
    let x = document.getElementById("code-input")

    fetch("http://localhost:8080/api/send", {
        method: "POST",
        headers:  {"Content-Type": "application/json"},
        body: JSON.stringify({ data: x.value })
    })
    .then(res => res.json())
    .then(response => document.getElementById("code-output").value = response.message)
}