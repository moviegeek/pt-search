const data = require('./cookies.json')
const fs = require('fs').promises

async function fix() {
    fixCapital()
    fixExpire()
    fixCanonicalHost()
    fixPersistent()
    fixTime()

    const fixedData = JSON.stringify(data, null, ' ');
    await fs.writeFile('cookies-fixed.json', fixedData);
}

function fixCapital() {
    function capitalize(s) {
        return s[0].toUpperCase() + s.substring(1);
    }

    data.forEach((entry, i) => {
        let newEntry = Object.fromEntries(
            Object.entries(entry).map(([k, v]) => [capitalize(k), v])
        )
        data[i] = newEntry
    })
}

function fixExpire() {
    data.forEach(entry => {
        if (entry.Session === true || entry.Expires < 0) {
            entry.Persistent = false
            entry.Expires = new Date('9999-12-31T23:59:59Z')
            return
        }
        let d = new Date(entry.Expires * 1000);
        entry.Expires = d.toISOString();
    });
}

function fixCanonicalHost() {
    data.forEach(entry => {
        if (entry.CanonicalHost) return;
        let d = entry.Domain;
        if (d.startsWith('.')) d = d.substring(1)
        entry.CanonicalHost = d;
    });
}

function fixPersistent() {
    data.forEach(entry => {
        if (entry.Persistent !== undefined) return;
        entry.Persistent = true;
    });
}

function fixTime() {
    let now = new Date().toISOString()
    data.forEach(entry => {
        if (!entry.Creation) entry.Creation = now;
        if (!entry.LastAccess) entry.LastAccess = now;
        if (!entry.Updated) entry.Updated = now;
    });
}

fix().then(
    () => console.log('done'),
    e => console.log(e)
)