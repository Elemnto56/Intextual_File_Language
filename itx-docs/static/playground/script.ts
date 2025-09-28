// ===================================
// Code Playground TypeScript
// ===================================
// @ts-ignore
// Theme Management
class ThemeManager {
    private currentTheme: 'light' | 'dark' = 'light';
    private themeToggleBtn: HTMLElement | null;

    constructor() {
        this.themeToggleBtn = document.querySelector('.theme-toggle');
        this.init();
    }

    private init(): void {
        // Check for saved theme preference or default to 'light'
        const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null;
        if (savedTheme) {
            this.currentTheme = savedTheme;
        } else {
            // Check system preference
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                this.currentTheme = 'dark';
            }
        }

        // Apply the theme
        this.applyTheme(this.currentTheme);

        // Add click listener to toggle button
        if (this.themeToggleBtn) {
            this.themeToggleBtn.addEventListener('click', () => this.toggleTheme());
        }
    }

    private applyTheme(theme: 'light' | 'dark'): void {
        if (theme === 'dark') {
            document.body.classList.add('dark-theme');
        } else {
            document.body.classList.remove('dark-theme');
        }
        this.currentTheme = theme;
        localStorage.setItem('theme', theme);
    }

    private toggleTheme(): void {
        const newTheme = this.currentTheme === 'light' ? 'dark' : 'light';
        this.applyTheme(newTheme);
    }

    public getTheme(): 'light' | 'dark' {
        return this.currentTheme;
    }
}
// File Explorer Management
class FileExplorer {
    private fileTree: HTMLElement | null;
    private selectedFile: HTMLElement | null = null;
    private fileContents: Map<string, string> = new Map();

    constructor() {
        this.fileTree = document.getElementById('file-tree');
        this.init();
    }

    private init(): void {
        // Add click listener for file selection
        if (this.fileTree) {
            this.fileTree.addEventListener('click', (e: MouseEvent) => {
                const target = e.target as HTMLElement;
                const fileElement = target.closest('.file') as HTMLElement;
                const folderElement = target.closest('.folder') as HTMLElement;

                if (fileElement) {
                    this.selectFile(fileElement);
                } else if (folderElement && target.closest('.folder')?.firstChild === target) {
                    // Only toggle if clicking on the folder itself, not its content
                    this.toggleFolder(folderElement);
                }
            });
        }
    }

    private selectFile(fileElement: HTMLElement): void {
        // Remove previous selection
        if (this.selectedFile) {
            this.selectedFile.classList.remove('selected');
        }

        // Add selection to clicked file
        fileElement.classList.add('selected');
        this.selectedFile = fileElement;

        const fileName = fileElement.textContent?.trim() || '';
        const fileExt = fileElement.getAttribute('data-ext') || '';
        
        // Load file content into the code input if it exists
        const content = this.fileContents.get(fileName);
        if (content !== undefined && editor) {
            editor.setInput(content);
        }

        // Trigger custom event for file selection
        const event = new CustomEvent('fileSelected', {
            detail: { fileName, fileExt, content }
        });
        document.dispatchEvent(event);
    }

    private toggleFolder(folderElement: HTMLElement): void {
        folderElement.classList.toggle('open');
    }

    /**
     * Add a file to the file tree
     * @param fileName - The name of the file to add
     * @param dataExt - Optional data extension attribute
     * @param content - Optional file content
     */
    public addFile(fileName: string, dataExt?: string, content?: string): void {
        if (!this.fileTree) return;

        const fileDiv = document.createElement('div');
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
    }

    public removeFile(filename: string) {
        let possibleFile = document.querySelector("#file-tree > *");

        if (possibleFile?.innerHTML == filename) {
            document.querySelector("#file-tree")?.removeChild(possibleFile);
        } else {
            console.error("The file set for deletion did not exist");
        }
    }
    /**
     * Set or update file content
     * @param fileName - The name of the file
     * @param content - The content to set
     */
    public setFileContent(fileName: string, content: string): void {
        this.fileContents.set(fileName, content);
    }

    /**
     * Get file content
     * @param fileName - The name of the file
     * @returns The file content or undefined if not found
     */
    public getFileContent(fileName: string): string | undefined {
        return this.fileContents.get(fileName);
    }

    /**
     * Update content of currently selected file with current editor input
     */
    public updateSelectedFileContent(): void {
        if (this.selectedFile && editor) {
            const fileName = this.selectedFile.textContent?.trim() || '';
            const currentContent = editor.getInput();
            this.fileContents.set(fileName, currentContent);
        }
    }

