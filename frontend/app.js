// -------------------------------------------------------
// frontend/app.js
// -------------------------------------------------------
// Purpose Summary:
//   - Load folders, manage file tabs, save/load content.
//   - Adds "+ New File" button with safe creation flow.
//   - Restricts to 10 open tabs. Logs all user actions.
//   - Adds global search, file rename/delete, keyboard shortcuts.
// Audit:
//   - All actions log to console with UTC ISO 8601.
//   - All backend requests validated; hard-fail on non-2xx.
//   - Right pane remains empty until a file is opened.
// -------------------------------------------------------

const MAX_TABS = 10;
// const API_BASE = "http://localhost:8888";
const API_BASE = "";

let activeFolder = null;
const tabs = {};
let currentTab = null;

// -------------------------------------------------------
// function utcNow()
// -------------------------------------------------------
// Purpose:
//   - Return current UTC time in ISO 8601 format.
// Audit:
//   - Used for all console log timestamps.
// -------------------------------------------------------
function utcNow() {
    return new Date().toISOString();
}

// -------------------------------------------------------
// function log(level, msg)
// -------------------------------------------------------
// Purpose:
//   - Structured console logging with UTC timestamps.
// Audit:
//   - Levels: [INFO], [WARN], [ERROR], [DEBUG].
// -------------------------------------------------------
function log(level, msg) {
    console.log(`[${level}] ${utcNow()} ${msg}`);
}

// -------------------------------------------------------
// function requireOk(res)
// -------------------------------------------------------
// Purpose:
//   - Enforce HTTP 2xx status before parsing responses.
// Audit:
//   - Throws with explicit HTTP status; callers must catch.
// -------------------------------------------------------
function requireOk(res) {
    if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
    }
    return res;
}

// -------------------------------------------------------
// function asArray(x)
// -------------------------------------------------------
// Purpose:
//   - Normalize possibly-null JSON into an array.
// Audit:
//   - Prevents runtime errors on unexpected null.
// -------------------------------------------------------
function asArray(x) {
    return Array.isArray(x) ? x : [];
}

document.addEventListener("DOMContentLoaded", () => {
    // Ensure New File button exists (create if missing in HTML)
    let newFileBtn = document.getElementById("new-file-btn");
    const newFolderBtn = document.getElementById("new-folder-btn");
    if (!newFileBtn) {
        newFileBtn = document.createElement("button");
        newFileBtn.id = "new-file-btn";
        newFileBtn.textContent = "+ New File";
        newFileBtn.style.display = "none";
        newFolderBtn.insertAdjacentElement("afterend", newFileBtn);
    }

    // Wire New File creation
    newFileBtn.addEventListener("click", () => {
        createNewFile();
    });

    // Initial folder load hides the New File button until a folder is selected
    document.getElementById("new-file-btn").style.display = "none";
    loadFolders();

    document.getElementById("new-folder-btn").addEventListener("click", () => {
        const name = prompt("New folder name:");
        if (!name) return;
        fetch(`${API_BASE}/folders`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name })
        })
        .then(requireOk)
        .then(() => {
            log("INFO", "Created folder: " + name);
            loadFolders();
        })
        .catch(err => {
            log("ERROR", "Create folder failed: " + err.message);
            alert("Failed to create folder");
        });
    });

    document.getElementById("search-box").addEventListener("input", event => {
        const query = event.target.value.trim();
        if (query.length >= 3) {
            searchFiles(query);
        } else {
            document.getElementById("search-results").innerHTML = "";
        }
    });

    document.addEventListener("keydown", event => {
        if (event.ctrlKey && event.key === "s") {
            event.preventDefault();
            if (currentTab) {
                const editors = document.querySelectorAll("#editor-area textarea");
                for (let i = 0; i < Object.keys(tabs).length; i++) {
                    if (Object.keys(tabs)[i] === currentTab) {
                        saveFile(currentTab, editors[i].value);
                        break;
                    }
                }
            }
        }

        if (event.ctrlKey && event.key === "Tab") {
            event.preventDefault();
            const keys = Object.keys(tabs);
            const idx = keys.indexOf(currentTab);
            if (idx >= 0) {
                const nextIdx = (idx + 1) % keys.length;
                switchTab(keys[nextIdx]);
            }
        }
    });
});

