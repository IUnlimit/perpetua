/* ═══════════════════════════════════════════════
   perpetua dashboard — application logic
   ═══════════════════════════════════════════════ */

const API_BASE = '/api/web';
const POLL_INTERVAL = 5000;
const PAGE_SIZE = 30;

// ── State ──
let currentLinkFilter = 'all';       // 'all' | 'ntqq' | 'client'
let currentDirFilter = 'all';        // 'all' | 'inbound' | 'outbound'
let currentHandlerFilter = '';
let currentPage = 0;
let totalPackets = 0;
let connections = [];
let pollTimer = null;

// ── Init ──
document.addEventListener('DOMContentLoaded', () => {
  refreshAll();
  pollTimer = setInterval(refreshAll, POLL_INTERVAL);
});

async function refreshAll() {
  await Promise.all([
    refreshSystem(),
    refreshConnections(),
    refreshPackets(),
  ]);
}

// ── API helpers ──
async function api(path, opts = {}) {
  try {
    const resp = await fetch(`${API_BASE}${path}`, opts);
    return await resp.json();
  } catch (e) {
    console.warn('[API]', path, e);
    return null;
  }
}

// ── System Info ──
async function refreshSystem() {
  const res = await api('/system');
  if (!res || res.status !== 'ok') {
    setOffline();
    return;
  }
  setOnline();
  const d = res.data;
  document.getElementById('statConnections').textContent = d.connections_count ?? 0;
  document.getElementById('statPackets').textContent = formatNumber(d.packets_count ?? 0);

  if (d.heartbeat && d.heartbeat.time) {
    const tsMs = d.heartbeat.time < 1e12 ? d.heartbeat.time * 1000 : d.heartbeat.time;
    const ago = timeSince(tsMs);
    document.getElementById('statHeartbeat').textContent = ago;
  } else {
    document.getElementById('statHeartbeat').textContent = '等待中';
  }
}

// ── Connections ──
async function refreshConnections() {
  const res = await api('/connections');
  if (!res || res.status !== 'ok') return;

  connections = res.data || [];
  const container = document.getElementById('connectionList');

  if (connections.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">📡</div>
        <div class="empty-state-text">暂无活跃连接</div>
        <div class="empty-state-sub">等待客户端接入</div>
      </div>`;
    return;
  }

  container.innerHTML = connections.map(conn => {
    const name = conn.client_name || '未命名客户端';
    const isSelected = currentHandlerFilter === conn.app_id;
    return `
      <div class="connection-row ${isSelected ? 'selected' : ''}"
           onclick="toggleHandlerFilter('${conn.app_id}')">
        <div class="conn-avatar ws">🔌</div>
        <div class="conn-info">
          <div class="conn-name">${escapeHtml(name)}</div>
          <div class="conn-id">${conn.app_id}</div>
        </div>
        <div class="conn-chevron">›</div>
      </div>`;
  }).join('');
}

function toggleHandlerFilter(handlerId) {
  if (currentHandlerFilter === handlerId) {
    currentHandlerFilter = '';
  } else {
    currentHandlerFilter = handlerId;
  }
  currentPage = 0;
  refreshConnections();
  refreshPackets();
}

// ── Filter controls ──
function setLink(link, el) {
  currentLinkFilter = link;
  document.querySelectorAll('#linkTabs .filter-chip').forEach(c => c.classList.remove('active'));
  el.classList.add('active');
  currentPage = 0;
  refreshPackets();
}

function setDirection(dir, el) {
  currentDirFilter = dir;
  document.querySelectorAll('#dirTabs .filter-chip').forEach(c => c.classList.remove('active'));
  el.classList.add('active');
  currentPage = 0;
  refreshPackets();
}

// ── Packets ──
async function refreshPackets() {
  const offset = currentPage * PAGE_SIZE;
  let url = `/packets?offset=${offset}&limit=${PAGE_SIZE}`;
  if (currentHandlerFilter) {
    url += `&handler_id=${currentHandlerFilter}`;
  }

  updateRefreshIndicator(true);
  const res = await api(url);
  updateRefreshIndicator(false);

  if (!res || res.status !== 'ok') return;

  const { packets, total } = res.data;
  totalPackets = total;
  const container = document.getElementById('packetList');

  let filtered = packets || [];

  // Apply link filter
  if (currentLinkFilter !== 'all') {
    filtered = filtered.filter(p => p.link === currentLinkFilter);
  }
  // Apply direction filter
  if (currentDirFilter !== 'all') {
    filtered = filtered.filter(p => p.direction === currentDirFilter);
  }

  if (filtered.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">📭</div>
        <div class="empty-state-text">暂无数据包</div>
        <div class="empty-state-sub">${currentHandlerFilter ? '该连接暂无记录' : '等待连接建立后将自动记录'}</div>
      </div>`;
    document.getElementById('pagination').style.display = 'none';
    return;
  }

  container.innerHTML = filtered.map(p => {
    const isInbound = p.direction === 'inbound';
    const dirClass = isInbound ? 'inbound' : 'outbound';
    const dirIcon = isInbound ? '↓' : '↑';
    const linkLabel = p.link === 'ntqq' ? 'NTQQ' : '客户端';
    const typeName = extractType(p.data);
    const preview = JSON.stringify(p.data).slice(0, 120);
    const clientLabel = p.client_name || p.handler_id?.slice(0, 8) || '';

    return `
      <div class="packet-row" onclick='showPacketDetail(${escapeJsonAttr(JSON.stringify(p))})'>
        <div class="packet-dir ${dirClass}">${dirIcon}</div>
        <div class="packet-body">
          <div class="packet-meta">
            <span class="packet-type">${escapeHtml(typeName)}</span>
            <span class="packet-link-badge ${p.link}">${linkLabel}</span>
            ${clientLabel ? `<span class="packet-client">${escapeHtml(clientLabel)}</span>` : ''}
            <span class="packet-time">${formatTime(p.timestamp)}</span>
          </div>
          <div class="packet-preview">${escapeHtml(preview)}</div>
        </div>
      </div>`;
  }).join('');

  // Pagination
  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE));
  const pag = document.getElementById('pagination');
  pag.style.display = total > PAGE_SIZE ? 'flex' : 'none';
  document.getElementById('prevBtn').disabled = currentPage === 0;
  document.getElementById('nextBtn').disabled = currentPage >= totalPages - 1;
  document.getElementById('pageInfo').textContent = `${currentPage + 1} / ${totalPages}`;
}

