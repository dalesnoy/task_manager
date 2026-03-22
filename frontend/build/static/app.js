// === Конфигурация ===
const API = '/api';
let token = localStorage.getItem('token');
let currentProjectId = null;
let currentProjectTitle = '';

// === Утилиты ===
function api(method, path, body) {
    const headers = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = 'Bearer ' + token;

    return fetch(API + path, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
    }).then(async res => {
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Ошибка сервера');
        return data;
    });
}

function showError(id, msg) {
    const el = document.getElementById(id);
    el.textContent = msg;
    el.classList.remove('hidden');
    setTimeout(() => el.classList.add('hidden'), 4000);
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU');
}

// === Навигация ===
function showPage(page) {
    document.querySelectorAll('.page').forEach(p => p.classList.add('hidden'));
    document.getElementById(page + '-page').classList.remove('hidden');
}

function showTab(tab) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    event.target.classList.add('active');
    document.getElementById('login-form').classList.toggle('hidden', tab !== 'login');
    document.getElementById('register-form').classList.toggle('hidden', tab !== 'register');
}

// === Авторизация ===
async function register(e) {
    e.preventDefault();
    try {
        await api('POST', '/auth/register', {
            name: document.getElementById('register-name').value,
            email: document.getElementById('register-email').value,
            password: document.getElementById('register-password').value,
        });
        // Автологин после регистрации
        await login(e, document.getElementById('register-email').value, document.getElementById('register-password').value);
    } catch (err) {
        showError('auth-error', err.message);
    }
}

async function login(e, emailOverride, passwordOverride) {
    if (e) e.preventDefault();
    try {
        const email = emailOverride || document.getElementById('login-email').value;
        const password = passwordOverride || document.getElementById('login-password').value;

        const data = await api('POST', '/auth/login', { email, password });
        token = data.token;
        localStorage.setItem('token', token);

        // Декодируем имя из токена (простой base64 decode payload)
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            document.getElementById('user-name').textContent = email;
        } catch(_) {}

        document.getElementById('navbar').classList.remove('hidden');
        showProjects();
    } catch (err) {
        showError('auth-error', err.message);
    }
}

function logout() {
    token = null;
    localStorage.removeItem('token');
    document.getElementById('navbar').classList.add('hidden');
    showPage('auth');
}

