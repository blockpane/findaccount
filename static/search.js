async function doSearch() {
    hideTable()
    pleaseWait()
    const addr = document.getElementById('inputAddress').value
    document.getElementById('tableDiv').hidden = false

    let data
    try {
        const response = await fetch("/q?addr=" + addr, {
            method: 'GET',
            mode: 'cors',
            cache: 'no-cache',
            credentials: 'same-origin',
            redirect: 'error',
            referrerPolicy: 'no-referrer'
        });
        data = await response.json()
    } catch (e) {
        setStatus(e)
        hideTable()
        return
    }
    if (typeof data.error !== 'undefined') {
        setStatus(data.error)
        hideTable()
    } else {
        clearStatus()
        showTable(data)
    }
}

function hideTable() {
    document.getElementById('tableDiv').hidden = true
    document.getElementById('tableDiv').innerHTML = ""
}

function showTable(data) {
    let rows = `
    <table class="table table-striped">
      <thead>
      <tr>
        <th scope="col">Chain</th>
        <th scope="col">Address</th>
        <th scope="col">Validator moniker</th>
        <th scope="col">Coins</th>
      </tr>
      </thead>
      <tbody>`
    data.forEach(row => {
        if (row.hasBalance === true) {
            rows += `
              <tr>
              <td><a href="${row.link}/account/${row.address}" target="_new">${cap(row.chain)}</a></td>
              <td>${row.address}</td>
              <td>${row.is_validator}</td>
              <td>${row.coins}</td>
              </tr>`
        }
    })
    rows += `</tbody>
    </table>`
    document.getElementById('tableDiv').innerHTML = rows
    document.getElementById('tableDiv').hidden = false
}

function setStatus(msg) {
    document.getElementById('status').innerText = msg
}

function clearStatus() {
    document.getElementById('status').innerText = " "
}

function pleaseWait() {
    document.getElementById('status').hidden = false
    document.getElementById('status').innerText = "Searching dozens of IBC enabled chains ... this can take a while, please be patient."
}

function cap(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}