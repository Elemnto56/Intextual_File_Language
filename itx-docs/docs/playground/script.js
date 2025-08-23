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

    fetch("https://api.codekeg.dev/run-code", {
        method: "POST",
        headers:  {"Content-Type": "application/json"},
        body: JSON.stringify({ data: x.value })
    })
    .then(res => res.json())
    .then(response => { console.log("We have recieved: ", response.message), document.getElementById("code-output").value = response.message})
}

function Return() {
    window.open("https://docs.codekeg.dev", '_top').focus()
}