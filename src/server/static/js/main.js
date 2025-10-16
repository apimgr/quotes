/**
 * Quotes API - Main JavaScript
 * Utilities for theme toggle, toast notifications, modals, and API helpers
 * Version: 1.0.0
 */

// ============================================
// INITIALIZATION
// ============================================

document.addEventListener('DOMContentLoaded', function() {
  initTheme();
  initMobileMenu();
  initKeyboardShortcuts();
  console.log('Quotes API loaded');
});

// ============================================
// THEME MANAGEMENT
// ============================================

function initTheme() {
  const savedTheme = localStorage.getItem('theme') || 'dark';
  document.documentElement.setAttribute('data-theme', savedTheme);
  updateThemeIcon(savedTheme);
}

function toggleTheme() {
  const currentTheme = document.documentElement.getAttribute('data-theme');
  const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

  document.documentElement.setAttribute('data-theme', newTheme);
  localStorage.setItem('theme', newTheme);
  updateThemeIcon(newTheme);

  showToast(`Switched to ${newTheme} theme`, 'info', 2000);
}

function updateThemeIcon(theme) {
  const themeIcon = document.querySelector('.theme-icon');
  if (themeIcon) {
    themeIcon.textContent = theme === 'dark' ? '🌙' : '☀️';
  }
}

// ============================================
// MOBILE MENU
// ============================================

function initMobileMenu() {
  const toggle = document.querySelector('.mobile-menu-toggle');
  const nav = document.querySelector('#main-nav');

  if (toggle && nav) {
    toggle.addEventListener('click', function() {
      nav.classList.toggle('active');
    });

    // Close menu when clicking outside
    document.addEventListener('click', function(e) {
      if (!toggle.contains(e.target) && !nav.contains(e.target)) {
        nav.classList.remove('active');
      }
    });
  }
}

function toggleMobileMenu() {
  const nav = document.querySelector('#main-nav');
  if (nav) {
    nav.classList.toggle('active');
  }
}

// ============================================
// TOAST NOTIFICATIONS
// ============================================

function showToast(message, type = 'info', duration = 3000) {
  const container = getToastContainer();

  const toast = document.createElement('div');
  toast.className = `toast ${type}`;

  const icon = document.createElement('span');
  icon.className = 'toast-icon';
  icon.textContent = getToastIcon(type);

  const messageEl = document.createElement('div');
  messageEl.className = 'toast-message';
  messageEl.textContent = message;

  const closeBtn = document.createElement('button');
  closeBtn.className = 'toast-close';
  closeBtn.textContent = '×';
  closeBtn.onclick = () => toast.remove();

  toast.appendChild(icon);
  toast.appendChild(messageEl);
  toast.appendChild(closeBtn);

  container.appendChild(toast);

  // Auto remove
  setTimeout(() => {
    toast.style.opacity = '0';
    setTimeout(() => toast.remove(), 300);
  }, duration);
}

function getToastContainer() {
  let container = document.getElementById('toast-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'toast-container';
    document.body.appendChild(container);
  }
  return container;
}

function getToastIcon(type) {
  const icons = {
    success: '✓',
    error: '✗',
    warning: '⚠',
    info: 'ℹ'
  };
  return icons[type] || icons.info;
}

// ============================================
// MODAL DIALOGS
// ============================================

function showModal(title, content, buttons = []) {
  const modalContainer = getModalContainer();

  const modal = document.createElement('div');
  modal.className = 'modal active';

  const backdrop = document.createElement('div');
  backdrop.className = 'modal-backdrop';
  backdrop.onclick = () => closeModal(modal);

  const modalContent = document.createElement('div');
  modalContent.className = 'modal-content';

  // Header
  const header = document.createElement('div');
  header.className = 'modal-header';

  const titleEl = document.createElement('h2');
  titleEl.className = 'modal-title';
  titleEl.textContent = title;

  const closeBtn = document.createElement('button');
  closeBtn.className = 'modal-close';
  closeBtn.textContent = '×';
  closeBtn.onclick = () => closeModal(modal);

  header.appendChild(titleEl);
  header.appendChild(closeBtn);

  // Body
  const body = document.createElement('div');
  body.className = 'modal-body';
  if (typeof content === 'string') {
    body.innerHTML = content;
  } else {
    body.appendChild(content);
  }

  // Footer
  const footer = document.createElement('div');
  footer.className = 'modal-footer';

  buttons.forEach(btn => {
    const button = document.createElement('button');
    button.className = `btn ${btn.class || 'btn-secondary'}`;
    button.textContent = btn.text;
    button.onclick = () => {
      if (btn.onClick) btn.onClick();
      closeModal(modal);
    };
    footer.appendChild(button);
  });

  // Assemble modal
  modalContent.appendChild(header);
  modalContent.appendChild(body);
  if (buttons.length > 0) {
    modalContent.appendChild(footer);
  }

  modal.appendChild(backdrop);
  modal.appendChild(modalContent);
  modalContainer.appendChild(modal);
}

