// ===================================
// Code Playground TypeScript
// ===================================
// @ts-ignore
// Theme Management
var ThemeManager = /** @class */ (function () {
    function ThemeManager() {
        this.currentTheme = 'light';
        this.themeToggleBtn = document.querySelector('.theme-toggle');
        this.init();
    }
    ThemeManager.prototype.init = function () {
        var _this = this;
        // Check for saved theme preference or default to 'light'
        var savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            this.currentTheme = savedTheme;
        }
        else {
            // Check system preference
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                this.currentTheme = 'dark';
            }
        }
        // Apply the theme
        this.applyTheme(this.currentTheme);
        // Add click listener to toggle button
        if (this.themeToggleBtn) {
            this.themeToggleBtn.addEventListener('click', function () { return _this.toggleTheme(); });
        }
    };
    ThemeManager.prototype.applyTheme = function (theme) {
        if (theme === 'dark') {
            document.body.classList.add('dark-theme');
        }
        else {
            document.body.classList.remove('dark-theme');
        }
        this.currentTheme = theme;
        localStorage.setItem('theme', theme);
    };
    ThemeManager.prototype.toggleTheme = function () {
        var newTheme = this.currentTheme === 'light' ? 'dark' : 'light';
        this.applyTheme(newTheme);
    };
    ThemeManager.prototype.getTheme = function () {
        return this.currentTheme;
    };
    return ThemeManager;
}());
// File Explorer Management
var FileExplorer = /** @class */ (function () {
    function FileExplorer() {
        this.selectedFile = null;
        this.fileContents = new Map();
        this.fileTree = document.getElementById('file-tree');
        this.init();
    }
    FileExplorer.prototype.init = function () {
        var _this = this;
        // Add click listener for file selection
        if (this.fileTree) {
            this.fileTree.addEventListener('click', function (e) {
                var _b;
                var target = e.target;
                var fileElement = target.closest('.file');
                var folderElement = target.closest('.folder');
                if (fileElement) {
                    _this.selectFile(fileElement);
                }
                else if (folderElement && ((_b = target.closest('.folder')) === null || _b === void 0 ? void 0 : _b.firstChild) === target) {
                    // Only toggle if clicking on the folder itself, not its content
                    _this.toggleFolder(folderElement);
                }
            });
        }
    };
    FileExplorer.prototype.selectFile = function (fileElement) {
        var _b;
        // Remove previous selection
        if (this.selectedFile) {
            this.selectedFile.classList.remove('selected');
        }
        // Add selection to clicked file
        fileElement.classList.add('selected');
        this.selectedFile = fileElement;
        var fileName = ((_b = fileElement.textContent) === null || _b === void 0 ? void 0 : _b.trim()) || '';
        var fileExt = fileElement.getAttribute('data-ext') || '';
        // Load file content into the code input if it exists
        var content = this.fileContents.get(fileName);
        if (content !== undefined && editor) {
            editor.setInput(content);
        }
        // Trigger custom event for file selection
        var event = new CustomEvent('fileSelected', {
            detail: { fileName: fileName, fileExt: fileExt, content: content }
        });
        document.dispatchEvent(event);
    };
    FileExplorer.prototype.toggleFolder = function (folderElement) {
        folderElement.classList.toggle('open');
    };
    /**
     * Add a file to the file tree
     * @param fileName - The name of the file to add
     * @param dataExt - Optional data extension attribute
     * @param content - Optional file content
     */
    FileExplorer.prototype.addFile = function (fileName, dataExt, content) {
        if (!this.fileTree)
            return;
        var fileDiv = document.createElement('div');
        fileDiv.className = 'file';
        fileDiv.textContent = fileName;
        // Only add data-ext attribute if provided
        if (dataExt && dataExt.trim() !== '') {
            fileDiv.setAttribute('data-ext', dataExt);
        }
        // Store content if provided
        if (content !== undefined) {
            this.fileContents.set(fileName, content);
        }
        this.fileTree.appendChild(fileDiv);
    };
    FileExplorer.prototype.removeFile = function (filename) {
        var _b;
        var possibleFile = document.querySelector("#file-tree > *");
        if ((possibleFile === null || possibleFile === void 0 ? void 0 : possibleFile.innerHTML) == filename) {
            (_b = document.querySelector("#file-tree")) === null || _b === void 0 ? void 0 : _b.removeChild(possibleFile);
        }
        else {
            console.error("The file set for deletion did not exist");
        }
    };
    /**
     * Set or update file content
     * @param fileName - The name of the file
     * @param content - The content to set
     */
    FileExplorer.prototype.setFileContent = function (fileName, content) {
        this.fileContents.set(fileName, content);
    };
    /**
     * Get file content
     * @param fileName - The name of the file
     * @returns The file content or undefined if not found
     */
    FileExplorer.prototype.getFileContent = function (fileName) {
        return this.fileContents.get(fileName);
    };
    /**
     * Update content of currently selected file with current editor input
     */
    FileExplorer.prototype.updateSelectedFileContent = function () {
        var _b;
        if (this.selectedFile && editor) {
            var fileName = ((_b = this.selectedFile.textContent) === null || _b === void 0 ? void 0 : _b.trim()) || '';
            var currentContent = editor.getInput();
            this.fileContents.set(fileName, currentContent);
        }
    };
    /**
     * Get all files and their contents
     */
    FileExplorer.prototype.getAllFiles = function () {
        return new Map(this.fileContents);
    };
    /**
     * Add a folder to the file tree
     * @param folderName - The name of the folder
     * @returns The folder content div where child items can be added
     */
    FileExplorer.prototype.addFolder = function (folderName) {
        if (!this.fileTree)
            return null;
        var folderDiv = document.createElement('div');
        folderDiv.className = 'folder';
        folderDiv.textContent = folderName;
        var folderContent = document.createElement('div');
        folderContent.className = 'folder-content';
        folderDiv.appendChild(folderContent);
        this.fileTree.appendChild(folderDiv);
        return folderContent;
    };
    /**
     * Clear all files from the file tree
     */
    FileExplorer.prototype.clearFiles = function () {
        if (this.fileTree) {
            this.fileTree.innerHTML = '';
        }
        this.selectedFile = null;
    };
    /**
     * Get the currently selected file
     */
    FileExplorer.prototype.getSelectedFile = function () {
        var _b;
        if (!this.selectedFile)
            return null;
        return {
            fileName: ((_b = this.selectedFile.textContent) === null || _b === void 0 ? void 0 : _b.trim()) || '',
            fileExt: this.selectedFile.getAttribute('data-ext') || ''
        };
    };
    return FileExplorer;
}());
// Editor Management
var Editor = /** @class */ (function () {
    function Editor() {
        this.codeInput = document.getElementById('code-input');
        this.codeOutput = document.getElementById('code-output');
        this.init();
    }
    Editor.prototype.init = function () {
        var _this = this;
        // Add tab support in the input textarea
        if (this.codeInput) {
            this.codeInput.addEventListener('keydown', function (e) {
                if (e.key === 'Tab') {
                    e.preventDefault();
                    var start = _this.codeInput.selectionStart;
                    var end = _this.codeInput.selectionEnd;
                    var value = _this.codeInput.value;
                    // Insert tab at cursor
                    _this.codeInput.value = value.substring(0, start) + '\t' + value.substring(end);
                    // Move cursor after the tab
                    _this.codeInput.selectionStart = _this.codeInput.selectionEnd = start + 1;
                }
                // Ctrl/Cmd + Enter to run
                if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
                    e.preventDefault();
                    _this.run();
                }
            });
        }
    };
    Editor.prototype.run = function () {
        // Call the Run function that should be defined elsewhere
        if (typeof window.Run === 'function') {
            window.Run();
        }
    };
    Editor.prototype.getInput = function () {
        var _b;
        return ((_b = this.codeInput) === null || _b === void 0 ? void 0 : _b.value) || '';
    };
    Editor.prototype.setInput = function (value) {
        if (this.codeInput) {
            this.codeInput.value = value;
        }
    };
    Editor.prototype.getOutput = function () {
        var _b;
        return ((_b = this.codeOutput) === null || _b === void 0 ? void 0 : _b.value) || '';
    };
    Editor.prototype.setOutput = function (value) {
        if (this.codeOutput) {
            this.codeOutput.value = value;
        }
    };
    Editor.prototype.clearInput = function () {
        this.setInput('');
    };
    Editor.prototype.clearOutput = function () {
        this.setOutput('');
    };
    return Editor;
}());
// Server Management
var Server = /** @class */ (function () {
    function Server() {
    }
    Server.prototype.runCode = function () {
        var codeInput = document.getElementById("code-input");
        var codeOutput = document.getElementById("code-output");
        fetch("https://api.codekeg.dev/run-code", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ data: codeInput.value })
        })
            .then(function (res) { return res.json(); })
            .then(function (response) {
            console.log(response)
            codeOutput.innerHTML = response.return;

            if (response.filedata != null) {
                for (let file of response.filedata) {
                    switch (file.Type) {
                        case "Write":
                            addFile(file.FileName, undefined, file.Content);
                            break;
                        case "Append":
                            appendFileContent(file.FileName, file.Content);
                            break;
                        case "Delete":
                            deleteFile(file.FileName); // Deletion could be arrays, so possibly change this
                            break;
                    }
                }
            }
        });
    };
    ;
    Server.prototype.return = function () {
        var x = window.open("https://docs.codekeg.dev", '_top');
        x.focus();
    };
    return Server;
}());
// Initialize everything when DOM is loaded
var themeManager = ThemeManager.prototype;
var fileExplorer = FileExplorer.prototype;
var editor = Editor.prototype;
var server = Server.prototype;
document.addEventListener('DOMContentLoaded', function () {
    themeManager = new ThemeManager();
    fileExplorer = new FileExplorer();
    editor = new Editor();
});
//@ts-ignore
// Global functions
function addFile(fileName, dataExt, content) {
    if (fileExplorer) {
        fileExplorer.addFile(fileName, dataExt, content);
    }
}
function appendFileContent(filename, newContent) {
    fileExplorer.setFileContent(filename, fileExplorer.getFileContent(filename) + newContent);
}
function deleteFile(filename) {
    fileExplorer.removeFile(filename);
}