// -------------------------------------------------------
// function loadFolders()
// -------------------------------------------------------
// Purpose:
//   - Loads folder tree and populates sidebar.
// Audit:
//   - Validates response; initializes UI state safely.
// -------------------------------------------------------
function loadFolders() {
    fetch(`${API_BASE}/folders`)
        .then(requireOk)
        .then(res => res.json())
        .then(folders => {
            folders = asArray(folders);
            const list = document.getElementById("folder-list");
            list.innerHTML = "";

            folders.forEach(folder => {
                const details = document.createElement("details");
                details.className = "folder";

                const summary = document.createElement("summary");
                summary.textContent = folder;
                details.appendChild(summary);

                const filesUl = document.createElement("ul");
                filesUl.className = "file-list";
                details.appendChild(filesUl);

                // Lazy-load files the first time this folder is opened
                details.addEventListener("toggle", () => {
                    if (details.open && !details.dataset.loaded) {
                        activeFolder = folder;
                        document.getElementById("new-file-btn").style.display = "inline-block";
                        loadFilesInto(details, folder);
                    }
                });

                list.appendChild(details);
            });

            // No active folder until user expands one
            activeFolder = null;
            document.getElementById("new-file-btn").style.display = "none";
            log("INFO", "Loaded folders");
        })
        .catch(err => {
            log("ERROR", "Failed to load folders: " + err.message);
            alert("Folder load failed");
        });
}

// New helper; do NOT show confirm()s—just list files
function loadFilesInto(detailsEl, folder) {
    const filesUl = detailsEl.querySelector(".file-list");
    filesUl.innerHTML = "<li>(Loading...)</li>";

    fetch(`${API_BASE}/files?folder=${encodeURIComponent(folder)}`)
        .then(requireOk)
        .then(res => res.json())
        .then(files => {
            files = asArray(files);
            filesUl.innerHTML = "";

            if (files.length === 0) {
                const empty = document.createElement("li");
                empty.textContent = "(empty)";
                empty.className = "muted";
                filesUl.appendChild(empty);
            } else {
                files.forEach(file => {
                    const li = document.createElement("li");
                    li.className = "file";
                    li.textContent = file;
                    li.title = `${folder}/${file}`;
                    li.addEventListener("click", (e) => {
                        e.stopPropagation(); // don't collapse the folder
                        openFile(`${folder}/${file}`);
                    });
                    filesUl.appendChild(li);
                });
            }
            detailsEl.dataset.loaded = "1";
            log("INFO", `Loaded ${files.length} files from ${folder}`);
        })
        .catch(err => {
            log("ERROR", "Failed to load files: " + err.message);
            filesUl.innerHTML = "<li class='error'>(Failed to load)</li>";
        });
}

// -------------------------------------------------------
// function loadFiles(folder)
// -------------------------------------------------------
// Purpose:
//   - Fetch list of files for selected folder; enable New File button.
// Audit:
//   - Validates responses; logs actions; fails visibly.
// -------------------------------------------------------
function loadFiles(folder) {
    activeFolder = folder;
    document.getElementById("new-file-btn").style.display = "inline-block";

    fetch(`${API_BASE}/files?folder=${encodeURIComponent(folder)}`)
        .then(requireOk)
        .then(res => res.json())
        .then(files => {
            files = asArray(files);
            if (files.length === 0) {
                log("INFO", `No files in ${folder}`);
                return;
            }
            files.forEach(file => {
                const open = confirm(`Open ${file}?`);
                if (open) openFile(`${folder}/${file}`);
            });
            log("INFO", `Loaded ${files.length} files from ${folder}`);
        })
        .catch(err => {
            log("ERROR", "Failed to load files: " + err.message);
            alert("File load failed");
        });
}

// -------------------------------------------------------
// function createNewFile()
// -------------------------------------------------------
// Purpose:
//   - Create a new uniquely named .txt file in current folder.
//   - Uses -NN suffix to avoid accidental duplicates.
// Audit:
//   - Validates folder selection; checks existing files before naming.
// -------------------------------------------------------
function createNewFile() {
    if (!activeFolder) {
        alert("Select a folder first");
        log("WARN", "New File requested without active folder");
        return;
    }
    let baseName = prompt("Base file name (without extension):", "notes");
    if (!baseName) return;

    fetch(`${API_BASE}/files?folder=${encodeURIComponent(activeFolder)}`)
        .then(requireOk)
        .then(res => res.json())
        .then(files => {
            files = asArray(files);
            let counter = 1;
            let finalName;
            do {
                const suffix = String(counter).padStart(2, "0");
                finalName = `${baseName}-${suffix}.txt`;
                counter++;
            } while (files.includes(finalName));

            const path = `${activeFolder}/${finalName}`;
            fetch(`${API_BASE}/file/save`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ path, content: "" })
            })
            .then(requireOk)
            .then(() => {
                log("INFO", "Created file: " + path);
                openFile(path);
            })
            .catch(err => {
                log("ERROR", "Create file failed: " + err.message);
                alert("Failed to create file");
            });
        })
        .catch(err => {
            log("ERROR", "Failed to check existing files: " + err.message);
            alert("Failed to check for duplicates");
        });
}

