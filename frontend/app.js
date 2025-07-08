// Dark/Light mode toggle
const modeToggle = document.getElementById('modeToggle');
const body = document.body;
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

function setMode(dark) {
  if (dark) {
    body.classList.add('dark');
    modeToggle.textContent = '☀️';
    modeToggle.title = 'Chuyển sang chế độ sáng';
    localStorage.setItem('theme', 'dark');
  } else {
    body.classList.remove('dark');
    modeToggle.textContent = '🌙';
    modeToggle.title = 'Chuyển sang chế độ tối';
    localStorage.setItem('theme', 'light');
  }
}

// Khởi tạo theme
const savedTheme = localStorage.getItem('theme');
setMode(savedTheme === 'dark' || (!savedTheme && prefersDark));

modeToggle.onclick = () => setMode(!body.classList.contains('dark'));

// Modal logic
const modal = document.getElementById('modal');
const modalLogin = document.getElementById('modal-login');
const modalRegister = document.getElementById('modal-register');

function openModal(type) {
  modal.style.display = 'flex';
  if (type === 'login') {
    modalLogin.style.display = 'block';
    modalRegister.style.display = 'none';
  } else {
    modalLogin.style.display = 'none';
    modalRegister.style.display = 'block';
  }
}
function closeModal() {
  modal.style.display = 'none';
}
window.openModal = openModal;
window.closeModal = closeModal;

// Đăng nhập/Đăng ký demo (localStorage)
function login() {
  const username = document.getElementById('login-username').value;
  if (username) {
    localStorage.setItem('ctfme_user', username);
    updateAuthUI();
    closeModal();
  }
}
function register() {
  const username = document.getElementById('register-username').value;
  if (username) {
    localStorage.setItem('ctfme_user', username);
    updateAuthUI();
    closeModal();
  }
}
function logout() {
  localStorage.removeItem('ctfme_user');
  updateAuthUI();
}
window.login = login;
window.register = register;
window.logout = logout;

// Hiển thị menu theo trạng thái đăng nhập
function updateAuthUI() {
  const user = localStorage.getItem('ctfme_user');
  document.querySelectorAll('.auth-item').forEach(e => e.style.display = user ? 'none' : 'inline-block');
  document.querySelectorAll('.profile-item').forEach(e => e.style.display = user ? 'inline-block' : 'none');
  if (user) {
    document.getElementById('profileBtn').querySelector('a').textContent = user;
  }
}
updateAuthUI();

// Năm hiện tại cho footer
const year = document.getElementById('year');
if (year) year.textContent = new Date().getFullYear();

// Đóng modal khi bấm ngoài
modal.addEventListener('click', (e) => {
  if (e.target === modal) closeModal();
});

// Đa ngôn ngữ
const translations = {
  en: {
    rules: 'Rules',
    rulesDesc: 'CTF rules, regulations and instructions...',
    challenges: 'Challenges',
    challengesDesc: 'Challenge list will be displayed here.',
    scoreboard: 'Scoreboard',
    scoreboardDesc: 'Ranking board for teams/players.',
    teams: 'Teams',
    teamsDesc: 'Information about participating teams.',
    login: 'Login',
    register: 'Register',
    profile: 'Profile',
    settings: 'Settings',
    logout: 'Logout',
    username: 'Username',
    password: 'Password',
  },
  vi: {
    rules: 'Luật chơi',
    rulesDesc: 'Luật chơi CTF, quy định và hướng dẫn...',
    challenges: 'Thử thách',
    challengesDesc: 'Danh sách thử thách sẽ hiển thị ở đây.',
    scoreboard: 'Bảng xếp hạng',
    scoreboardDesc: 'Bảng xếp hạng các đội/thi sinh.',
    teams: 'Đội',
    teamsDesc: 'Thông tin các đội tham gia.',
    login: 'Đăng nhập',
    register: 'Đăng ký',
    profile: 'Cá nhân',
    settings: 'Cài đặt',
    logout: 'Đăng xuất',
    username: 'Tên đăng nhập',
    password: 'Mật khẩu',
  },
  fr: {
    rules: 'Règles',
    rulesDesc: 'Règles CTF, règlements et instructions...',
    challenges: 'Défis',
    challengesDesc: 'La liste des défis sera affichée ici.',
    scoreboard: 'Classement',
    scoreboardDesc: 'Tableau de classement des équipes/joueurs.',
    teams: 'Équipes',
    teamsDesc: 'Informations sur les équipes participantes.',
    login: 'Connexion',
    register: 'Inscription',
    profile: 'Profil',
    settings: 'Paramètres',
    logout: 'Déconnexion',
    username: "Nom d'utilisateur",
    password: 'Mot de passe',
  },
  ja: {
    rules: 'ルール',
    rulesDesc: 'CTFのルール、規則、説明...',
    challenges: 'チャレンジ',
    challengesDesc: 'チャレンジリストがここに表示されます。',
    scoreboard: 'スコアボード',
    scoreboardDesc: 'チーム/プレイヤーのランキングボード。',
    teams: 'チーム',
    teamsDesc: '参加チームの情報。',
    login: 'ログイン',
    register: '登録',
    profile: 'プロフィール',
    settings: '設定',
    logout: 'ログアウト',
    username: 'ユーザー名',
    password: 'パスワード',
  },
  zh: {
    rules: '规则',
    rulesDesc: 'CTF规则、规定和说明...',
    challenges: '挑战',
    challengesDesc: '挑战列表将在此显示。',
    scoreboard: '排行榜',
    scoreboardDesc: '团队/选手排行榜。',
    teams: '队伍',
    teamsDesc: '参赛队伍信息。',
    login: '登录',
    register: '注册',
    profile: '个人',
    settings: '设置',
    logout: '登出',
    username: '用户名',
    password: '密码',
  },
};