    /**
     * Get all files and their contents
     */
    public getAllFiles(): Map<string, string> {
        return new Map(this.fileContents);
    }

    /**
     * Add a folder to the file tree
     * @param folderName - The name of the folder
     * @returns The folder content div where child items can be added
     */
    public addFolder(folderName: string): HTMLElement | null {
        if (!this.fileTree) return null;

        const folderDiv = document.createElement('div');
        folderDiv.className = 'folder';
        folderDiv.textContent = folderName;

        const folderContent = document.createElement('div');
        folderContent.className = 'folder-content';

        folderDiv.appendChild(folderContent);
        this.fileTree.appendChild(folderDiv);

        return folderContent;
    }

    /**
     * Clear all files from the file tree
     */
    public clearFiles(): void {
        if (this.fileTree) {
            this.fileTree.innerHTML = '';
        }
        this.selectedFile = null;
    }

    /**
     * Get the currently selected file
     */
    public getSelectedFile(): { fileName: string; fileExt: string } | null {
        if (!this.selectedFile) return null;

        return {
            fileName: this.selectedFile.textContent?.trim() || '',
            fileExt: this.selectedFile.getAttribute('data-ext') || ''
        };
    }
}

// Editor Management
class Editor {
    private codeInput: HTMLTextAreaElement | null;
    private codeOutput: HTMLTextAreaElement | null;

    constructor() {
        this.codeInput = document.getElementById('code-input') as HTMLTextAreaElement;
        this.codeOutput = document.getElementById('code-output') as HTMLTextAreaElement;
        this.init();
    }

    private init(): void {
        // Add tab support in the input textarea
        if (this.codeInput) {
            this.codeInput.addEventListener('keydown', (e: KeyboardEvent) => {
                if (e.key === 'Tab') {
                    e.preventDefault();
                    const start = this.codeInput!.selectionStart;
                    const end = this.codeInput!.selectionEnd;
                    const value = this.codeInput!.value;

                    // Insert tab at cursor
                    this.codeInput!.value = value.substring(0, start) + '\t' + value.substring(end);
                    
                    // Move cursor after the tab
                    this.codeInput!.selectionStart = this.codeInput!.selectionEnd = start + 1;
                }

                // Ctrl/Cmd + Enter to run
                if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
                    e.preventDefault();
                    this.run();
                }
            });
        }
    }

    private run(): void {
        // Call the Run function that should be defined elsewhere
        if (typeof (window as any).Run === 'function') {
            (window as any).Run();
        }
    }

    public getInput(): string {
        return this.codeInput?.value || '';
    }

    public setInput(value: string): void {
        if (this.codeInput) {
            this.codeInput.value = value;
        }
    }

    public getOutput(): string {
        return this.codeOutput?.value || '';
    }

    public setOutput(value: string): void {
        if (this.codeOutput) {
            this.codeOutput.value = value;
        }
    }

    public clearInput(): void {
        this.setInput('');
    }

    public clearOutput(): void {
        this.setOutput('');
    }
}

// Server Management
class Server {
    public runCode() {
        const codeInput = document.getElementById("code-input") as HTMLTextAreaElement;
        const codeOutput = document.getElementById("code-output") as HTMLTextAreaElement;

        fetch("https://api.codekeg.dev/run-code", {
            method: "POST",
            headers:  {"Content-Type": "application/json"},
            body: JSON.stringify({ data: codeInput.value })
        })
        .then(res => res.json())
        .then(response => { 
            codeOutput.innerHTML = response.return;

            for (var _i = 0, _a = response.filedata; _i < _a.length; _i++) {
                var file = _a[_i];
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
            console.debug(response);
        });
    };

    public return() {
       const x = window.open("https://docs.codekeg.dev", '_top') as Window; 
       x.focus();
    }
}

// Initialize everything when DOM is loaded
let themeManager = ThemeManager.prototype;
let fileExplorer = FileExplorer.prototype;
let editor = Editor.prototype;
let server = Server.prototype;

document.addEventListener('DOMContentLoaded', () => {
    themeManager = new ThemeManager();
    fileExplorer = new FileExplorer();
    editor = new Editor();
});
//@ts-ignore
// Global functions
function addFile(fileName: string, dataExt?: string, content: string): void {
    if (fileExplorer) {
        fileExplorer.addFile(fileName, dataExt, content);
    }
}

function appendFileContent(filename: string, newContent: string) {
    fileExplorer.setFileContent(filename, fileExplorer.getFileContent(filename) + newContent);
}

function deleteFile(filename: string) {
    fileExplorer.removeFile(filename);
}