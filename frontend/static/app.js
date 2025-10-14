async function fetchTime() {
    const res = await fetch('/api/time');
    if (!res.ok) throw new Error('HTTP ' + res.status);
    const data = await res.json();
    const isoEl = document.getElementById('iso');
    const zoneEl = document.getElementById('zone');
    const unixEl = document.getElementById('unix');
    const localEl = document.getElementById('local');


    isoEl.textContent = data.iso;
    zoneEl.textContent = `Zone (Backend): ${data.zone}`;
    unixEl.textContent = String(data.unix);
    try {
        const d = new Date(data.iso);
        localEl.textContent = d.toLocaleString();
    } catch { localEl.textContent = '—'; }
}


window.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById('refresh');
    btn.addEventListener('click', () => fetchTime().catch(console.error));
    fetchTime().catch(console.error);
});
