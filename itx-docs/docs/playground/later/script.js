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
    .then(response => { 
        // User output
        document.getElementById("code-output").value = response.return
        console.log(response.files)
        // File I/O
        for (const file of response.files) {
            for (const filename in file) {
                const fileData = file[filename]
                let parent = document.getElementById("files")
                let newFile = document.createElement('div')

                let exists = Array.from(parent.children).some(child => child.textContent === filename)
                if (exists && !fileData.type == "remove") {
                    newFile.onclick = () => {
                        document.getElementById("file-view").value = file[filename].value
                    }
                    continue
                }
                console.log(fileData.type)

                switch (fileData.type) {
                case "write":
                    newFile.textContent = filename
                    parent.appendChild(newFile)
                    break
                case "remove":
                    console.log("I reached here")
                    for (let child of parent.children) {
                        if (child.textContent == filename) {
                            console.log(filename)
                            parent.removeChild(child)
                            document.getElementById("file-view").value = ""
                            break
                        }
                    }
                    break
                case "append":
                    for (let child of parent.children) {
                        if (child.textContent == filename) {
                            //TODO
                            break
                        }
                    }
                    break
                }

                newFile.onclick = () => {
                    document.getElementById("file-view").value = fileData.value
                }
            }
        }
    })

    document.getElementById("file-view").value = ""
}

function Return() {
    window.open("https://docs.codekeg.dev", '_top').focus()
}