function closeModal(modal) {
  if (modal) {
    modal.classList.remove('active');
    setTimeout(() => modal.remove(), 300);
  } else {
    const activeModal = document.querySelector('.modal.active');
    if (activeModal) {
      activeModal.classList.remove('active');
      setTimeout(() => activeModal.remove(), 300);
    }
  }
}

function getModalContainer() {
  let container = document.getElementById('modal-container');
  if (!container) {
    container = document.createElement('div');
    container.id = 'modal-container';
    document.body.appendChild(container);
  }
  return container;
}

// ============================================
// API HELPERS
// ============================================

async function apiGet(endpoint) {
  try {
    const response = await fetch(endpoint);
    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error?.message || 'API request failed');
    }

    return data.data;
  } catch (error) {
    showToast(error.message, 'error');
    throw error;
  }
}

async function apiPost(endpoint, body) {
  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    });

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error?.message || 'API request failed');
    }

    return data.data;
  } catch (error) {
    showToast(error.message, 'error');
    throw error;
  }
}

async function apiDelete(endpoint) {
  try {
    const response = await fetch(endpoint, {
      method: 'DELETE'
    });

    const data = await response.json();

    if (!data.success) {
      throw new Error(data.error?.message || 'API request failed');
    }

    return data.data;
  } catch (error) {
    showToast(error.message, 'error');
    throw error;
  }
}

// ============================================
// KEYBOARD SHORTCUTS
// ============================================

function initKeyboardShortcuts() {
  document.addEventListener('keydown', function(e) {
    // Ignore if typing in input/textarea
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
      return;
    }

    // ESC - Close modal
    if (e.key === 'Escape') {
      closeModal();
    }

    // T - Toggle theme
    if (e.key === 't' || e.key === 'T') {
      toggleTheme();
    }
  });
}

// ============================================
// UTILITY FUNCTIONS
// ============================================

function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
}

function formatTime(dateString) {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit'
  });
}

function formatDateTime(dateString) {
  return `${formatDate(dateString)} at ${formatTime(dateString)}`;
}

function copyToClipboard(text) {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(text).then(() => {
      showToast('Copied to clipboard!', 'success', 2000);
    }).catch(() => {
      fallbackCopyToClipboard(text);
    });
  } else {
    fallbackCopyToClipboard(text);
  }
}

function fallbackCopyToClipboard(text) {
  const textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();

  try {
    document.execCommand('copy');
    showToast('Copied to clipboard!', 'success', 2000);
  } catch (err) {
    showToast('Failed to copy', 'error', 2000);
  }

  document.body.removeChild(textarea);
}

function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

function throttle(func, limit) {
  let inThrottle;
  return function executedFunction(...args) {
    if (!inThrottle) {
      func.apply(this, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

// ============================================
// LOADING SPINNER
// ============================================

function showSpinner(container) {
  const spinner = document.createElement('div');
  spinner.className = 'spinner';
  if (typeof container === 'string') {
    document.querySelector(container).appendChild(spinner);
  } else if (container) {
    container.appendChild(spinner);
  }
  return spinner;
}

function hideSpinner(spinner) {
  if (spinner && spinner.parentNode) {
    spinner.parentNode.removeChild(spinner);
  }
}

// ============================================
// EXPORT FOR MODULE USAGE (if needed)
// ============================================

if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    toggleTheme,
    showToast,
    showModal,
    closeModal,
    apiGet,
    apiPost,
    apiDelete,
    formatDate,
    formatTime,
    formatDateTime,
    copyToClipboard,
    debounce,
    throttle,
    showSpinner,
    hideSpinner
  };
}