function prevPage() {
  if (currentPage > 0) { currentPage--; refreshPackets(); }
}

function nextPage() {
  const totalPages = Math.ceil(totalPackets / PAGE_SIZE);
  if (currentPage < totalPages - 1) { currentPage++; refreshPackets(); }
}

// ── Packet Detail Modal ──
function showPacketDetail(packet) {
  const modal = document.getElementById('packetModal');
  const content = document.getElementById('modalContent');

  const isInbound = packet.direction === 'inbound';
  const linkLabel = packet.link === 'ntqq' ? 'NTQQ ↔ perpetua' : 'perpetua ↔ 客户端';
  const dirLabel = isInbound ? '入站 (→ perpetua)' : '出站 (perpetua →)';
  const clientLabel = packet.client_name || '-';

  content.innerHTML = `
    <div class="detail-row">
      <span class="detail-label">链路</span>
      <span class="detail-value">${linkLabel}</span>
    </div>
    <div class="detail-row">
      <span class="detail-label">方向</span>
      <span class="detail-value">${dirLabel}</span>
    </div>
    ${packet.link === 'client' ? `
    <div class="detail-row">
      <span class="detail-label">客户端</span>
      <span class="detail-value">${escapeHtml(clientLabel)}</span>
    </div>
    <div class="detail-row">
      <span class="detail-label">Handler</span>
      <span class="detail-value" style="font-family:'SF Mono','Menlo','Consolas',monospace;font-size:12px">${packet.handler_id}</span>
    </div>` : ''}
    <div class="detail-row">
      <span class="detail-label">时间</span>
      <span class="detail-value">${new Date(packet.timestamp).toLocaleString('zh-CN')}</span>
    </div>
    <div class="detail-row" style="border:none">
      <span class="detail-label">ID</span>
      <span class="detail-value" style="font-family:'SF Mono','Menlo','Consolas',monospace;font-size:12px">${packet.id}</span>
    </div>
    <div class="json-block">${syntaxHighlight(packet.data)}</div>
    <div class="trace-section" id="traceSection">
      <div class="trace-title">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
        关联转发链路
      </div>
      <div class="trace-loading" id="traceLoading">加载中…</div>
    </div>
  `;

  modal.classList.add('visible');
  document.body.style.overflow = 'hidden';

  // Load trace data
  loadTrace(packet.id);
}

async function loadTrace(packetId) {
  const res = await api(`/packets/trace?id=${packetId}`);
  const container = document.getElementById('traceSection');
  const loading = document.getElementById('traceLoading');

  if (!res || res.status !== 'ok' || !res.data.related || res.data.related.length <= 1) {
    loading.textContent = '无关联数据包';
    return;
  }

  const related = res.data.related;
  loading.remove();

  const items = related.map((p, i) => {
    const isCurrent = p.id === packetId;
    const dirIcon = p.direction === 'inbound' ? '↓' : '↑';
    const linkLabel = p.link === 'ntqq' ? 'NTQQ' : '客户端';
    const desc = describePacketFlow(p);
    const timeStr = new Date(p.timestamp).toLocaleString('zh-CN');
    const connector = i < related.length - 1
      ? `<div class="trace-connector"><svg width="12" height="16" viewBox="0 0 12 16" fill="none" stroke="currentColor" stroke-width="1.5"><line x1="6" y1="0" x2="6" y2="16"/><polyline points="3 12 6 16 9 12"/></svg></div>`
      : '';

    return `
      <div class="trace-item ${isCurrent ? 'current' : ''}"
           ${isCurrent ? '' : `onclick="jumpToPacket('${p.id}')"`}>
        <div class="trace-dir ${p.direction}">${dirIcon}</div>
        <div class="trace-info">
          <div class="trace-info-top">
            <span>${escapeHtml(desc)}</span>
            <span class="trace-link-badge ${p.link}">${linkLabel}</span>
          </div>
          <div class="trace-info-sub">${timeStr}${p.client_name ? ' · ' + escapeHtml(p.client_name) : ''}</div>
        </div>
        ${isCurrent ? '<span style="font-size:11px;color:var(--accent);font-weight:500">当前</span>' : '<span class="trace-arrow">›</span>'}
      </div>${connector}`;
  }).join('');

  container.innerHTML += `<div class="trace-list">${items}</div>`;
}

