const API = ""; // mismo host/puerto

async function fetchJSON(url, options = {}) {
  const res = await fetch(url, {
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
    ...options,
  });

  const text = await res.text();
  let data = null;
  try { data = text ? JSON.parse(text) : null; } catch { data = text; }

  if (!res.ok) {
    const msg = (data && data.error) ? data.error : (typeof data === "string" ? data : "Error");
    throw new Error(msg);
  }
  return data;
}

function money(n){
  const x = Number(n || 0);
  return x.toLocaleString("es-EC", { style:"currency", currency:"USD" });
}

function setMsg(id, message, isError=false){
  const el = document.getElementById(id);
  if(!el) return;
  el.textContent = message || "";
  el.className = isError ? "msg error" : "msg";
  el.style.display = message ? "block" : "none";
}

function setText(id, val){
  const el = document.getElementById(id);
  if(el) el.textContent = String(val);
}

function escapeHTML(s){
  return String(s ?? "").replace(/[&<>"']/g, m => ({
    "&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"
  }[m]));
}

/* ===================== PRODUCTOS ===================== */

async function loadProducts(){
  setMsg("msgProducts", "Cargando productos...");
  try{
    const list = await fetchJSON(`${API}/api/products`);
    renderProducts(list || []);
    setMsg("msgProducts", `Listo: ${list.length} producto(s).`);
  }catch(e){
    setMsg("msgProducts", e.message, true);
  }
}

function renderProducts(list){
  const tbody = document.getElementById("productsBody");
  if(!tbody) return;

  tbody.innerHTML = "";

  for(const p of list){
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${p.id}</td>
      <td>${escapeHTML(p.nombre)}</td>
      <td><span class="badge">${p.stock}</span></td>
      <td>${money(p.precio)}</td>
    `;
    tbody.appendChild(tr);
  }

  const total = list.length;
  const stockTotal = list.reduce((acc, p) => acc + Number(p.stock||0), 0);
  const valorAprox = list.reduce((acc, p) => acc + (Number(p.stock||0) * Number(p.precio||0)), 0);

  setText("kpiProdTotal", total);
  setText("kpiProdStock", stockTotal);
  setText("kpiProdValor", money(valorAprox));
}

async function onCreateProduct(e){
  e.preventDefault();

  const nombre = document.getElementById("pNombre").value.trim();
  const stock = Number(document.getElementById("pStock").value);
  const precio = Number(document.getElementById("pPrecio").value);

  if(!nombre || !Number.isFinite(stock) || !Number.isFinite(precio)){
    setMsg("msgCreateProduct", "Completa nombre, stock y precio correctamente.", true);
    return;
  }

  setMsg("msgCreateProduct", "Guardando...");

  try{
    await fetchJSON(`${API}/api/products`, {
      method: "POST",
      body: JSON.stringify({ nombre, stock, precio })
    });

    document.getElementById("pNombre").value = "";
    document.getElementById("pStock").value = "";
    document.getElementById("pPrecio").value = "";

    setMsg("msgCreateProduct", "Producto creado ✅");
    await loadProducts();
  }catch(e2){
    setMsg("msgCreateProduct", e2.message, true);
  }
}

/* ===================== CLIENTES ===================== */

async function loadClients(){
  setMsg("msgClients", "Cargando clientes...");
  try{
    const list = await fetchJSON(`${API}/api/clients`);
    renderClients(list || []);
    setMsg("msgClients", `Listo: ${list.length} cliente(s).`);
  }catch(e){
    setMsg("msgClients", e.message, true);
  }
}

function renderClients(list){
  const tbody = document.getElementById("clientsBody");
  if(!tbody) return;

  tbody.innerHTML = "";

  for(const c of list){
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${c.id}</td>
      <td>${escapeHTML(c.nombre)}</td>
      <td>${escapeHTML(c.cedula)}</td>
      <td>${escapeHTML(c.email)}</td>
    `;
    tbody.appendChild(tr);
  }
}

async function onCreateClient(e){
  e.preventDefault();

  const nombre = document.getElementById("cNombre").value.trim();
  const cedula = document.getElementById("cCedula").value.trim();
  const email  = document.getElementById("cEmail").value.trim();

  if(!nombre || !cedula || !email){
    setMsg("msgCreateClient", "Todos los campos son obligatorios.", true);
    return;
  }

  setMsg("msgCreateClient", "Guardando...");

  try{
    await fetchJSON(`${API}/api/clients`, {
      method: "POST",
      body: JSON.stringify({ nombre, cedula, email })
    });

    document.getElementById("cNombre").value = "";
    document.getElementById("cCedula").value = "";
    document.getElementById("cEmail").value = "";

    setMsg("msgCreateClient", "Cliente creado ✅");
    await loadClients();
  }catch(e2){
    setMsg("msgCreateClient", e2.message, true);
  }
}