// === Проекты ===
async function showProjects() {
    showPage('projects');
    try {
        const projects = await api('GET', '/projects');
        const list = document.getElementById('projects-list');

        if (!projects || projects.length === 0) {
            list.innerHTML = '<p style="color:#888;text-align:center;padding:40px">Нет проектов. Создайте первый!</p>';
            return;
        }

        list.innerHTML = projects.map(p => `
            <div class="project-card" onclick="openProject(${p.id}, '${p.title.replace(/'/g, "\\'")}')">
                <h3>${escapeHtml(p.title)}</h3>
                <p>${escapeHtml(p.description || 'Без описания')}</p>
                <div class="meta">Создан: ${formatDate(p.created_at)}</div>
                <div class="project-actions" onclick="event.stopPropagation()">
                    <button class="btn btn-danger btn-sm" onclick="deleteProject(${p.id})">Удалить</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        if (err.message.includes('авторизация') || err.message.includes('токен')) {
            logout();
        }
    }
}

function showCreateProject() {
    document.getElementById('create-project-form').classList.remove('hidden');
}

function hideCreateProject() {
    document.getElementById('create-project-form').classList.add('hidden');
    document.getElementById('project-title').value = '';
    document.getElementById('project-desc').value = '';
}

async function createProject() {
    const title = document.getElementById('project-title').value;
    if (!title) return;

    try {
        await api('POST', '/projects', {
            title,
            description: document.getElementById('project-desc').value,
        });
        hideCreateProject();
        showProjects();
    } catch (err) {
        alert(err.message);
    }
}

async function deleteProject(id) {
    if (!confirm('Удалить проект и все его задачи?')) return;
    try {
        await api('DELETE', '/projects/' + id);
        showProjects();
    } catch (err) {
        alert(err.message);
    }
}

// === Задачи ===
function openProject(id, title) {
    currentProjectId = id;
    currentProjectTitle = title;
    document.getElementById('project-name').textContent = title;
    showPage('tasks');
    loadTasks();
}

async function loadTasks() {
    const status = document.getElementById('filter-status').value;
    const priority = document.getElementById('filter-priority').value;
    const sort = document.getElementById('filter-sort').value;

    let query = `?limit=100&sort=${sort}`;
    if (status) query += `&status=${status}`;
    if (priority) query += `&priority=${priority}`;

    try {
        const data = await api('GET', `/projects/${currentProjectId}/tasks${query}`);
        const tasks = data.tasks || [];

        // Разделяем по колонкам
        const todo = tasks.filter(t => t.status === 'todo');
        const inProgress = tasks.filter(t => t.status === 'in_progress');
        const done = tasks.filter(t => t.status === 'done');

        document.getElementById('tasks-todo').innerHTML = todo.map(taskCard).join('') || '<p style="color:#aaa;font-size:13px;padding:8px">Пусто</p>';
        document.getElementById('tasks-in-progress').innerHTML = inProgress.map(taskCard).join('') || '<p style="color:#aaa;font-size:13px;padding:8px">Пусто</p>';
        document.getElementById('tasks-done').innerHTML = done.map(taskCard).join('') || '<p style="color:#aaa;font-size:13px;padding:8px">Пусто</p>';
    } catch (err) {
        alert(err.message);
    }
}

function taskCard(task) {
    const tags = (task.tags || []).map(t =>
        `<span class="tag" style="background:${t.color}22;color:${t.color}">${escapeHtml(t.name)}</span>`
    ).join('');

    const deadlineStr = task.deadline ? `<span class="deadline">${formatDate(task.deadline)}</span>` : '';

    const statusButtons = [];
    if (task.status !== 'todo') statusButtons.push(`<button onclick="changeStatus(${task.id},'todo')">К выполнению</button>`);
    if (task.status !== 'in_progress') statusButtons.push(`<button onclick="changeStatus(${task.id},'in_progress')">В работу</button>`);
    if (task.status !== 'done') statusButtons.push(`<button onclick="changeStatus(${task.id},'done')">Готово</button>`);

    return `
        <div class="task-card">
            <h4>${escapeHtml(task.title)}</h4>
            ${task.description ? `<div class="task-desc">${escapeHtml(task.description)}</div>` : ''}
            <div class="task-meta">
                <span class="priority-badge priority-${task.priority}">${task.priority}</span>
                ${deadlineStr}
            </div>
            ${tags ? `<div>${tags}</div>` : ''}
            <div class="task-actions">
                ${statusButtons.join('')}
                <button onclick="deleteTask(${task.id})" style="color:#e74c3c">Удалить</button>
            </div>
        </div>
    `;
}

function showCreateTask() {
    document.getElementById('create-task-form').classList.remove('hidden');
}

function hideCreateTask() {
    document.getElementById('create-task-form').classList.add('hidden');
    document.getElementById('task-title').value = '';
    document.getElementById('task-desc').value = '';
    document.getElementById('task-deadline').value = '';
}

async function createTask() {
    const title = document.getElementById('task-title').value;
    if (!title) return;

    const body = {
        title,
        description: document.getElementById('task-desc').value,
        priority: document.getElementById('task-priority').value,
    };

    const deadline = document.getElementById('task-deadline').value;
    if (deadline) {
        body.deadline = new Date(deadline).toISOString();
    }

    try {
        await api('POST', `/projects/${currentProjectId}/tasks`, body);
        hideCreateTask();
        loadTasks();
    } catch (err) {
        alert(err.message);
    }
}

async function changeStatus(taskId, status) {
    try {
        await api('PUT', '/tasks/' + taskId, { title: '', status });
        loadTasks();
    } catch (err) {
        alert(err.message);
    }
}

async function deleteTask(id) {
    if (!confirm('Удалить задачу?')) return;
    try {
        await api('DELETE', '/tasks/' + id);
        loadTasks();
    } catch (err) {
        alert(err.message);
    }
}

// === Защита от XSS ===
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// === Инициализация ===
if (token) {
    document.getElementById('navbar').classList.remove('hidden');
    showProjects();
} else {
    showPage('auth');
}