function setLang(lang) {
  const t = translations[lang] || translations['en'];
  // Navbar
  document.querySelector('a[href="#rules"]').textContent = t.rules;
  document.querySelector('a[href="#challenges"]').textContent = t.challenges;
  document.querySelector('a[href="#scoreboard"]').textContent = t.scoreboard;
  document.querySelector('a[href="#teams"]').textContent = t.teams;
  document.getElementById('loginBtn').querySelector('a').textContent = t.login;
  document.getElementById('registerBtn').querySelector('a').textContent = t.register;
  document.getElementById('profileBtn').querySelector('a').textContent = t.profile;
  document.getElementById('settingBtn').querySelector('a').textContent = t.settings;
  document.getElementById('logoutBtn').querySelector('a').textContent = t.logout;
  // Section titles & desc
  document.getElementById('rules-title').textContent = t.rules;
  document.getElementById('rules-desc').textContent = t.rulesDesc;
  document.getElementById('challenges-title').textContent = t.challenges;
  document.getElementById('challenges-desc').textContent = t.challengesDesc;
  document.getElementById('scoreboard-title').textContent = t.scoreboard;
  document.getElementById('scoreboard-desc').textContent = t.scoreboardDesc;
  document.getElementById('teams-title').textContent = t.teams;
  document.getElementById('teams-desc').textContent = t.teamsDesc;
  // Modal
  document.getElementById('login-title').textContent = t.login;
  document.getElementById('login-username').placeholder = t.username;
  document.getElementById('login-password').placeholder = t.password;
  document.getElementById('login-btn').textContent = t.login;
  document.getElementById('register-title').textContent = t.register;
  document.getElementById('register-username').placeholder = t.username;
  document.getElementById('register-password').placeholder = t.password;
  document.getElementById('register-btn').textContent = t.register;
}

const langSelect = document.getElementById('langSelect');
const savedLang = localStorage.getItem('lang') || 'en';
langSelect.value = savedLang;
setLang(savedLang);

langSelect.addEventListener('change', function() {
  setLang(this.value);
  localStorage.setItem('lang', this.value);
});

// Google Translate Widget toggle (an toàn hơn)
const langToggle = document.getElementById('langToggle');
const langWidget = document.getElementById('google_translate_element');

langToggle.addEventListener('click', function(e) {
  e.stopPropagation();
  // Nếu widget chưa có nội dung, thử khởi tạo lại
  if (!langWidget.innerHTML.trim()) {
    if (typeof google !== 'undefined' && google.translate && google.translate.TranslateElement) {
      new google.translate.TranslateElement({
        pageLanguage: 'en',
        includedLanguages: 'en,vi,fr,ja,zh,ko,ru,es,de,th,pt,it',
        layout: google.translate.TranslateElement.InlineLayout.SIMPLE
      }, 'google_translate_element');
    }
  }
  langWidget.classList.toggle('lang-widget-shown');
  langWidget.classList.toggle('lang-widget-hidden');
});
document.addEventListener('click', function(e) {
  if (!langWidget.contains(e.target) && e.target !== langToggle) {
    langWidget.classList.add('lang-widget-hidden');
    langWidget.classList.remove('lang-widget-shown');
  }
}); 