/* ===================== VENTAS ===================== */

let PRODUCTS_CACHE = [];
let CLIENTS_CACHE = [];
let SALE_ITEMS = [];

async function loadSalesPageData(){
  try{
    const [clients, products] = await Promise.all([
      fetchJSON(`${API}/api/clients`),
      fetchJSON(`${API}/api/products`)
    ]);

    CLIENTS_CACHE = clients || [];
    PRODUCTS_CACHE = products || [];

    fillClientsSelect(CLIENTS_CACHE);
    fillProductsSelect(PRODUCTS_CACHE);
  }catch(e){
    setMsg("msgSale", e.message, true);
  }
}

function fillClientsSelect(list){
  const sel = document.getElementById("saleClient");
  if(!sel) return;

  sel.innerHTML = "";
  const opt0 = document.createElement("option");
  opt0.value = "";
  opt0.textContent = "Seleccione cliente...";
  sel.appendChild(opt0);

  for(const c of list){
    const opt = document.createElement("option");
    opt.value = String(c.id);
    opt.textContent = `${c.nombre} • ${c.cedula}`;
    sel.appendChild(opt);
  }
}

function fillProductsSelect(list){
  const sel = document.getElementById("saleProduct");
  if(!sel) return;

  sel.innerHTML = "";
  const opt0 = document.createElement("option");
  opt0.value = "";
  opt0.textContent = "Seleccione producto...";
  sel.appendChild(opt0);

  for(const p of list){
    const opt = document.createElement("option");
    opt.value = String(p.id);
    opt.textContent = `${p.nombre} • stock:${p.stock} • ${money(p.precio)}`;
    sel.appendChild(opt);
  }
}

function recalcSale(){
  SALE_ITEMS.forEach(it => {
    it.subtotal = Number(it.cantidad) * Number(it.precio_unitario);
  });
  const total = SALE_ITEMS.reduce((a, it) => a + Number(it.subtotal||0), 0);
  const totalEl = document.getElementById("saleTotal");
  if(totalEl) totalEl.textContent = money(total);
  renderSaleItems();
}

function renderSaleItems(){
  const tbody = document.getElementById("saleItemsBody");
  if(!tbody) return;
  tbody.innerHTML = "";

  for(let i=0;i<SALE_ITEMS.length;i++){
    const it = SALE_ITEMS[i];
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${escapeHTML(it.nombre)}</td>
      <td>${it.cantidad}</td>
      <td>${money(it.precio_unitario)}</td>
      <td>${money(it.subtotal)}</td>
      <td><button data-i="${i}" type="button">Quitar</button></td>
    `;
    tbody.appendChild(tr);
  }

  tbody.querySelectorAll("button[data-i]").forEach(btn => {
    btn.addEventListener("click", () => {
      const idx = Number(btn.getAttribute("data-i"));
      SALE_ITEMS.splice(idx, 1);
      recalcSale();
    });
  });
}

function addSaleItem(){
  const selProd = document.getElementById("saleProduct");
  const qtyEl = document.getElementById("saleQty");
  if(!selProd || !qtyEl) return;

  const productID = Number(selProd.value);
  const cantidad = Number(qtyEl.value);

  if(!productID || !cantidad) return;

  const p = PRODUCTS_CACHE.find(x => Number(x.id) === productID);
  if(!p) return;

  // si ya existe, acumula cantidad
  const ex = SALE_ITEMS.find(it => it.product_id === productID);
  if(ex){
    ex.cantidad += cantidad;
  }else{
    SALE_ITEMS.push({
      product_id: productID,
      nombre: p.nombre,
      cantidad,
      precio_unitario: Number(p.precio),
      subtotal: 0
    });
  }

  qtyEl.value = "";
  recalcSale();
}

async function confirmSale(){
  const selClient = document.getElementById("saleClient");
  if(!selClient) return;

  const clientID = Number(selClient.value);
  if(!clientID || !SALE_ITEMS.length){
    setMsg("msgSale", "Selecciona cliente y agrega productos.", true);
    return;
  }

  const payload = {
    client_id: clientID,
    items: SALE_ITEMS.map(it => ({
      product_id: it.product_id,
      cantidad: it.cantidad,
      precio_unitario: it.precio_unitario
    }))
  };

  setMsg("msgSale", "Registrando venta...");

  try{
    await fetchJSON(`${API}/api/sales`, {
      method: "POST",
      body: JSON.stringify(payload)
    });

    setMsg("msgSale", "Venta creada ✅");
    SALE_ITEMS = [];
    recalcSale();

    // refrescar stock y lista
    PRODUCTS_CACHE = await fetchJSON(`${API}/api/products`);
    fillProductsSelect(PRODUCTS_CACHE);
    await loadSalesList();
  }catch(e){
    setMsg("msgSale", e.message, true);
  }
}

async function loadSalesList(){
  try{
    const list = await fetchJSON(`${API}/api/sales`);
    renderSales(list || []);
  }catch(e){
    setMsg("msgSalesList", e.message, true);
  }
}

function renderSales(list){
  const tbody = document.getElementById("salesBody");
  if(!tbody) return;

  tbody.innerHTML = "";

  for(const s of list){
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${s.id}</td>
      <td>${escapeHTML(s.client_name || ("ID " + s.client_id))}</td>
      <td>${formatDate(s.fecha)}</td>
      <td>${money(s.total)}</td>
      <td><button class="btn secondary" data-sale="${s.id}" type="button">Ver</button></td>
    `;
    tbody.appendChild(tr);
  }

  tbody.querySelectorAll("button[data-sale]").forEach(btn => {
    btn.addEventListener("click", async () => {
      const id = Number(btn.getAttribute("data-sale"));
      await openSaleDetail(id);
    });
  });
}

