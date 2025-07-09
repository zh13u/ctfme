// Thêm lại biến challengeTranslations để fix lỗi và hỗ trợ đa ngôn ngữ
const challengeTranslations = {
    en: {
        challenges: 'Challenges',
        challengesDesc: 'Challenge list will be displayed here.',
        all: 'All',
        web: 'Web',
        crypto: 'Crypto',
        forensics: 'Forensics',
        pwn: 'Pwn',
        reverse: 'Reverse',
        misc: 'Misc',
        startChallenge: 'Start Challenge',
        details: 'Details',
        loginRequired: 'Please login to access challenges',
        noChallenges: 'No challenges available at the moment.',
        errorLoading: 'Error loading challenges. Please try again later.'
    },
    // Có thể thêm các ngôn ngữ khác nếu muốn
};

// Challenge page specific functionality
document.addEventListener('DOMContentLoaded', function() {
    // Initialize challenge categories
    initializeChallengeCategories();
    
    // Load challenges
    loadChallenges();
    
    // Initialize language support
    initializeLanguageSupport();
});

// Challenge categories functionality
function initializeChallengeCategories() {
    const categoryButtons = document.querySelectorAll('.category-btn');
    
    categoryButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons
            categoryButtons.forEach(btn => btn.classList.remove('active'));
            
            // Add active class to clicked button
            this.classList.add('active');
            
            // Filter challenges
            const category = this.getAttribute('data-category');
            filterChallenges(category);
        });
    });
}

// Filter challenges by category
function filterChallenges(category) {
    const challenges = Array.from(document.querySelectorAll('.challenge-card'));
    // Lấy vị trí nút category đang active
    const activeBtn = document.querySelector('.category-btn.active');
    let btnX = window.innerWidth / 2, btnY = 0;
    if (activeBtn) {
        const rect = activeBtn.getBoundingClientRect();
        btnX = rect.left + rect.width / 2;
        btnY = rect.top + rect.height / 2;
    }

    // Lọc và hiệu ứng
    let visibleCards = [];
    challenges.forEach(card => {
        const challengeCategory = card.getAttribute('data-category');
        if (category !== 'all' && challengeCategory !== category) {
            card.classList.remove('slide-in', 'show', 'fade-in', 'fade-out');
            card.classList.add('fade-out');
            setTimeout(() => {
                card.style.display = 'none';
            }, 400);
        } else {
            card.style.display = 'block';
            card.classList.remove('fade-in', 'fade-out', 'show');
            card.classList.add('slide-in');
            visibleCards.push(card);
        }
    });

    // Reveal từng card kiểu nhân bản: card đầu từ nút, các card sau từ vị trí card trước
    setTimeout(() => {
        let prevX = btnX, prevY = btnY;
        visibleCards.forEach((card, idx) => {
            // Lấy vị trí card hiện tại
            const rect = card.getBoundingClientRect();
            const cardX = rect.left + rect.width / 2;
            const cardY = rect.top + rect.height / 2;
            // Vector từ card trước đến card hiện tại
            const dx = cardX - prevX;
            const dy = cardY - prevY;
            // Slide theo hướng từ card trước, giới hạn max 120px
            const slideX = Math.max(Math.min(dx, 120), -120);
            const slideY = Math.max(Math.min(dy, 120), -120);
            card.style.setProperty('--slide-x', `${slideX}px`);
            card.style.setProperty('--slide-y', `${slideY}px`);
            // Delay tăng dần
            const delay = idx * 180;
            setTimeout(() => {
                card.classList.add('show');
            }, delay);
            // Cập nhật vị trí cho card tiếp theo
            prevX = cardX;
            prevY = cardY;
        });
    }, 30);
}

// Load challenges from API
async function loadChallenges() {
    const container = document.getElementById('challenges-container');
    const lang = localStorage.getItem('lang') || 'en';
    const t = challengeTranslations[lang] || challengeTranslations['en'];

    try {
        const res = await fetch('http://localhost:3000/api/challenges');
        if (!res.ok) throw new Error('API error');
        const challenges = await res.json();

        renderChallenges(challenges);
    } catch (error) {
        console.error('Error loading challenges:', error);
        container.innerHTML = `<p>${t.errorLoading}</p>`;
    }
}

