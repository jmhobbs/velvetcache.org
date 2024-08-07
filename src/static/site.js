document.body.dataset['scheme'] = window.localStorage.getItem('scheme') || '';
document.body.dataset['theme'] = window.localStorage.getItem('theme') || '';

document.addEventListener('DOMContentLoaded', function() {
  const preloadStylesheetCallback = function () {
    this.removeEventListener('load', preloadStylesheetCallback);
    this.rel = 'stylesheet';
  };

  document.querySelectorAll('link[data-lazy]').forEach(function(lazyStylesheet) {
    lazyStylesheet.addEventListener('load', preloadStylesheetCallback);
    lazyStylesheet.href = lazyStylesheet.dataset.href;
    lazyStylesheet.rel = 'preload';
  });

  document.querySelectorAll('a.theme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['theme'] = e.target.dataset['theme'];
    window.localStorage.setItem('theme', e.target.dataset['theme']);
  }));
  document.querySelectorAll('a.scheme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['scheme'] = e.target.dataset['scheme'];
    window.localStorage.setItem('scheme', e.target.dataset['scheme']);
  }));
});