function formatDate(x){
  if(!x) return "-";
  try{
    const d = new Date(x);
    if(Number.isNaN(d.getTime())) return String(x);
    return d.toLocaleString("es-EC");
  }catch{
    return String(x);
  }
}

/* ===================== DETALLE DE VENTA ===================== */

async function openSaleDetail(id){
  const box = document.getElementById("saleDetailBox");
  const meta = document.getElementById("saleDetailMeta");
  const tbody = document.getElementById("saleDetailItems");
  if(!box || !meta || !tbody) return;

  box.style.display = "block";
  meta.textContent = "Cargando detalle...";
  tbody.innerHTML = "";

  try{
    const s = await fetchJSON(`${API}/api/sales/${id}`);
    meta.textContent = `Venta #${s.id} • Cliente ID: ${s.client_id} • Fecha: ${formatDate(s.fecha)} • Total: ${money(s.total)}`;

    for(const it of (s.items || [])){
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td>${it.product_id}</td>
        <td><span class="badge">${it.cantidad}</span></td>
        <td>${money(it.precio_unitario)}</td>
        <td>${money(it.subtotal)}</td>
      `;
      tbody.appendChild(tr);
    }
  }catch(e){
    meta.textContent = `Error: ${e.message}`;
  }
}

/* ===================== REPORTE VENTAS DEL DÍA ===================== */

async function loadVentasHoy(){
  const box = document.getElementById("ventasHoyBox");
  if(!box) return;

  box.style.display = "block";
  box.className = "msg";
  box.textContent = "Cargando reporte del día...";

  try{
    const r = await fetchJSON(`${API}/api/report/ventas-hoy`);
    box.className = "msg";
    box.textContent = `Hoy: ${r.ventas} venta(s) • Total: ${money(r.total)}`;
  }catch(e){
    box.className = "msg error";
    box.textContent = e.message;
  }
}

/* ===================== INIT ===================== */

document.addEventListener("DOMContentLoaded", () => {

  const formProduct = document.getElementById("formCreateProduct");
  if(formProduct){
    formProduct.addEventListener("submit", onCreateProduct);
    loadProducts();
  }

  const formClient = document.getElementById("formCreateClient");
  if(formClient){
    formClient.addEventListener("submit", onCreateClient);
    loadClients();
  }

  const btnAddItem = document.getElementById("btnAddItem");
  const btnConfirmSale = document.getElementById("btnConfirmSale");
  const btnCloseDetail = document.getElementById("btnCloseDetail");
  const btnVentasHoy = document.getElementById("btnVentasHoy");

  // Sales page hooks
  if(btnAddItem && btnConfirmSale){
    btnAddItem.addEventListener("click", addSaleItem);
    btnConfirmSale.addEventListener("click", confirmSale);
    loadSalesPageData();
    loadSalesList();
    recalcSale();
  }

  if(btnCloseDetail){
    btnCloseDetail.addEventListener("click", () => {
      const box = document.getElementById("saleDetailBox");
      if(box) box.style.display = "none";
    });
  }

  if(btnVentasHoy){
    btnVentasHoy.addEventListener("click", loadVentasHoy);
  }

});