// Render challenges in the container
function renderChallenges(challenges) {
    const container = document.getElementById('challenges-container');
    const lang = localStorage.getItem('lang') || 'en';
    const t = challengeTranslations[lang] || challengeTranslations['en'];
    
    if (challenges.length === 0) {
        container.innerHTML = `<p>${t.noChallenges}</p>`;
        return;
    }
    
    const categoryIcons = {
        web: '🌐',
        crypto: '🔐',
        forensics: '🕵️',
        pwn: '💣',
        reverse: '��',
        misc: '✨'
    };

    const challengesHTML = challenges.map(challenge => {
        const title = challenge.title || challenge.name || 'No Title';
        const category = (challenge.category || 'misc').toLowerCase();
        const points = challenge.points || challenge.score || 0;
        const difficultyRaw = challenge.difficulty || 'Easy';
        const difficulty = difficultyRaw.toLowerCase();
        const icon = categoryIcons[category] || '❓';
        return `
            <div class="challenge-card ripple-reveal" data-category="${category}" data-id="${challenge.id}">
                <div class="challenge-header">
                    <div class="challenge-title-row">
                      <h3 class="challenge-title">${title}</h3>
                      <span class="challenge-points"><span class="challenge-icon">${icon}</span> ${points} pts</span>
                    </div>
                    <div class="challenge-title-divider"></div>
                    <div class="challenge-category">
                      <span class="category-badge ${category}">${category.toUpperCase()}</span>
                      <span class="difficulty-badge ${difficulty}">${difficultyRaw}</span>
                    </div>
                </div>
                <div class="challenge-actions">
                    <button class="btn-primary" onclick="openChallenge(${challenge.id})">${t.startChallenge}</button>
                </div>
            </div>
        `;
    }).join('');
    
    container.innerHTML = challengesHTML;
}

// Open a specific challenge
function openChallenge(challengeId) {
    // Check if user is logged in
    const user = localStorage.getItem('ctfme_user');
    if (!user) {
        alert('Please login to access challenges');
        openModal('login');
        return;
    }
    
    // In a real app, this would redirect to the challenge page
    alert(`Opening challenge ${challengeId}. This would redirect to the challenge interface.`);
}

// View challenge details
function viewChallengeDetails(challengeId) {
    // In a real app, this would show a modal with challenge details
    alert(`Showing details for challenge ${challengeId}`);
}

function initializeLanguageSupport() {
    const savedLang = localStorage.getItem('lang') || 'en';
    setChallengeLang(savedLang);
}

function setChallengeLang(lang) {
    const t = challengeTranslations[lang] || challengeTranslations['en'];
    
    // Update page title and description
    const titleElement = document.getElementById('challenges-title');
    const descElement = document.getElementById('challenges-desc');
    if (titleElement) titleElement.textContent = t.challenges;
    if (descElement) descElement.textContent = t.challengesDesc;
    
    // Update category buttons
    const categoryButtons = document.querySelectorAll('.category-btn');
    categoryButtons.forEach(btn => {
        const category = btn.getAttribute('data-category');
        if (t[category]) {
            btn.textContent = t[category];
        }
    });
    
    // Update button texts
    const startButtons = document.querySelectorAll('.btn-primary');
    const detailButtons = document.querySelectorAll('.btn-secondary');
    startButtons.forEach(btn => {
        if (btn.textContent.includes('Start Challenge')) {
            btn.textContent = t.startChallenge;
        }
    });
    detailButtons.forEach(btn => {
        if (btn.textContent.includes('Details')) {
            btn.textContent = t.details;
        }
    });
}