async function jumpToPacket(packetId) {
  const res = await api(`/packets/trace?id=${packetId}`);
  if (!res || res.status !== 'ok') return;
  showPacketDetail(res.data.source);
}

function describePacketFlow(p) {
  if (p.link === 'ntqq' && p.direction === 'inbound') return 'NTQQ → perpetua';
  if (p.link === 'ntqq' && p.direction === 'outbound') return 'perpetua → NTQQ';
  if (p.link === 'client' && p.direction === 'outbound') return 'perpetua → 客户端';
  if (p.link === 'client' && p.direction === 'inbound') return '客户端 → perpetua';
  return p.link + ' ' + p.direction;
}

function closeModal() {
  document.getElementById('packetModal').classList.remove('visible');
  document.body.style.overflow = '';
}

function closeModalOutside(e) {
  if (e.target === e.currentTarget) closeModal();
}

// ── Delete Modal ──
function showDeleteConfirm() {
  document.getElementById('deleteModal').classList.add('visible');
  document.body.style.overflow = 'hidden';
}

function closeDeleteModal() {
  document.getElementById('deleteModal').classList.remove('visible');
  document.body.style.overflow = '';
}

function closeDeleteOutside(e) {
  if (e.target === e.currentTarget) closeDeleteModal();
}

async function confirmDelete() {
  const before = Date.now();
  const res = await api(`/packets?before=${before}`, { method: 'DELETE' });
  closeDeleteModal();
  if (res && res.status === 'ok') {
    currentPage = 0;
    refreshAll();
  }
}

// ── Helpers ──
function setOnline() {
  document.getElementById('statusDot').classList.remove('offline');
  document.getElementById('statusText').textContent = '运行中';
}

function setOffline() {
  document.getElementById('statusDot').classList.add('offline');
  document.getElementById('statusText').textContent = '离线';
}

function updateRefreshIndicator(loading) {
  const el = document.getElementById('refreshIndicator');
  if (loading) {
    el.innerHTML = '<span class="spin">↻</span> 刷新中';
  } else {
    el.textContent = formatTime(Date.now());
  }
}

function extractType(data) {
  if (!data) return '未知';
  if (data.post_type) {
    if (data.post_type === 'message') return `消息 · ${data.message_type || ''}`;
    if (data.post_type === 'notice') return `通知 · ${data.notice_type || ''}`;
    if (data.post_type === 'request') return `请求 · ${data.request_type || ''}`;
    if (data.post_type === 'meta_event') return `元事件 · ${data.meta_event_type || ''}`;
    return data.post_type;
  }
  if (data.action) return `API · ${data.action}`;
  if (data.status !== undefined) return 'API 响应';
  return '数据';
}

function formatTime(ts) {
  const d = new Date(ts);
  const now = new Date();
  const pad = n => String(n).padStart(2, '0');

  if (d.toDateString() === now.toDateString()) {
    return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
  }
  return `${pad(d.getMonth()+1)}/${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function timeSince(ts) {
  const seconds = Math.floor((Date.now() - ts) / 1000);
  if (seconds < 5) return '刚刚';
  if (seconds < 60) return `${seconds}秒前`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}分钟前`;
  const hours = Math.floor(minutes / 60);
  return `${hours}小时前`;
}

function formatNumber(n) {
  if (n >= 10000) return (n / 10000).toFixed(1) + '万';
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k';
  return String(n);
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}

function escapeJsonAttr(str) {
  return str.replace(/'/g, '&#39;').replace(/"/g, '&quot;');
}

function syntaxHighlight(obj) {
  const json = JSON.stringify(obj, null, 2);
  return json.replace(
    /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
    match => {
      let cls = 'json-number';
      if (/^"/.test(match)) {
        cls = /:$/.test(match) ? 'json-key' : 'json-string';
      } else if (/true|false/.test(match)) {
        cls = 'json-bool';
      } else if (/null/.test(match)) {
        cls = 'json-null';
      }
      return `<span class="${cls}">${match}</span>`;
    }
  );
}