// -------------------------------------------------------
// function openFile(path)
// -------------------------------------------------------
// Purpose:
//   - Opens a .txt file in a new tab, if under limit.
// Audit:
//   - Validates response; logs open; prevents duplicates.
// -------------------------------------------------------
function openFile(path) {
    if (tabs[path]) {
        switchTab(path);
        return;
    }
    if (Object.keys(tabs).length >= MAX_TABS) {
        alert("Max 10 tabs reached");
        log("WARN", "Tab limit reached");
        return;
    }

    fetch(`${API_BASE}/file?path=${encodeURIComponent(path)}`)
        .then(requireOk)
        .then(res => res.text())
        .then(content => {
            tabs[path] = content;
            renderTab(path, content);
            switchTab(path);
            log("INFO", "Opened file: " + path);
        })
        .catch(err => {
            log("ERROR", "Open file failed: " + err.message);
        });
}

// -------------------------------------------------------
// function renderTab(path, content)
// -------------------------------------------------------
// Purpose:
//   - Adds a new tab and editor textarea to DOM.
//   - Includes clickable "X" to close tab without deleting file.
// Audit:
//   - Logs DOM mutations; no silent failures.
//   - Prevents accidental file deletion when closing tab.
// -------------------------------------------------------
function renderTab(path, content) {
    const tab = document.createElement("div");
    tab.className = "tab";
    tab.title = path;

    const nameSpan = document.createElement("span");
    nameSpan.textContent = path.split("/").pop();
    nameSpan.addEventListener("click", () => switchTab(path));

    const closeBtn = document.createElement("span");
    closeBtn.textContent = " ✕";
    closeBtn.style.cursor = "pointer";
    closeBtn.style.marginLeft = "6px";
    closeBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        closeTab(path);
    });

    tab.appendChild(nameSpan);
    tab.appendChild(closeBtn);
    document.getElementById("tabs").appendChild(tab);

    const textarea = document.createElement("textarea");
    textarea.value = content;
    textarea.style.display = "none";
    document.getElementById("editor-area").appendChild(textarea);

    log("DEBUG", "Rendered tab and editor for: " + path);
}

// -------------------------------------------------------
// function closeTab(path)
// -------------------------------------------------------
// Purpose:
//   - Closes tab + editor for given path without touching disk.
// Audit:
//   - Logs closure; switches to first available tab if any remain.
// -------------------------------------------------------
function closeTab(path) {
    const keys = Object.keys(tabs);
    const idx = keys.indexOf(path);
    if (idx >= 0) {
        delete tabs[path];
        const tabsDom = document.getElementById("tabs").children;
        const editorsDom = document.getElementById("editor-area").children;
        tabsDom[idx].remove();
        editorsDom[idx].remove();
        log("INFO", "Closed tab: " + path);
        currentTab = null;
        if (Object.keys(tabs).length > 0) {
            switchTab(Object.keys(tabs)[0]);
        }
    }
}

// -------------------------------------------------------
// function switchTab(path)
// -------------------------------------------------------
// Purpose:
//   - Activates tab and editor for given path.
// Audit:
//   - Logs tab switches; hides non-active editors.
// -------------------------------------------------------
function switchTab(path) {
    currentTab = path;
    document.querySelectorAll(".tab").forEach(t => t.classList.remove("active"));
    document.querySelectorAll("textarea").forEach(t => t.style.display = "none");

    const tabsDom = document.getElementById("tabs").children;
    const editorsDom = document.getElementById("editor-area").children;

    for (let i = 0; i < tabsDom.length; i++) {
        if (Object.keys(tabs)[i] === path) {
            tabsDom[i].classList.add("active");
            editorsDom[i].style.display = "block";
            break;
        }
    }

    log("INFO", "Switched tab to: " + path);
}

// -------------------------------------------------------
// function saveFile(path, content)
// -------------------------------------------------------
// Purpose:
//   - Saves file to backend on content change.
// Audit:
//   - Logs before/after snapshot (truncated) and save success.
// -------------------------------------------------------
function saveFile(path, content) {
    fetch(`${API_BASE}/file/save`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ path, content })
    })
    .then(requireOk)
    .then(() => {
        log("INFO", "Saved file: " + path);
        log("DEBUG", "Snapshot: " + (content.slice(0, 200) + (content.length > 200 ? " [truncated]" : "")));
    })
    .catch(err => {
        log("ERROR", "Save failed: " + err.message);
    });
}

