// Thêm lại biến translations để fix lỗi và hỗ trợ đa ngôn ngữ cho giao diện
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
  // ... các ngôn ngữ khác nếu có
};

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

// Navigation active state
function setActiveNavItem(clickedLink) {
  // Remove active class from all nav links
  document.querySelectorAll('.nav-links a').forEach(link => {
    link.classList.remove('active');
  });
  
  // Add active class to clicked link
  clickedLink.classList.add('active');
}

// Add click event listeners to nav links
document.querySelectorAll('.nav-links a').forEach(link => {
  link.addEventListener('click', function(e) {
    setActiveNavItem(this);
  });
});

function setLang(lang) {
  const t = translations[lang] || translations['en'];
  // Navbar
  const rulesLink = document.querySelector('a[href="#rules"]');
  if (rulesLink) rulesLink.textContent = t.rules;
  const challengesLink = document.querySelector('a[href="challenges.html"]');
  if (challengesLink) challengesLink.textContent = t.challenges;
  const scoreboardLink = document.querySelector('a[href="#scoreboard"]');
  if (scoreboardLink) scoreboardLink.textContent = t.scoreboard;
  const teamsLink = document.querySelector('a[href="#teams"]');
  if (teamsLink) teamsLink.textContent = t.teams;
  const loginBtn = document.getElementById('loginBtn');
  if (loginBtn) loginBtn.querySelector('a').textContent = t.login;
  const registerBtn = document.getElementById('registerBtn');
  if (registerBtn) registerBtn.querySelector('a').textContent = t.register;
  const profileBtn = document.getElementById('profileBtn');
  if (profileBtn) profileBtn.querySelector('a').textContent = t.profile;
  const settingBtn = document.getElementById('settingBtn');
  if (settingBtn) settingBtn.querySelector('a').textContent = t.settings;
  const logoutBtn = document.getElementById('logoutBtn');
  if (logoutBtn) logoutBtn.querySelector('a').textContent = t.logout;
  // Section titles & desc
  const rulesTitle = document.getElementById('rules-title');
  if (rulesTitle) rulesTitle.textContent = t.rules;
  const rulesDesc = document.getElementById('rules-desc');
  if (rulesDesc) rulesDesc.textContent = t.rulesDesc;
  const challengesTitle = document.getElementById('challenges-title');
  if (challengesTitle) challengesTitle.textContent = t.challenges;
  const challengesDesc = document.getElementById('challenges-desc');
  if (challengesDesc) challengesDesc.textContent = t.challengesDesc;
  const scoreboardTitle = document.getElementById('scoreboard-title');
  if (scoreboardTitle) scoreboardTitle.textContent = t.scoreboard;
  const scoreboardDesc = document.getElementById('scoreboard-desc');
  if (scoreboardDesc) scoreboardDesc.textContent = t.scoreboardDesc;
  const teamsTitle = document.getElementById('teams-title');
  if (teamsTitle) teamsTitle.textContent = t.teams;
  const teamsDesc = document.getElementById('teams-desc');
  if (teamsDesc) teamsDesc.textContent = t.teamsDesc;
  // Modal
  const loginTitle = document.getElementById('login-title');
  if (loginTitle) loginTitle.textContent = t.login;
  const loginUsername = document.getElementById('login-username');
  if (loginUsername) loginUsername.placeholder = t.username;
  const loginPassword = document.getElementById('login-password');
  if (loginPassword) loginPassword.placeholder = t.password;
  const loginBtn2 = document.getElementById('login-btn');
  if (loginBtn2) loginBtn2.textContent = t.login;
  const registerTitle = document.getElementById('register-title');
  if (registerTitle) registerTitle.textContent = t.register;
  const registerUsername = document.getElementById('register-username');
  if (registerUsername) registerUsername.placeholder = t.username;
  const registerPassword = document.getElementById('register-password');
  if (registerPassword) registerPassword.placeholder = t.password;
  const registerBtn2 = document.getElementById('register-btn');
  if (registerBtn2) registerBtn2.textContent = t.register;
}

const langSelect = document.getElementById('langSelect');
const savedLang = localStorage.getItem('lang') || 'en';
setLang(savedLang);
if (langSelect) {
  langSelect.value = savedLang;
  langSelect.addEventListener('change', function() {
    setLang(this.value);
    localStorage.setItem('lang', this.value);
  });
}

// Google Translate Widget toggle (an toàn hơn)
const langToggle = document.getElementById('langToggle');
const langWidget = document.getElementById('google_translate_element');

if (langToggle && langWidget) {
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
}

 