// Override functions to use translations
function openChallenge(challengeId) {
    const user = localStorage.getItem('ctfme_user');
    const lang = localStorage.getItem('lang') || 'en';
    const t = challengeTranslations[lang] || challengeTranslations['en'];
    
    if (!user) {
        alert(t.loginRequired);
        openModal('login');
        return;
    }
    
    alert(`Opening challenge ${challengeId}. This would redirect to the challenge interface.`);
}

// Make functions globally available
window.openChallenge = openChallenge;
window.viewChallengeDetails = viewChallengeDetails;
window.setChallengeLang = setChallengeLang;

// Matrix rain effect
function startMatrixRain() {
  const canvas = document.getElementById('matrix-bg');
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  let width = window.innerWidth;
  let height = window.innerHeight;
  canvas.width = width;
  canvas.height = height;

  const fontSize = 18;
  const columns = Math.floor(width / fontSize);
  const drops = Array(columns).fill(1);

  function draw() {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.13)';
    ctx.fillRect(0, 0, width, height);

    ctx.font = fontSize + "px monospace";
    ctx.fillStyle = "#00ffe0";
    for (let i = 0; i < drops.length; i++) {
      const text = String.fromCharCode(0x30A0 + Math.random() * 96);
      ctx.fillText(text, i * fontSize, drops[i] * fontSize);

      if (drops[i] * fontSize > height && Math.random() > 0.975) {
        drops[i] = 0;
      }
      drops[i]++;
    }
  }

  setInterval(draw, 40);

  // Responsive
  window.addEventListener('resize', () => {
    width = window.innerWidth;
    height = window.innerHeight;
    canvas.width = width;
    canvas.height = height;
  });
}
startMatrixRain();

// Ripple effect and reveal cards
function triggerRippleAndReveal() {
  const overlay = document.getElementById('rippleOverlay');
  if (!overlay) return;

  // Lấy vị trí nút Challenges trên navbar
  const navBtn = document.querySelector('.nav-links a[href="challenges.html"], .nav-links a.active');
  let x = window.innerWidth / 2, y = 0;
  if (navBtn) {
    const rect = navBtn.getBoundingClientRect();
    x = rect.left + rect.width / 2;
    y = rect.top + rect.height / 2;
  }
  overlay.style.left = '0';
  overlay.style.top = '0';
  overlay.style.background = `radial-gradient(circle at ${x}px ${y}px, rgba(0,255,255,0.35) 0%, rgba(0,255,255,0.18) 30%, rgba(120,0,255,0.13) 60%, transparent 80%), radial-gradient(circle at ${x}px ${y}px, rgba(0,255,255,0.18) 0%, transparent 70%)`;

  overlay.classList.add('active');

  // Reveal tất cả .ripple-reveal theo delay dựa trên khoảng cách từ nút đến element
  const reveals = Array.from(document.querySelectorAll('.ripple-reveal'));
  reveals.forEach(el => {
    el.classList.remove('revealed');
    const rect = el.getBoundingClientRect();
    const cx = rect.left + rect.width / 2;
    const cy = rect.top + rect.height / 2;
    const dist = Math.sqrt((cx - x) ** 2 + (cy - y) ** 2);
    const delay = 600 + dist * 1.1; // ms, chậm hơn và sóng mạnh hơn
    setTimeout(() => {
      el.classList.add('revealed');
    }, delay);
  });

  // Ẩn overlay sau khi xong
  setTimeout(() => {
    overlay.classList.remove('active');
  }, 2600);
}

// Gắn sự kiện cho nút Challenges trên navbar
window.addEventListener('DOMContentLoaded', () => {
  const navBtn = document.querySelector('.nav-links a[href="challenges.html"], .nav-links a.active');
  if (navBtn) {
    navBtn.addEventListener('click', e => {
      // Nếu đang ở trang challenges.html thì trigger hiệu ứng
      if (window.location.pathname.includes('challenges.html')) {
        e.preventDefault();
        triggerRippleAndReveal();
      }
    });
  }
  // Khi load trang, nếu là trang challenges thì cũng trigger hiệu ứng
  if (window.location.pathname.includes('challenges.html')) {
    setTimeout(triggerRippleAndReveal, 200);
  }
}); 