// -------------------------------------------------------
// function renameFile(oldPath)
// -------------------------------------------------------
// Purpose:
//   - Renames a file using the backend move API.
// Audit:
//   - Logs old/new paths; updates in-memory tab map.
// -------------------------------------------------------
function renameFile(oldPath) {
    const base = oldPath.substring(0, oldPath.lastIndexOf("/"));
    const oldName = oldPath.split("/").pop();
    const newName = prompt("New file name:", oldName);
    if (!newName || newName === oldName) return;

    const newPath = `${base}/${newName}`;
    fetch(`${API_BASE}/file/move`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ from: oldPath, to: newPath })
    })
    .then(requireOk)
    .then(() => {
        const content = tabs[oldPath];
        delete tabs[oldPath];
        tabs[newPath] = content;
        log("INFO", `Renamed: ${oldPath} -> ${newPath}`);
        // Update UI without full reload
        const tabsDom = document.getElementById("tabs").children;
        for (let i = 0; i < tabsDom.length; i++) {
            if (tabsDom[i].title === oldPath) {
                tabsDom[i].title = newPath;
                tabsDom[i].textContent = newPath.split("/").pop();
                break;
            }
        }
        switchTab(newPath);
    })
    .catch(err => {
        log("ERROR", "Rename failed: " + err.message);
        alert("Rename failed");
    });
}

// -------------------------------------------------------
// function deleteFile(path)
// -------------------------------------------------------
// Purpose:
//   - Permanently deletes the specified file.
// Audit:
//   - Removes tab from memory; updates UI without full reload.
// -------------------------------------------------------
function deleteFile(path) {
    const confirmDelete = confirm(`Delete file ${path}?`);
    if (!confirmDelete) return;

    fetch(`${API_BASE}/file?path=${encodeURIComponent(path)}`, {
        method: "DELETE"
    })
    .then(requireOk)
    .then(() => {
        delete tabs[path];
        log("INFO", "Deleted file: " + path);
        // Remove tab/editor DOM nodes
        const tabsDom = document.getElementById("tabs").children;
        const editorsDom = document.getElementById("editor-area").children;
        for (let i = 0; i < tabsDom.length; i++) {
            if (tabsDom[i].title === path) {
                tabsDom[i].remove();
                editorsDom[i].remove();
                break;
            }
        }
        currentTab = null;
    })
    .catch(err => {
        log("ERROR", "Delete failed: " + err.message);
        alert("Delete failed");
    });
}

// -------------------------------------------------------
// function searchFiles(query)
// -------------------------------------------------------
// Purpose:
//   - Client-side full-text search across all files.
// Audit:
//   - Validates responses; logs queries and result counts.
// -------------------------------------------------------
function searchFiles(query) {
    fetch(`${API_BASE}/folders`)
        .then(requireOk)
        .then(res => res.json())
        .then(folders => {
            folders = asArray(folders);
            const allPromises = folders.map(folder =>
                fetch(`${API_BASE}/files?folder=${encodeURIComponent(folder)}`)
                    .then(requireOk)
                    .then(res => res.json())
                    .then(files => ({ folder, files: asArray(files) }))
            );

            Promise.all(allPromises)
                .then(fileGroups => {
                    const fetches = [];
                    fileGroups.forEach(group => {
                        group.files.forEach(file => {
                            const path = `${group.folder}/${file}`;
                            fetches.push(
                                fetch(`${API_BASE}/file?path=${encodeURIComponent(path)}`)
                                    .then(requireOk)
                                    .then(res => res.text())
                                    .then(text => ({ path, content: text }))
                            );
                        });
                    });

                    Promise.all(fetches).then(allFiles => {
                        const results = allFiles.filter(f =>
                            f.content.toLowerCase().includes(query.toLowerCase())
                        );
                        showSearchResults(results, query);
                        log("INFO", `Search '${query}' found ${results.length} result(s)`);
                    });
                });
        })
        .catch(err => {
            log("ERROR", "Search failed: " + err.message);
        });
}

// -------------------------------------------------------
// function showSearchResults(results, query)
// -------------------------------------------------------
// Purpose:
//   - Displays clickable list of matched files from search.
// Audit:
//   - Logs selection; clears UI on open.
// -------------------------------------------------------
function showSearchResults(results, query) {
    const list = document.getElementById("search-results");
    list.innerHTML = "";

    results.forEach(result => {
        const item = document.createElement("li");
        item.textContent = result.path;
        item.addEventListener("click", () => {
            openFile(result.path);
            document.getElementById("search-box").value = "";
            list.innerHTML = "";
            log("INFO", `Opened search match: ${result.path}`);
        });
        list.appendChild(item);
    });
}
