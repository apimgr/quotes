/**
 * Extra JavaScript for MkDocs Material
 * Quotes API Documentation
 */

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', function() {
  console.log('Quotes API Documentation loaded');

  // Add copy button feedback
  const copyButtons = document.querySelectorAll('.md-clipboard');
  copyButtons.forEach(button => {
    button.addEventListener('click', function() {
      const originalTitle = this.title;
      this.title = 'Copied!';
      setTimeout(() => {
        this.title = originalTitle;
      }, 2000);
    });
  });

  // Add external link icons
  const externalLinks = document.querySelectorAll('a[href^="http"]');
  externalLinks.forEach(link => {
    if (!link.hostname.includes('quotes.readthedocs.io') &&
        !link.hostname.includes('localhost')) {
      link.setAttribute('target', '_blank');
      link.setAttribute('rel', 'noopener noreferrer');
    }
  });

  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
      const targetId = this.getAttribute('href').substring(1);
      const targetElement = document.getElementById(targetId);

      if (targetElement) {
        e.preventDefault();
        targetElement.scrollIntoView({
          behavior: 'smooth',
          block: 'start'
        });

        // Update URL without jumping
        history.pushState(null, null, `#${targetId}`);
      }
    });
  });

  // Add keyboard shortcuts hint
  const keyboardShortcuts = {
    '/': 'Focus search',
    'f': 'Search in page',
    's': 'Focus search',
    'p': 'Previous page',
    'n': 'Next page'
  };

  // Handle keyboard shortcuts
  document.addEventListener('keydown', function(e) {
    // Ignore if user is typing in an input
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
      return;
    }

    switch(e.key) {
      case '/':
      case 's':
        e.preventDefault();
        const searchInput = document.querySelector('.md-search__input');
        if (searchInput) searchInput.focus();
        break;
      case 'f':
        // Browser's default find functionality
        break;
    }
  });

  // Add version info to footer
  const footer = document.querySelector('.md-footer-meta');
  if (footer) {
    const versionInfo = document.createElement('div');
    versionInfo.className = 'md-footer-version';
    versionInfo.style.textAlign = 'center';
    versionInfo.style.padding = '1rem';
    versionInfo.style.color = '#6272a4';
    versionInfo.innerHTML = 'Version 0.0.1 | Last updated: 2025-10-14';
    footer.appendChild(versionInfo);
  }

  // Highlight current section in navigation
  const currentPath = window.location.pathname;
  const navLinks = document.querySelectorAll('.md-nav__link');
  navLinks.forEach(link => {
    if (link.getAttribute('href') === currentPath) {
      link.classList.add('md-nav__link--active');
    }
  });

  // Add "Back to Top" button functionality
  const backToTop = document.createElement('button');
  backToTop.innerHTML = '↑';
  backToTop.className = 'back-to-top';
  backToTop.style.cssText = `
    position: fixed;
    bottom: 2rem;
    right: 2rem;
    width: 3rem;
    height: 3rem;
    border-radius: 50%;
    background-color: #bd93f9;
    color: #f8f8f2;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.3s;
    z-index: 1000;
  `;
  document.body.appendChild(backToTop);

  // Show/hide back to top button
  window.addEventListener('scroll', function() {
    if (window.pageYOffset > 300) {
      backToTop.style.opacity = '1';
    } else {
      backToTop.style.opacity = '0';
    }
  });

  // Scroll to top on click
  backToTop.addEventListener('click', function() {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  });

  // Add loading animation for code blocks
  const codeBlocks = document.querySelectorAll('pre code');
  codeBlocks.forEach(block => {
    block.style.opacity = '0';
    block.style.transition = 'opacity 0.3s';
    setTimeout(() => {
      block.style.opacity = '1';
    }, 100);
  });
});

// Analytics (placeholder - add your analytics code here)
// Example: Google Analytics, Plausible, etc.

// Add print styles
if (window.matchMedia) {
  const mediaQuery = window.matchMedia('print');
  mediaQuery.addListener(function(mql) {
    if (mql.matches) {
      console.log('Preparing for print...');
      // Hide navigation and other non-essential elements
      document.querySelectorAll('.md-header, .md-sidebar').forEach(el => {
        el.style.display = 'none';
      });
    }
  });
}
